package rdctx

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type subscriber struct {
	*redis.PubSub
}

func (c *subscriber) OnMessage(ctx context.Context, f func(message string, err error)) {
	go func() {
		for {
			msg, err := c.ReceiveMessage(ctx)
			if err != nil {
				f("", err)
				continue
			}
			f(msg.Payload, nil)
		}
	}()
}
