package ajwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type UserClaims struct {
	jwt.RegisteredClaims
	Login string
}

func newUserClaims(registeredClaims jwt.RegisteredClaims, login string) *UserClaims {
	return &UserClaims{RegisteredClaims: registeredClaims,
		Login: login,
	}
}

func (m *jwtManager) NewUser(login string) (string, error) {
	now := time.Now()

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, newUserClaims(
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.jwtLiveTime)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "saina auth service",
		},
		login,
	))

	token, err := jwtToken.SignedString(m.secretKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (m *jwtManager) ParseUser(token string) (*UserClaims, error) {
	parsed, err := jwt.ParseWithClaims(token, &UserClaims{}, m.validateMethod)

	if err != nil {
		return nil, err
	}
	if !parsed.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	uClaims, ok := parsed.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("wrong claims type")
	}

	err = uClaims.Valid()
	if err != nil {
		return nil, err
	}

	return uClaims, nil
}
