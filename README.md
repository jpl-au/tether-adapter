# Tether Adapter

Adapters that connect [Tether](https://github.com/jpl-au/tether) to something outside it: a store for its sessions and diffs, a cluster for its buses and values, or a desktop app that hosts its pages. Each adapter is its own Go module so `go get` only pulls the dependencies you need.

## Available adapters

| Adapter | Package | Install | Description |
|---------|---------|---------|-------------|
| Filesystem | [`fs`](fs/) | `go get github.com/jpl-au/tether-adapter/fs` | `SessionStore` and `DiffStore` as a file per session on disk. Zero dependencies. Ideal for development and single-node deployments |
| Memory | [`memory`](memory/) | `go get github.com/jpl-au/tether-adapter/memory` | `SessionStore` and `DiffStore` as an in-process map. No persistence. Ideal for testing and short-lived deployments |
| Redis | [`redis`](redis/) | `go get github.com/jpl-au/tether-adapter/redis` | `Cluster` over Redis Pub/Sub, so buses and values reach every node |
| Wails | [`wails`](wails/) | `go get github.com/jpl-au/tether-adapter/wails` | `mode.IPC` inside a Wails v3 window, over the window's own stream, with no network port |

## Quick start

```go
import "github.com/jpl-au/tether-adapter/fs"

tether.Stateful(app, tether.StatefulConfig[State]{
    SessionStore: fs.NewSessionStore("tmp/sessions"),
    DiffStore:    fs.NewDiffStore("tmp/diffs"),
    // ...
})
```

The Redis adapter's package is named `redis`, as go-redis's is, so import go-redis under another name:

```go
import (
    goredis "github.com/redis/go-redis/v9"
    "github.com/jpl-au/tether-adapter/redis"
)

rdb := goredis.NewClient(&goredis.Options{Addr: "localhost:6379"})
app := tether.App{Cluster: redis.New(rdb)}
```

## Writing your own

The `SessionStore` and `DiffStore` interfaces are intentionally simple - three methods each. Implement them with your preferred backend:

```go
type SessionStore interface {
    Save(ctx context.Context, id string, data []byte, ttl time.Duration) error
    Load(ctx context.Context, id string) ([]byte, error)
    Delete(ctx context.Context, id string) error
}

type DiffStore interface {
    Save(ctx context.Context, id string, data []byte) error
    Load(ctx context.Context, id string) ([]byte, error)
    Delete(ctx context.Context, id string) error
}
```

See the [SessionStore documentation](https://github.com/jpl-au/tether/blob/main/docs/session-store.md), [DiffStore documentation](https://github.com/jpl-au/tether/blob/main/docs/store.md) and [Cluster documentation](https://github.com/jpl-au/tether/blob/main/docs/cluster.md) for the full interface contracts.
