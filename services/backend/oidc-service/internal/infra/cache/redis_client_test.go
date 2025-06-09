package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

func setupRedisContainerAndGetServiceURL(t *testing.T) (string, error) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	testcontainers.CleanupContainer(t, redisC)

	if err != nil {
		t.Errorf("expected no error, got %s", err)
		return "", err
	}

	return redisC.Endpoint(ctx, "")
}

func TestValidRedisClient(t *testing.T) {
	endpoint, err := setupRedisContainerAndGetServiceURL(t)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
		return
	}

	rdb := redis.NewClient(&redis.Options{Addr: endpoint})

	redisClientInstance := newRedisCacheClient(rdb)

	if redisClientInstance == nil {
		t.Errorf("expected valid instance, got nil")
		return
	}
}

func TestRedisClient_ValidWriteAndRead(t *testing.T) {
	endpoint, err := setupRedisContainerAndGetServiceURL(t)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
		return
	}

	rdb := redis.NewClient(&redis.Options{Addr: endpoint})

	redisClientInstance := newRedisCacheClient(rdb)

	ctx := context.Background()

	err = redisClientInstance.Set(ctx, "key", []byte("value"))
	if err != nil {
		t.Errorf("expected no error, got %s", err)
		return
	}

	value, err := redisClientInstance.Get(ctx, "key")
	if err != nil {
		t.Errorf("expected no error, got %s", err)
		return
	}

	if string(value) != "value" {
		t.Errorf("expected value, got %s", value)
	}
}
