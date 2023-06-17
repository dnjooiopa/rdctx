package rdctx

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func New(addr string, password string, db int) *Client {
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &Client{c}
}

type Client struct {
	*redis.Client
}

// connection ok if no error
func (c *Client) ConnectionOK() error {
	_, err := c.Ping(context.Background()).Result()
	return err
}

type ctxKeyClient struct{}

func NewWithContext(ctx context.Context, addr string, password string, db int) (*Client, context.Context) {
	c := New(addr, password, db)
	return c, NewContext(ctx, c)
}

func NewContext(ctx context.Context, c *Client) context.Context {
	return context.WithValue(ctx, ctxKeyClient{}, c)
}

func Middleware(c *Client) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(NewContext(r.Context(), c))
			h.ServeHTTP(w, r)
		})
	}
}

func c(ctx context.Context) *Client {
	return ctx.Value(ctxKeyClient{}).(*Client)
}

var keyPrefix string

func SetKeyPrefix(prefix string) {
	keyPrefix = prefix
}

func addPrefix(keys ...string) []string {
	if keyPrefix == "" {
		return keys
	}
	for i, v := range keys {
		keys[i] = keyPrefix + ":" + v
	}
	return keys
}

func replacePrefix(keys ...string) []string {
	if keyPrefix == "" {
		return keys
	}
	for i, v := range keys {
		if strings.HasPrefix(v, keyPrefix+":") {
			keys[i] = strings.TrimPrefix(v, keyPrefix+":")
		}
	}
	return keys
}

func SetEx(ctx context.Context, key string, value interface{}, exp time.Duration) (string, error) {
	return c(ctx).Set(ctx, addPrefix(key)[0], value, exp).Result()
}

func Set(ctx context.Context, key string, value interface{}) (string, error) {
	return SetEx(ctx, key, value, 0)
}

func Del(ctx context.Context, keys ...string) (int64, error) {
	return c(ctx).Del(ctx, addPrefix(keys...)...).Result()
}

func Incr(ctx context.Context, key string) (int64, error) {
	return c(ctx).Incr(ctx, addPrefix(key)[0]).Result()
}

func Get(ctx context.Context, key string) (string, error) {
	return c(ctx).Get(ctx, addPrefix(key)[0]).Result()
}

func Expire(ctx context.Context, key string, exp time.Duration) (bool, error) {
	return c(ctx).Expire(ctx, addPrefix(key)[0], exp).Result()
}

func Keys(ctx context.Context, pattern string) ([]string, error) {
	keys, err := c(ctx).Keys(ctx, addPrefix(pattern)[0]).Result()
	if err != nil {
		return nil, err
	}
	return replacePrefix(keys...), nil
}

func MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	return c(ctx).MGet(ctx, addPrefix(keys...)...).Result()
}

func Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return c(ctx).Scan(ctx, cursor, addPrefix(match)[0], count).Result()
}

type KeyValue struct {
	Key   string
	Value interface{}
}

func MSetEx(ctx context.Context, keyValues []KeyValue, exp time.Duration) error {
	if len(keyValues) == 0 {
		return nil
	}

	pipeline := c(ctx).Pipeline()
	for _, v := range keyValues {
		pipeline.Set(ctx, addPrefix(v.Key)[0], v.Value, exp)
	}
	_, err := pipeline.Exec(ctx)
	return err
}
