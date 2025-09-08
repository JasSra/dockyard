package nginx

import (
    "fmt"
    "os"
    "path/filepath"
)

type Manager struct {
    ConfDir string
    UpstreamHost string // defaults to 127.0.0.1 if empty
}

func (m *Manager) ConfPath(projectID string) string {
    return filepath.Join(m.ConfDir, fmt.Sprintf("%s.conf", projectID))
}

func (m *Manager) WriteServer(projectID, host string, upstreamPort int) (string, error) {
    path := m.ConfPath(projectID)
    if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil { return "", err }
        up := m.UpstreamHost
        if up == "" { up = "127.0.0.1" }
        conf := fmt.Sprintf(`server {
    listen 80;
    server_name %s;
    gzip on; # { SPECULATION }
    location / {
            proxy_pass http://%s:%d;
      proxy_http_version 1.1;
      proxy_set_header Host $host;
      proxy_set_header X-Forwarded-For $remote_addr;
      proxy_set_header Upgrade $http_upgrade;
      proxy_set_header Connection $connection_upgrade;
      proxy_read_timeout 60s; # { SPECULATION }
    }
}
`, host, up, upstreamPort)
    if err := os.WriteFile(path, []byte(conf), 0o644); err != nil { return "", err }
    return path, nil
}

// Reload performs `docker exec nginx nginx -s reload` outside of this package.
// Here we just expose a hook point for tests. // { SPECULATION }
type Reloader interface { Reload() error }
