# rdctx

rdctx is a go-redis client wrapper with more features

### Example

```go
package main

import (
	"context"
	"log"

	"github.com/dnjooiopa/rdctx"
)

func main() {
	rc, ctx := rdctx.NewWithContext(context.Background(), "localhost:6379", "", 3)
	defer rc.Close()

	if err := rc.ConnectionOK(); err != nil {
		log.Println("cannot connect to redis:", err)
	}
	log.Println("redis connected")

	performTask(ctx)
}

func performTask(ctx context.Context) {
	if _, err := rdctx.Set(ctx, "foo", "bar"); err != nil {
		log.Println("cannot set key to redis:", err)
	}

	v, err := rdctx.Get(ctx, "foo")
	if err != nil {
		log.Println("cannot get key from redis:", err)
	}
	log.Println("result is", v)
}
```
