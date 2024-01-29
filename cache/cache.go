package cache

import (
	"context"
	"time"
)

type Cacher[T any] interface {
	Set(ctx context.Context, key string, value T, expiration time.Duration) error
	Get(ctx context.Context, key string) (T, bool)
	Remove(ctx context.Context, key string) error
}
