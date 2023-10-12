package mysync

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
)

const MaxPoolSize = 1048576

var ErrExceededMaxPoolSize = errors.New("exceeded max pool size")

// StartPool maps f to each element of ins in errgroup with limit of poolSize returns err of first function
func StartPool[T any](ctx context.Context, poolSize int, ins []T, f func(context.Context, T) error) error {
	if poolSize > MaxPoolSize {
		return ErrExceededMaxPoolSize
	}

	grp, ctx := errgroup.WithContext(ctx)
	grp.SetLimit(poolSize)
	for _, a := range ins {
		_a := a
		grp.Go(func() error {
			return f(ctx, _a)
		})
	}
	err := grp.Wait()
	if err != nil {
		return fmt.Errorf("caught in pool: %w", err)
	}
	return nil
}
