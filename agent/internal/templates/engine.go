package templates

import (
    "errors"
    "fmt"
    "io/fs"
    "os"
    "path/filepath"
    "strings"

    yaml "gopkg.in/yaml.v3"
)

type DetectRule struct {
    FilesAny   []string `yaml:"filesAny"`
    JsonPathAll []struct{ Path string `yaml:"path"` } `yaml:"jsonPathAll"`
}

type Manifest struct {
    Kind     string `yaml:"kind"`
    Name     string `yaml:"name"`
    Type     string `yaml:"type"`
    Detect   []DetectRule `yaml:"detect"`
    Defaults struct {
        ContainerPort int `yaml:"containerPort"`
        HealthPath    string `yaml:"healthPath"`
        Env           []struct{ Name, Value string }
    } `yaml:"defaults"`
    Build struct {
        Mode       string `yaml:"mode"`
        Dockerfile string `yaml:"dockerfile"`
        Args       []struct{ Name, From string }
    } `yaml:"build"`
    Run struct {
        Command []string `yaml:"command"`
        ExposePortFrom string `yaml:"exposePortFrom"`
    } `yaml:"run"`
    Nginx struct {
        Enabled      bool   `yaml:"enabled"`
        ConfTemplate string `yaml:"confTemplate"`
        Websocket    bool   `yaml:"websocket"`
        CacheStatic  bool   `yaml:"cacheStatic"`
    } `yaml:"nginx"`
    Dir string `yaml:"-"` // filled at load
}

type Engine struct {
    Root string
    Manifests []Manifest
}

func Load(root string) (*Engine, error) {
    e := &Engine{Root: root}
    err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if err != nil { return err }
        if d.IsDir() { return nil }
        if filepath.Base(path) == "template.yaml" {
            b, err := os.ReadFile(path)
            if err != nil { return err }
            var m Manifest
            if err := yaml.Unmarshal(b, &m); err != nil { return fmt.Errorf("parse %s: %w", path, err) }
            m.Dir = filepath.Dir(path)
            e.Manifests = append(e.Manifests, m)
        }
        return nil
    })
    if err != nil { return nil, err }
    if len(e.Manifests) == 0 { return nil, errors.New("no templates found") }
    return e, nil
}

// Detect chooses a template by rules. Simplified: only checks filesAny existence.
func (e *Engine) Detect(projectDir string, override string) (*Manifest, error) {
    if override != "" {
        for i := range e.Manifests { if e.Manifests[i].Name == override { return &e.Manifests[i], nil } }
        return nil, fmt.Errorf("override template %s not found", override)
    }
    for i := range e.Manifests {
        m := &e.Manifests[i]
        for _, r := range m.Detect {
            if anyExists(projectDir, r.FilesAny) {
                return m, nil
            }
        }
    }
    // fallback to static_html
    for i := range e.Manifests { if e.Manifests[i].Name == "static_html" { return &e.Manifests[i], nil } }
    return nil, errors.New("no matching template and no fallback")
}

func anyExists(base string, files []string) bool {
    for _, f := range files {
        if _, err := os.Stat(filepath.Join(base, f)); err == nil { return true }
    }
    return false
}

// RenderDockerfile returns the raw template content with variable placeholders untouched.
// The agent will later replace {{ containerPort }} etc. Minimal for MVP. // { SPECULATION }
func (m *Manifest) RenderDockerfile() (string, error) {
    if m.Build.Mode != "dockerfile" || m.Build.Dockerfile == "" { return "", errors.New("no dockerfile mode") }
    path := filepath.Join(m.Dir, m.Build.Dockerfile)
    b, err := os.ReadFile(path)
    return string(b), err
}

// RenderNginxConf returns conf template if present.
func (m *Manifest) RenderNginxConf() (string, error) {
    if !m.Nginx.Enabled || m.Nginx.ConfTemplate == "" { return "", nil }
    b, err := os.ReadFile(filepath.Join(m.Dir, m.Nginx.ConfTemplate))
    return string(b), err
}
