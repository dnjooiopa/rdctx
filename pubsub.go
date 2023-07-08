package rdctx

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Subscriber struct {
	*redis.PubSub
}

func (c *Subscriber) OnMessage(ctx context.Context, f func(message string)) {
	go func() {
		if _, err := c.ReceiveTimeout(ctx, -1); err != nil {
			return
		}
		for msg := range c.Channel() {
			f(msg.Payload)
		}
	}()
}
