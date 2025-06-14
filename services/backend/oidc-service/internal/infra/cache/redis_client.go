package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisClient struct {
	rdb *redis.Client
}

// newRedisCacheClient initializes a new RedisClient instance that implements the Cache interface.
//
// Parameters:
//   - rdb: A pointer to an already initialized Redis client instance.
//
// Returns:
//   - A Cache interface implemented by the RedisClient.
func newRedisCacheClient(rdb *redis.Client) Cache {
	return &RedisClient{
		rdb: rdb,
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

// Set stores the given value in Redis under the specified key with a custom expiration duration.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control.
//   - key: The key under which the value is stored.
//   - value: The value to store as a byte slice.
//   - expirationTime: Duration for which the key should remain in Redis before expiring.
//
// Returns:
//   - An error if the operation fails.
func (c *RedisClient) Set(ctx context.Context, key string, value []byte, expirationTime time.Duration) error {
	return c.rdb.Set(ctx, key, value, expirationTime).Err()
}
