package ports

import (
    "os"
    "path/filepath"
    "testing"
)

func TestAllocateAndFree(t *testing.T) {
    dir := t.TempDir()
    file := filepath.Join(dir, "ports.json")
    r, err := NewRegistry(file, 20000, 20002)
    if err != nil { t.Fatal(err) }

    p1, err := r.Allocate("a")
    if err != nil { t.Fatal(err) }
    if p1 < 20000 || p1 > 20002 { t.Fatalf("port out of range: %d", p1) }

    p2, err := r.Allocate("b")
    if err != nil { t.Fatal(err) }
    if p2 == p1 { t.Fatalf("collision: %d", p2) }

    // restart, ensure persistence
    r2, err := NewRegistry(file, 20000, 20002)
    if err != nil { t.Fatal(err) }
    if p, _ := r2.Allocate("a"); p != p1 { t.Fatalf("expected persisted %d got %d", p1, p) }

    // free and re-alloc reuse
    ok, err := r2.Free("a")
    if err != nil || !ok { t.Fatalf("free failed: %v %v", ok, err) }
    p3, err := r2.Allocate("c")
    if err != nil { t.Fatal(err) }
    // could be same as p1 but must be within range
    if p3 < 20000 || p3 > 20002 { t.Fatalf("port out of range: %d", p3) }

    // exhaust
    _, _ = r2.Allocate("d")
    if _, err := r2.Allocate("e"); err == nil { t.Fatalf("expected error when exhausted") }

    // cleanup file exists
    if _, err := os.Stat(file); err != nil { t.Fatalf("state not saved: %v", err) }
}
