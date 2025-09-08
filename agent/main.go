package main

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "bytes"
    "io"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/gorilla/mux"
    "github.com/JasSra/dockyard/agent/internal/nginx"
    "github.com/JasSra/dockyard/agent/internal/ports"
    dkr "github.com/JasSra/dockyard/agent/internal/docker"
    cmp "github.com/JasSra/dockyard/agent/internal/compose"
)

const (
    defaultPortRangeStart = 20000
    defaultPortRangeEnd   = 29999
)

// Wire minimal state
type Server struct {
    router *mux.Router
    secret string
    ports  *ports.Registry
    nginx  *nginx.Manager
    docker *dkr.Client
    nginxContainer string
    compose *cmp.Manager
}

func (s *Server) routes() {
    r := s.router
    r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200); _, _ = w.Write([]byte("ok")) }).Methods("GET")

    // HMAC protected group
    h := func(fn http.HandlerFunc) http.HandlerFunc { return s.hmac(fn) }
    r.HandleFunc("/v1/deploy", h(s.handleDeploy)).Methods("POST")
    r.HandleFunc("/v1/compose/up", h(s.handleComposeUp)).Methods("POST")
    r.HandleFunc("/v1/compose/down", h(s.handleComposeDown)).Methods("POST")
    r.HandleFunc("/v1/nginx/apply", h(s.handleNginxApply)).Methods("POST")
    r.HandleFunc("/v1/ports/allocate", h(s.handlePortsAllocate)).Methods("POST")
    r.HandleFunc("/v1/ports/free", h(s.handlePortsFree)).Methods("POST")
    r.HandleFunc("/v1/status", h(s.handleStatus)).Methods("GET")
    r.HandleFunc("/v1/logs", h(s.handleLogs)).Methods("GET")
    r.HandleFunc("/v1/ops/run", h(s.handleOpsRun)).Methods("POST")
}

func (s *Server) hmac(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        sig := r.Header.Get("X-Signature")
        if s.secret == "" {
            http.Error(w, "agent not configured", http.StatusInternalServerError)
            return
        }
    body, err := io.ReadAll(r.Body)
        if err != nil { http.Error(w, "bad body", 400); return }
    r.Body = io.NopCloser(bytes.NewReader(body))
        mac := hmac.New(sha256.New, []byte(s.secret))
        mac.Write(body)
        want := hex.EncodeToString(mac.Sum(nil))
        if sig != want {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }
    r.Body = io.NopCloser(bytes.NewReader(body))
        next(w, r)
    }
}


// Payloads (minimal stubs)
type DeployReq struct { ProjectID string `json:"projectId"`; ArchiveURL string `json:"archiveUrl,omitempty"`; KV map[string]string `json:"kv,omitempty"`; Type string `json:"type,omitempty"`; TemplateOverride string `json:"templateOverride,omitempty"`; Domain string `json:"domain,omitempty"`; HealthPath string `json:"healthPath,omitempty"` }
type ComposeUpReq struct { ProjectID string `json:"projectId"`; ImageTag string `json:"imageTag,omitempty"`; Env map[string]string `json:"env,omitempty"`; ContainerPort int `json:"containerPort,omitempty"`; HostPort int `json:"hostPort,omitempty"` }
type ComposeDownReq struct { ProjectID string `json:"projectId"`; RemoveVolumes bool `json:"removeVolumes,omitempty"` }
type NginxApplyReq struct { ProjectID string `json:"projectId"`; Host string `json:"host"`; UpstreamPort int `json:"upstreamPort"` }
type PortsAllocReq struct { ProjectID string `json:"projectId"` }
type PortsAllocResp struct { HostPort int `json:"hostPort"` }
type StatusResp struct { ProjectID string `json:"projectId"`; Container string `json:"container"`; Ports map[string]int `json:"ports"`; Health string `json:"health"`; CPU float64 `json:"cpu"`; MemBytes uint64 `json:"mem"` }
type LogsResp struct { Lines []string `json:"lines"` }
type OpsRunReq struct { OperationID string `json:"operationId"`; Action string `json:"action"`; Args map[string]any `json:"args"` }

// Handlers (stubs return 200 with minimal JSON)
func (s *Server) handleDeploy(w http.ResponseWriter, r *http.Request) {
    var req DeployReq
    _ = json.NewDecoder(r.Body).Decode(&req)
    if req.ProjectID == "" { http.Error(w, "projectId required", 400); return }
    if s.docker == nil { http.Error(w, "docker not available", 500); return }
    // Minimal: support remote URL builds (git or tarball URL). Tag as dy-{project}:latest
    tag := "dy-" + req.ProjectID + ":latest"
    if req.ArchiveURL == "" {
        writeJSON(w, map[string]any{"status":"skipped","reason":"no archiveUrl"})
        return
    }
    if err := s.docker.BuildRemote(tag, req.ArchiveURL, ""); err != nil { http.Error(w, err.Error(), 500); return }
    writeJSON(w, map[string]any{"status":"built","image": tag})
}
func (s *Server) handleComposeUp(w http.ResponseWriter, r *http.Request) {
    var req ComposeUpReq
    _ = json.NewDecoder(r.Body).Decode(&req)
    if req.ContainerPort == 0 { req.ContainerPort = 3000 }
    // Prefer Docker Engine API if available, fallback to compose no-op
    if s.docker != nil {
        name := "dy-" + req.ProjectID
        if _, err := s.docker.Create(name, req.ImageTag, req.Env, req.ContainerPort, req.HostPort); err != nil {
            http.Error(w, err.Error(), 500); return
        }
        if err := s.docker.Start(name); err != nil { http.Error(w, err.Error(), 500); return }
        writeJSON(w, map[string]any{"status":"up","driver":"engine"})
        return
    }
    if s.compose != nil {
        svc := cmp.Service{ Name: "app", Image: req.ImageTag, Env: req.Env }
        if req.HostPort != 0 { svc.Ports = map[int]int{ req.HostPort: req.ContainerPort } }
        _, _ = s.compose.Write(req.ProjectID, svc)
        if err := s.compose.Up(req.ProjectID); err != nil { http.Error(w, err.Error(), 500); return }
        writeJSON(w, map[string]any{"status":"up","driver":"compose"})
        return
    }
    http.Error(w, "no runtime configured", 500)
}
func (s *Server) handleComposeDown(w http.ResponseWriter, r *http.Request) {
    var req ComposeDownReq
    _ = json.NewDecoder(r.Body).Decode(&req)
    if s.docker != nil {
        name := "dy-" + req.ProjectID
        _ = s.docker.Remove(name, true)
        writeJSON(w, map[string]any{"status":"down","driver":"engine"})
        return
    }
    if s.compose != nil {
        if err := s.compose.Down(req.ProjectID); err != nil { http.Error(w, err.Error(), 500); return }
        writeJSON(w, map[string]any{"status":"down","driver":"compose"})
        return
    }
    http.Error(w, "no runtime configured", 500)
}
func (s *Server) handleNginxApply(w http.ResponseWriter, r *http.Request) {
    var req NginxApplyReq
    _ = json.NewDecoder(r.Body).Decode(&req)
    if s.nginx == nil { http.Error(w, "nginx not configured", 500); return }
    _, err := s.nginx.WriteServer(req.ProjectID, req.Host, req.UpstreamPort)
    if err != nil { http.Error(w, err.Error(), 500); return }
    // try nginx reload via docker kill -HUP on container name. // { SPECULATION }
    if s.docker != nil && s.nginxContainer != "" {
        if err := s.docker.Kill(s.nginxContainer, "HUP"); err != nil {
            log.Printf("warn: nginx reload failed: %v", err)
        }
    }
    writeJSON(w, map[string]any{"status":"written","reloaded": s.docker != nil && s.nginxContainer != ""})
}
func (s *Server) handlePortsAllocate(w http.ResponseWriter, r *http.Request) {
    var req PortsAllocReq
    _ = json.NewDecoder(r.Body).Decode(&req)
    if s.ports == nil { http.Error(w, "ports not configured", 500); return }
    p, err := s.ports.Allocate(req.ProjectID)
    if err != nil { http.Error(w, err.Error(), 500); return }
    writeJSON(w, PortsAllocResp{HostPort: p})
}
func (s *Server) handlePortsFree(w http.ResponseWriter, r *http.Request) {
    var req PortsAllocReq
    _ = json.NewDecoder(r.Body).Decode(&req)
    if s.ports == nil { http.Error(w, "ports not configured", 500); return }
    _, err := s.ports.Free(req.ProjectID)
    if err != nil { http.Error(w, err.Error(), 500); return }
    writeJSON(w, map[string]any{"status":"freed"})
}
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) { writeJSON(w, StatusResp{ProjectID: r.URL.Query().Get("projectId"), Container: "stub", Ports: map[string]int{"host":defaultPortRangeStart,"container":3000}, Health: "ok"}) }
func (s *Server) handleLogs(w http.ResponseWriter, r *http.Request) { writeJSON(w, LogsResp{Lines: []string{"stub log"}}) }
func (s *Server) handleOpsRun(w http.ResponseWriter, r *http.Request) { writeJSON(w, map[string]any{"status":"accepted"}) }

func writeJSON(w http.ResponseWriter, v any) {
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(v)
}

func main() {
    r := mux.NewRouter()
    portsFile := "/data/ports.json"
    pr, _ := ports.NewRegistry(portsFile, defaultPortRangeStart, defaultPortRangeEnd)
    nm := &nginx.Manager{ConfDir: "/etc/nginx/conf.d", UpstreamHost: getenvDefault("NGINX_UPSTREAM_HOST", "127.0.0.1")}
    // set up minimal docker client using unix socket
    sock := "/var/run/docker.sock"
    var dk *dkr.Client
    if _, err := os.Stat(sock); err == nil { dk = dkr.NewClient(sock) }
    s := &Server{
        router: r,
        secret: os.Getenv("AGENT_SHARED_SECRET"),
        ports: pr,
        nginx: nm,
        docker: dk,
        nginxContainer: getenvDefault("NGINX_CONTAINER", "nginx"),
        compose: &cmp.Manager{Dir: "/data"},
    }
    s.routes()
    addr := ":8080" // {{ SPECULATION }} agent listens on 8080 inside host
    srv := &http.Server{Addr: addr, Handler: r, ReadTimeout: 15 * time.Second, WriteTimeout: 30 * time.Second}
    log.Printf("agent listening on %s", addr)
    log.Fatal(srv.ListenAndServe())
}

func getenvDefault(k, d string) string { v := os.Getenv(k); if v == "" { return d }; return v }
