package hash

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
)

//type InstantHasher func(...[]byte) string

type InstantHasher interface {
	MD5(elems ...[]byte) string
	SHA512(elems ...[]byte) string
}

func MD5(elems ...[]byte) string {
	hash := md5.New()

	for _, elem := range elems {
		hash.Write(elem)
	}

	return hex.EncodeToString(hash.Sum(nil))
}

func SHA512(elems ...[]byte) string {
	hash := sha512.New()

	for _, elem := range elems {
		hash.Write(elem)
	}

	return hex.EncodeToString(hash.Sum(nil))
}
