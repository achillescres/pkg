package passlib

import (
	"golang.org/x/crypto/bcrypt"
)

type HashManager interface {
	Hash(s string) (string, error)
	Compare(hashed, password string) error
}

type hashManager struct {
	salt string
}

var _ HashManager = (*hashManager)(nil)

func NewHashManager(salt string) HashManager {
	return &hashManager{salt: salt}
}

func (hM *hashManager) Salt(s string) string {
	return s + hM.salt
}

func (hM *hashManager) Hash(s string) (string, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(hM.Salt(s)), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(password), nil
}

func (hM *hashManager) Compare(hashed, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(hM.Salt(password)))
	if err != nil {
		return err
	}
	return nil
}
