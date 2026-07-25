package cache

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"github.com/furkancmn57/go-architecture-template/src/config"
	"github.com/furkancmn57/go-architecture-template/src/interfaces"
)

type redisCache struct {
	client *goredis.Client
}

// NewClient creates and pings a Redis client.
func NewClient(cfg config.Redis) (*goredis.Client, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("cache: failed to connect to redis: %w", err)
	}
	return client, nil
}

// New returns an interfaces.Cache backed by Redis.
func New(client *goredis.Client) interfaces.Cache {
	return &redisCache{client: client}
}

func (c *redisCache) Get(ctx context.Context, key string) (string, error) {
	value, err := c.client.Get(ctx, key).Result()
	if err == goredis.Nil {
		return "", nil
	}
	return value, err
}

func (c *redisCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *redisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
