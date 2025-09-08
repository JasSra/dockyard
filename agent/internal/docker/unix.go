package docker

import (
    "context"
    "fmt"
    "bytes"
    "encoding/json"
    "net"
    "net/http"
    "time"
)

// Minimal Docker Engine HTTP client over UNIX socket. // { SPECULATION }
type Client struct { sock string; http *http.Client }

func NewClient(sock string) *Client {
    tr := &http.Transport{
    DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
            // network and addr are ignored; connect to unix socket
            return (&net.Dialer{Timeout: 5 * time.Second}).DialContext(ctx, "unix", sock)
        },
    }
    return &Client{sock: sock, http: &http.Client{Transport: tr, Timeout: 10 * time.Second}}
}

func (c *Client) do(method, path string) (*http.Response, error) {
    req, err := http.NewRequest(method, "http://unix"+path, nil)
    if err != nil { return nil, err }
    return c.http.Do(req)
}

func (c *Client) doJSON(method, path string, v any) (*http.Response, error) {
    var body *bytes.Reader
    if v != nil {
        b, err := json.Marshal(v)
        if err != nil { return nil, err }
        body = bytes.NewReader(b)
    } else {
        body = bytes.NewReader(nil)
    }
    req, err := http.NewRequest(method, "http://unix"+path, body)
    if err != nil { return nil, err }
    req.Header.Set("Content-Type", "application/json")
    return c.http.Do(req)
}

// Kill sends a signal to a container by name or ID.
func (c *Client) Kill(container, signal string) error {
    if signal == "" { signal = "HUP" }
    resp, err := c.do("POST", fmt.Sprintf("/containers/%s/kill?signal=%s", container, signal))
    if err != nil { return err }
    defer resp.Body.Close()
    if resp.StatusCode >= 300 { return fmt.Errorf("docker kill failed: %s", resp.Status) }
    return nil
}

// Create creates a container with optional host port mapping.
func (c *Client) Create(name, image string, env map[string]string, containerPort, hostPort int) (string, error) {
    type hostBinding struct{ HostIp, HostPort string }
    body := map[string]any{
        "Image": image,
        "Env": toEnvArray(env),
    }
    if containerPort > 0 {
        portKey := fmt.Sprintf("%d/tcp", containerPort)
        body["ExposedPorts"] = map[string]any{ portKey: map[string]any{} }
        if hostPort > 0 {
            body["HostConfig"] = map[string]any{
                "PortBindings": map[string]any{
                    portKey: []hostBinding{{ HostIp: "0.0.0.0", HostPort: fmt.Sprintf("%d", hostPort) }},
                },
            }
        }
    }
    resp, err := c.doJSON("POST", "/containers/create?name="+name, body)
    if err != nil { return "", err }
    defer resp.Body.Close()
    if resp.StatusCode >= 300 { return "", fmt.Errorf("docker create failed: %s", resp.Status) }
    var out struct{ Id string `json:"Id"` }
    json.NewDecoder(resp.Body).Decode(&out)
    if out.Id == "" { out.Id = name }
    return out.Id, nil
}

func (c *Client) Start(id string) error {
    resp, err := c.do("POST", "/containers/"+id+"/start")
    if err != nil { return err }
    defer resp.Body.Close()
    if resp.StatusCode >= 300 { return fmt.Errorf("docker start failed: %s", resp.Status) }
    return nil
}

func (c *Client) Remove(nameOrID string, force bool) error {
    q := ""
    if force { q = "?force=1" }
    resp, err := c.do("DELETE", "/containers/"+nameOrID+q)
    if err != nil { return err }
    defer resp.Body.Close()
    if resp.StatusCode >= 300 && resp.StatusCode != 404 { return fmt.Errorf("docker remove failed: %s", resp.Status) }
    return nil
}

func toEnvArray(m map[string]string) []string {
    if len(m) == 0 { return nil }
    out := make([]string, 0, len(m))
    for k, v := range m { out = append(out, fmt.Sprintf("%s=%s", k, v)) }
    return out
}

// BuildRemote triggers a Docker build using a remote context URL (git/tar), tagging the image.
// docker API: POST /build?remote=<url>&t=<tag>&dockerfile=<path>
func (c *Client) BuildRemote(tag, remoteURL, dockerfile string) error {
    q := fmt.Sprintf("/build?remote=%s&t=%s", remoteURL, tag)
    if dockerfile != "" { q += "&dockerfile=" + dockerfile }
    resp, err := c.do("POST", q)
    if err != nil { return err }
    defer resp.Body.Close()
    if resp.StatusCode >= 300 { return fmt.Errorf("docker build failed: %s", resp.Status) }
    return nil
}
