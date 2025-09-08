package nginx

import (
    "os"
    "path/filepath"
    "strings"
    "testing"
)

func TestWriteServer(t *testing.T) {
    dir := t.TempDir()
    m := &Manager{ConfDir: dir}
    path, err := m.WriteServer("proj1", "abc.example.com", 23456)
    if err != nil { t.Fatal(err) }
    if !strings.HasSuffix(path, filepath.Join(dir, "proj1.conf")) { t.Fatalf("unexpected path: %s", path) }
    b, err := os.ReadFile(path)
    if err != nil { t.Fatal(err) }
    s := string(b)
    if !strings.Contains(s, "server_name abc.example.com") { t.Fatalf("missing host: %s", s) }
    if !strings.Contains(s, "127.0.0.1:23456") { t.Fatalf("missing upstream: %s", s) }
}
