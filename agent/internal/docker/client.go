package docker

// Minimal interfaces and a no-op client for MVP. // { SPECULATION }

type BuildOptions struct { ContextDir string; Dockerfile string; Tag string; Args map[string]string }
type RunOptions struct { Name string; Image string; Env map[string]string; Ports map[int]int }

type Status struct { Running bool; ContainerID string }
type Stats struct { CPU float64; MemBytes uint64 }

type Client interface {
    Build(opts BuildOptions) error
    Run(opts RunOptions) error
    Stop(name string) error
    Remove(name string) error
    Logs(name string, tail int) ([]string, error)
    Status(name string) (Status, error)
    PruneImages(olderThanHours int) error
    ContainerStats(name string) (Stats, error)
}

type NoopClient struct{}

func (NoopClient) Build(opts BuildOptions) error                       { return nil }
func (NoopClient) Run(opts RunOptions) error                           { return nil }
func (NoopClient) Stop(name string) error                              { return nil }
func (NoopClient) Remove(name string) error                            { return nil }
func (NoopClient) Logs(name string, tail int) ([]string, error)        { return []string{"noop"}, nil }
func (NoopClient) Status(name string) (Status, error)                  { return Status{Running: false, ContainerID: ""}, nil }
func (NoopClient) PruneImages(olderThanHours int) error                { return nil }
func (NoopClient) ContainerStats(name string) (Stats, error)           { return Stats{}, nil }
