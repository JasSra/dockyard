package ops

import "sync"

type Runner struct { done sync.Map }

// Execute a whitelisted action if idempotencyKey not seen.
func (r *Runner) Do(key string, fn func() error) error {
    if _, loaded := r.done.LoadOrStore(key, struct{}{}); loaded { return nil }
    return fn()
}
