package ttl

import (
    "time"
)

type Stopper interface { Stop(projectID string) error }
type PortFree interface { Free(projectID string) (bool, error) }
type Nginx interface { ConfPath(projectID string) string }

type Item struct { ProjectID string; ExpiresAt time.Time }

type Reaper struct { Interval time.Duration; Items func() []Item; Stop Stopper; Ports PortFree; Nginx Nginx }

func (r *Reaper) Run(stop <-chan struct{}) {
    t := time.NewTicker(r.Interval)
    defer t.Stop()
    for {
        select {
        case <-t.C:
            now := time.Now()
            for _, it := range r.Items() {
                if now.After(it.ExpiresAt) {
                    _ = r.Stop.Stop(it.ProjectID)
                    _, _ = r.Ports.Free(it.ProjectID)
                    // Removing nginx conf is handled by Stop // { SPECULATION }
                }
            }
        case <-stop:
            return
        }
    }
}
