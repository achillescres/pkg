package redisCache

import (
	"context"
	"github.com/achillescres/pkg/cache"
	"github.com/redis/go-redis/v9"
	"time"
)

type redisCacher struct {
	client *redis.Client
}

func NewRedisCacher(client *redis.Client) cache.Cacher[[]byte] {
	return &redisCacher{client: client}
}

func (rC *redisCacher) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	err := rC.client.Set(
		ctx,
		key,
		value,
		expiration,
	).Err()
	if err != nil {
		return err
	}

	return nil
}

func (rC *redisCacher) Get(ctx context.Context, key string) ([]byte, bool) {
	val, err := rC.client.Get(
		ctx,
		key,
	).Bytes()
	if err != nil {
		return nil, false
	}

	if val == nil {
		return nil, false
	}

	return val, true
}

func (rC *redisCacher) Remove(ctx context.Context, key string) error {
	err := rC.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}
