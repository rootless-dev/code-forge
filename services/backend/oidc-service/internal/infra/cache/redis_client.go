package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisClient struct {
	rdb            *redis.Client
	expirationTime time.Duration
}

// newRedisCacheClient creates a new Redis cache client with an optional expiration time.
// If no expiration time is provided, it defaults to 5 minutes.
//
// Parameters:
//   - rdb: A pointer to a Redis client.
//   - expirationTime: Optional expiration time for cache entries.
//
// Returns:
//   - A Cache interface backed by Redis.
func newRedisCacheClient(rdb *redis.Client, expirationTime ...time.Duration) Cache {
	expiration := 5 * time.Minute

	if len(expirationTime) > 0 {
		expiration = expirationTime[0]
	}

	return &RedisClient{
		rdb:            rdb,
		expirationTime: expiration,
	}
}

// Get retrieves the value associated with the given key from Redis.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control.
//   - key: The key to retrieve the value for.
//
// Returns:
//   - The value as a byte slice, or an error if the key does not exist or retrieval fails.
func (c *RedisClient) Get(ctx context.Context, key string) ([]byte, error) {
	return c.rdb.Get(ctx, key).Bytes()
}

// Set stores the given value in Redis with the configured expiration time.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control.
//   - key: The key under which the value is stored.
//   - value: The value to store as a byte slice.
//
// Returns:
//   - An error if the operation fails.
func (c *RedisClient) Set(ctx context.Context, key string, value []byte) error {
	return c.rdb.Set(ctx, key, value, c.expirationTime).Err()
}
