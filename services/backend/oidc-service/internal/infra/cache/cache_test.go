package cache

import (
	"github.com/redis/go-redis/v9"
	"testing"
)

func TestNew_ValidClient(t *testing.T) {
	rdb := &redis.Client{}
	validInstance, err := New(RedisCacheType, rdb)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
		return
	}

	if validInstance == nil {
		t.Errorf("expected valid instance, got nil")
		return
	}

}

func TestNew_InvalidCastingRDB(t *testing.T) {
	invalidInstance, err := New(RedisCacheType, "invalid")
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	if invalidInstance != nil {
		t.Errorf("expected nil, got %v", invalidInstance)
		return
	}
}

func TestNew_InvalidCacheType(t *testing.T) {
	rdb := &redis.Client{}

	_, err := New('z', rdb)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}
