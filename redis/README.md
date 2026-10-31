# redis

Redis Pub/Sub backend for [Tether](https://github.com/jpl-au/tether)'s
Cluster interface. Enables cross-node communication for Bus and Value
primitives.

## Install

```bash
go get github.com/jpl-au/tether-adapter/redis
```

## Usage

The package is named `redis`, as go-redis's is, so import go-redis under another name.

```go
import (
    goredis "github.com/redis/go-redis/v9"
    "github.com/jpl-au/tether"
    "github.com/jpl-au/tether-adapter/redis"
)

func main() {
    rdb := goredis.NewClient(&goredis.Options{Addr: "localhost:6379"})
    defer rdb.Close()

    app := tether.App{
        Cluster: redis.New(rdb),
    }

    // Buses and Values with topic names now publish/subscribe
    // across all nodes connected to the same Redis.
}
```

## How it works

- `Publish` delegates to Redis `PUBLISH`
- `Subscribe` creates a Redis Pub/Sub subscription, confirms it
  with `Receive` before returning (so the caller can publish
  immediately without a race), and runs a goroutine to read
  messages from the channel
- The returned unsubscribe function closes the Redis subscription

The caller owns the `*goredis.Client` lifetime. Create it before
the Cluster, close it after the tether application shuts down.

## Testing

Tests use [miniredis](https://github.com/alicebob/miniredis) - a
pure Go in-process Redis server. No external Redis required.

```bash
go test ./...
```

## Licence

MIT
