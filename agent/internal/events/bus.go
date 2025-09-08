package events

import (
    "encoding/json"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type Event struct { OperationID string `json:"operation_id"`; Type string `json:"type"`; Payload any `json:"payload"`; Ts int64 `json:"ts"` }

type Bus struct { mu sync.Mutex; file string; mem []Event }

func NewBus(file string) *Bus { return &Bus{file: file} }

func (b *Bus) Emit(t string, op string, payload any) error {
    b.mu.Lock(); defer b.mu.Unlock()
    e := Event{OperationID: op, Type: t, Payload: payload, Ts: time.Now().UnixMilli()}
    b.mem = append(b.mem, e)
    if err := os.MkdirAll(filepath.Dir(b.file), 0o755); err != nil { return err }
    f, err := os.OpenFile(b.file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
    if err != nil { return err }
    defer f.Close()
    enc := json.NewEncoder(f)
    return enc.Encode(&e)
}

func (b *Bus) List() []Event { b.mu.Lock(); defer b.mu.Unlock(); return append([]Event{}, b.mem...) }
