package cache

import (
	"context"
	"errors"
	"github.com/phuslu/log"
	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte) error
}

type CacheType rune

const (
	RedisCacheType CacheType = 'r'
)

var (
	ErrInvalidCasting   = errors.New("invalid casting of redis client")
	ErrInvalidCacheType = errors.New("invalid cache type")
)

func New(cacheType CacheType, cacheParams any) (Cache, error) {
	switch cacheType {
	case RedisCacheType:
		rdb, ok := cacheParams.(*redis.Client)
		if !ok {
			log.Error().Err(ErrInvalidCasting).Msg("")
			return nil, ErrInvalidCasting
		}
		return newRedisCacheClient(rdb), nil
	default:
		return nil, ErrInvalidCacheType
	}
}
