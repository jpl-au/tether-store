# Tether Store

Storage implementations for [Tether](https://github.com/jpl-au/tether)'s `SessionStore` and `DiffStore` interfaces. Each store is its own Go module so `go get` only pulls the dependencies you need.

## Available stores

| Store | Package | Install | Description |
|-------|---------|---------|-------------|
| Filesystem | [`fs`](fs/) | `go get github.com/jpl-au/tether-store/fs` | File-per-session on disk. Zero dependencies. Ideal for development and single-node deployments |
| Memory | [`memory`](memory/) | `go get github.com/jpl-au/tether-store/memory` | In-process map. No persistence. Ideal for testing and short-lived deployments |

## Quick start

```go
import "github.com/jpl-au/tether-store/fs"

tether.Stateful(app, tether.StatefulConfig[State]{
    SessionStore: fs.NewSessionStore("tmp/sessions"),
    DiffStore:    fs.NewDiffStore("tmp/diffs"),
    // ...
})
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

See the [SessionStore documentation](https://github.com/jpl-au/tether/blob/main/docs/session-store.md) and [DiffStore documentation](https://github.com/jpl-au/tether/blob/main/docs/store.md) for the full interface contracts.
