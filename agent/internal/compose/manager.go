package compose

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"
)

type Service struct {
    Name string
    Image string
    Env map[string]string
    Ports map[int]int // host:container
}

type Manager struct { Dir string }

func (m *Manager) file(projectID string) string { return filepath.Join(m.Dir, fmt.Sprintf("%s.compose.yaml", projectID)) }

func (m *Manager) Write(projectID string, s Service) (string, error) {
    path := m.file(projectID)
    if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil { return "", err }
    // very minimal yaml // { SPECULATION }
    b := &strings.Builder{}
    b.WriteString("version: '3.9'\nservices:\n  app:\n    image: ")
    b.WriteString(s.Image)
    b.WriteString("\n")
    if len(s.Env) > 0 {
        b.WriteString("    environment:\n")
        for k, v := range s.Env { b.WriteString(fmt.Sprintf("      %s: %q\n", k, v)) }
    }
    if len(s.Ports) > 0 {
        b.WriteString("    ports:\n")
        for host, cont := range s.Ports { b.WriteString(fmt.Sprintf("      - \"%d:%d\"\n", host, cont)) }
    }
    content := b.String()
    return path, os.WriteFile(path, []byte(content), 0o644)
}

// Up runs docker compose up -d on the generated file. // { SPECULATION }
func (m *Manager) Up(projectID string) error {
    // TODO: shell out to docker compose or use Engine API. Placeholder no-op. // { SPECULATION }
    log.Printf("compose up requested for %s (noop)", projectID)
    return nil
}

// Down runs docker compose down on the generated file.
func (m *Manager) Down(projectID string) error {
    log.Printf("compose down requested for %s (noop)", projectID)
    return nil
}
