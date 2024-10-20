package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type AbstractCacheService interface {
	SetItem(key string, item interface{}, expire time.Duration) error
	GetItem(key string) (any, error)
}

type CacheService struct {
	redis redis.Client
	ctx   context.Context
}

func NewCacheService(redis redis.Client, ctx context.Context) AbstractCacheService {
	return &CacheService{
		redis: redis,
		ctx:   ctx,
	}
}

func (cacheService *CacheService) SetItem(key string, item interface{}, expire time.Duration) error {
	err := cacheService.redis.Set(cacheService.ctx, key, item, expire).Err()
	if err != nil {
		panic(err)
	}

	return nil
}

func (cacheService *CacheService) GetItem(key string) (any, error) {
	item, err := cacheService.redis.Get(cacheService.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return item, nil
}
