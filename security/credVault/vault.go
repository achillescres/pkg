package credVault

import (
	"context"
	"errors"
	"sync"
	"time"
)

// PullCred func pulls new cred and returns it and its expiration time
type PullCred[Cred any] func(ctx context.Context) (Cred, *time.Time, error)

// CredVault stores and updates credentials
// Uses PullCred function to pull new credentials
// Automatically updates credentials if expired or will in expirationTimeMinimumBound
type CredVault[Cred any] struct {
	mu                         *sync.RWMutex
	pull                       PullCred[Cred]
	expirationTimeMinimumBound time.Duration

	cred     Cred
	expireAt *time.Time
}

// New creates valid instance of CredVault[Cred]
func New[Cred any](ctx context.Context, pullFn PullCred[Cred], expirationTimeMinimumBound time.Duration) (*CredVault[Cred], error) {
	cred, expireAt, err := pullFn(ctx)
	if err != nil {
		return nil, err
	}

	return &CredVault[Cred]{
		pull:                       pullFn,
		expirationTimeMinimumBound: expirationTimeMinimumBound,

		cred:     cred,
		expireAt: expireAt,

		mu: &sync.RWMutex{},
	}, err
}

// maintain checks cred relevancy and updates it if cred is outdated
func (cv CredVault[Cred]) maintain(ctx context.Context) error {
	cv.mu.Lock()
	defer cv.mu.Unlock()

	if cv.expireAt == nil {
		return nil
	}
	if cv.expireAt.Sub(time.Now()) > cv.expirationTimeMinimumBound {
		return nil
	}
	cred, expireAt, err := cv.pull(ctx)
	if err != nil {
		return err
	}

	cv.cred = cred

	if expireAt == nil {
		cv.expireAt = nil
		return nil
	}

	if expireAt.Sub(time.Now()) < cv.expirationTimeMinimumBound {
		return errors.New("error new token is expired or don't satisfy expirationTimeMinimumBound")
	}
	cv.expireAt = expireAt
	return nil
}

// Cred updates credentials if it outdated and returns relevant credentials
func (cv CredVault[Cred]) Cred(ctx context.Context) (Cred, error) {
	if err := cv.maintain(ctx); err != nil {
		return *new(Cred), err
	}

	cv.mu.RLock()
	defer cv.mu.RUnlock()
	return cv.cred, nil
}

// ExpireAt updates credentials if it outdated and returns expire time
func (cv CredVault[Cred]) ExpireAt(ctx context.Context) (*time.Time, error) {
	if err := cv.maintain(ctx); err != nil {
		return nil, err
	}

	cv.mu.RLock()
	defer cv.mu.RUnlock()
	return cv.expireAt, nil
}
