package ports

import (
    "encoding/json"
    "errors"
    "fmt"
    "math/rand/v2"
    "os"
    "path/filepath"
    "sync"
)

type state struct {
    Alloc map[string]int `json:"alloc"`
}

type Registry struct {
    mu    sync.Mutex
    file  string
    start int
    end   int
    st    state
}

func NewRegistry(file string, start, end int) (*Registry, error) {
    r := &Registry{file: file, start: start, end: end, st: state{Alloc: map[string]int{}}}
    if err := r.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
        return nil, err
    }
    return r, nil
}

func (r *Registry) load() error {
    b, err := os.ReadFile(r.file)
    if err != nil { return err }
    return json.Unmarshal(b, &r.st)
}

func (r *Registry) save() error {
    if err := os.MkdirAll(filepath.Dir(r.file), 0o755); err != nil { return err }
    b, _ := json.MarshalIndent(r.st, "", "  ")
    return os.WriteFile(r.file, b, 0o644)
}

func (r *Registry) Allocate(projectID string) (int, error) {
    r.mu.Lock(); defer r.mu.Unlock()
    if p, ok := r.st.Alloc[projectID]; ok {
        return p, nil
    }
    // random start to reduce collisions across agents // { SPECULATION }
    start := r.start + rand.IntN(r.end-r.start+1)
    for i := 0; i <= (r.end-r.start); i++ {
        port := r.start + ((start - r.start + i) % (r.end - r.start + 1))
        if !r.inUse(port) {
            r.st.Alloc[projectID] = port
            if err := r.save(); err != nil { return 0, err }
            return port, nil
        }
    }
    return 0, fmt.Errorf("no free ports in range %d-%d", r.start, r.end)
}

func (r *Registry) inUse(p int) bool {
    for _, v := range r.st.Alloc { if v == p { return true } }
    return false
}

func (r *Registry) Free(projectID string) (bool, error) {
    r.mu.Lock(); defer r.mu.Unlock()
    if _, ok := r.st.Alloc[projectID]; ok {
        delete(r.st.Alloc, projectID)
        return true, r.save()
    }
    return false, nil
}
