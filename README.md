# rdctx

rdctx is a go-redis client wrapper with more features

## Examples

#### Simple usage

```go
package main

import (
 "context"
 "log"

 "github.com/dnjooiopa/rdctx"
)

func main() {
  host := "localhost:6379"
  pwd := ""
  db := 3

  // init
  client, ctx := rdctx.NewWithContext(context.Background(), host, pwd, db)
  defer client.Close()

  // Check if connection is ok
  if err := client.ConnOK(); err != nil {
    log.Println("cannot connect to redis:", err)
  }
  log.Println("redis connected")

  performTask(ctx)
}

func performTask(ctx context.Context) {
  // Set key-value
  _, err := rdctx.Set(ctx, "foo", "bar")
  if err != nil {
    log.Println("cannot set key to redis:", err)
  }

  // Get value from key
  v, err := rdctx.Get(ctx, "foo")
  if err != nil {
    log.Println("cannot get key from redis:", err)
  }
  log.Println("result is", v)
}
```

#### Inject to echo middleware

```go
e.Use(echo.WrapMiddleware(rdctx.Middleware(client)))
```

#### Set multiple key-value pairs with a single command

```go
pairs := []rdctx.KeyValue{
  {"foo1", "bar1"},
  {"foo2", "bar2"},
}
err := rdctx.MSetEx(ctx, pairs, 24*time.Hour)
```

## License

MIT
