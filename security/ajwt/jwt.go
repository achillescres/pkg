package ajwt

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Hasher interface {
	Hash(s string) (string, error)
}

type JWTManager interface {
	validateMethod(token *jwt.Token) (any, error)
	NewUser(login string) (string, error)
	ParseUser(token string) (*UserClaims, error)
	NewRefreshToken() (string, int64, int64, error)
	NewTokenPair(login string) (*TokenPair, error)
	ParseRefreshToken(ctx context.Context, token string) (*RefreshTokenClaims, error)
}

type TokenPair struct {
	JWT          string `json:"JWT" binding:"required"`
	RefreshToken string `json:"refreshToken" binding:"required"`
}

func NewTokenPair(JWT string, refreshToken string) *TokenPair {
	return &TokenPair{JWT: JWT, RefreshToken: refreshToken}
}

type jwtManager struct {
	hasher          Hasher
	secretKey       []byte
	jwtLiveTime     time.Duration
	refreshLiveTime time.Duration
	keyFuncFabric   func(secret []byte) jwt.Keyfunc
}

func NewJWTManager(hasher Hasher, secretKey string, jwtLiveTime time.Duration, refreshLiveTime time.Duration) JWTManager {
	return &jwtManager{
		hasher:          hasher,
		secretKey:       []byte(secretKey),
		jwtLiveTime:     jwtLiveTime,
		refreshLiveTime: refreshLiveTime,
	}
}

func (m *jwtManager) NewTokenPair(login string) (*TokenPair, error) {
	user, err := m.NewUser(login)
	if err != nil {
		return nil, err
	}
	rt, _, _, err := m.NewRefreshToken()
	if err != nil {
		return nil, err
	}
	return NewTokenPair(user, rt), nil
}

func (m *jwtManager) validateMethod(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("wrong signing method")
	}
	return m.secretKey, nil
}
