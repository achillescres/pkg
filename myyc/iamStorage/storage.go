package iamStorage

import (
	"context"
	"errors"
	"saina.gitlab.yandexcloud.net/saina/backend/pkg/myyc"
	"time"
)

const (
	tokenExpireAdditionalBarrier = time.Minute * 5
	tokenUpdateTimeout           = time.Second * 10
)

// PullIAM must return new valid myyc.IAM and nil error or empty myyc.IAM and non-nil error
// must stop working if context is closed
type PullIAM func(context.Context) (myyc.IAM, error)

type IAMStorage struct {
	pull PullIAM

	iam      *myyc.IAM
	expireAt time.Time
}

func New(pull PullIAM) *IAMStorage {
	return &IAMStorage{pull: pull}
}

func (s *IAMStorage) updateToken(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, tokenUpdateTimeout)
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()

	iam, err := s.pull(ctx)
	if err != nil {
		return err
	}
	if len(iam.AccessToken) == 0 {
		return errors.New("error iam token is empty")
	}
	if iam.TokenType != "Bearer" {
		return errors.New("error iam token type isn't Bearer")
	}

	s.iam = &iam
	s.expireAt = time.Now().Add(time.Duration(s.iam.ExpiresIn) * time.Second).Add(tokenExpireAdditionalBarrier)
	return nil
}

func (s *IAMStorage) maintainToken(ctx context.Context) error {
	if s.iam == nil {
		return s.updateToken(ctx)
	}

	// Is token relevant enough?
	if last := s.expireAt.Sub(time.Now()); last > tokenExpireAdditionalBarrier {
		// then just chill
		return nil
	}

	return s.updateToken(ctx)
}

func (s *IAMStorage) Token(ctx context.Context) (myyc.IAM, error) {
	err := s.maintainToken(ctx)
	if err != nil {
		return myyc.IAM{}, err
	}

	return *s.iam, nil
}

func (s *IAMStorage) ExpireAt(ctx context.Context) (time.Time, error) {
	err := s.maintainToken(ctx)
	if err != nil {
		return time.Time{}, err
	}

	return s.expireAt, nil
}
