package hash

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
)

//type InstantHasher func(...[]byte) string

type InstantHasher interface {
	Hash(elems ...[]byte) string
}

type md5H struct{}

func NewMD5H() InstantHasher {
	return &md5H{}
}

func (md5H) Hash(elems ...[]byte) string {
	hash := md5.New()

	for _, elem := range elems {
		hash.Write(elem)
	}

	return hex.EncodeToString(hash.Sum(nil))
}

type sha512H struct{}

func NewSHA512H() InstantHasher {
	return &sha512H{}
}

func (sha512H) Hash(elems ...[]byte) string {
	hash := sha512.New()

	for _, elem := range elems {
		hash.Write(elem)
	}

	return hex.EncodeToString(hash.Sum(nil))
}
