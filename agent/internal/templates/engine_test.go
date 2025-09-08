package templates

import (
    "os"
    "path/filepath"
    "testing"
)

func write(dir, file, content string) { _ = os.MkdirAll(filepath.Dir(filepath.Join(dir, file)), 0o755); _ = os.WriteFile(filepath.Join(dir, file), []byte(content), 0o644) }

func TestLoadDetectRender(t *testing.T) {
    root := t.TempDir()
    // static_html
    write(root, "static_html/template.yaml", "kind: Template/v1\nname: static_html\ntype: web\ndetect: []\nnginx:\n  enabled: true\n")
    write(root, "static_html/nginx_static.conf.tmpl", "server { listen 80; }")
    // nextjs_ssr
    write(root, "nextjs_ssr/template.yaml", "kind: Template/v1\nname: nextjs_ssr\ntype: web\ndetect:\n  - filesAny: [\"package.json\"]\nbuild:\n  mode: dockerfile\n  dockerfile: Dockerfile.tmpl\n")
    write(root, "nextjs_ssr/Dockerfile.tmpl", "FROM node:20\n")

    eng, err := Load(root)
    if err != nil { t.Fatal(err) }
    if len(eng.Manifests) != 2 { t.Fatalf("expected 2 manifests, got %d", len(eng.Manifests)) }

    proj := t.TempDir()
    write(proj, "package.json", "{\n \"dependencies\": {\n  \"next\": \"14\"\n }\n}")

    m, err := eng.Detect(proj, "")
    if err != nil { t.Fatal(err) }
    if m.Name != "nextjs_ssr" { t.Fatalf("expected nextjs_ssr, got %s", m.Name) }

    df, err := m.RenderDockerfile()
    if err != nil { t.Fatal(err) }
    if df == "" { t.Fatalf("empty dockerfile") }
}
