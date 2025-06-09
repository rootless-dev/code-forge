package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	rdb *redis.Client
}

func newRedisCacheClient(rdb *redis.Client) Cache {
	return &RedisClient{
		rdb: rdb,
	}
}

func (c *RedisClient) Get(ctx context.Context, key string) ([]byte, error) {
	return c.rdb.Get(ctx, key).Bytes()
}

func (c *RedisClient) Set(ctx context.Context, key string, value []byte) error {
	return c.rdb.Set(ctx, key, value, 0).Err()
}
