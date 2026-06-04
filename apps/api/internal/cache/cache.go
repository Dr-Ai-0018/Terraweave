package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

type Status struct {
	OK    bool   `json:"ok"`
	Ping  string `json:"ping,omitempty"`
	Error string `json:"error,omitempty"`
}

func Open(redisURL string) (*Cache, error) {
	if redisURL == "" {
		return nil, fmt.Errorf("missing TERRAWEAVE_REDIS_URL")
	}
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	return &Cache{client: redis.NewClient(opts)}, nil
}

func (c *Cache) Close() error {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Close()
}

func (c *Cache) Check(ctx context.Context) Status {
	if c == nil || c.client == nil {
		return Status{OK: false, Error: "redis not configured"}
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	pong, err := c.client.Ping(ctx).Result()
	if err != nil {
		return Status{OK: false, Error: err.Error()}
	}
	return Status{OK: true, Ping: pong}
}
