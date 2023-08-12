package ajwt

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"time"
)

type RefreshTokenClaims struct {
	jwt.RegisteredClaims
	Salt string
}

func newRefreshTokenClaims(registeredClaims jwt.RegisteredClaims, salt string) *RefreshTokenClaims {
	return &RefreshTokenClaims{RegisteredClaims: registeredClaims, Salt: salt}
}

func (m *jwtManager) NewRefreshToken() (string, int64, int64, error) {
	salt, err := m.hasher.Hash(uuid.New().String())
	if err != nil {
		return "", -1, -1, err
	}

	issuedTime, expiresTime := jwt.NewNumericDate(time.Now()), jwt.NewNumericDate(time.Now().Add(m.refreshLiveTime))
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, newRefreshTokenClaims(
		jwt.RegisteredClaims{
			IssuedAt:  issuedTime,
			ExpiresAt: expiresTime,
			Subject:   "rtsaina",
		},
		salt,
	))

	token, err := jwtToken.SignedString(m.secretKey)
	if err != nil {
		return "", -1, -1, nil
	}

	return token, issuedTime.Unix(), expiresTime.Unix(), nil
}

func (m *jwtManager) ParseRefreshToken(ctx context.Context, token string) (*RefreshTokenClaims, error) {
	parsed, err := jwt.ParseWithClaims(token, &RefreshTokenClaims{}, m.validateMethod)
	if err != nil {
		return nil, err
	}
	if !parsed.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	rtClaims, ok := parsed.Claims.(*RefreshTokenClaims)
	if !ok {
		return nil, fmt.Errorf("wrong claims type")
	}

	err = rtClaims.Valid()
	if err != nil {
		return nil, err
	}

	return rtClaims, nil
}
