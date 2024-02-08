package hash

import (
	"crypto/md5"
	"encoding/hex"
)

type InstantHasher func(...[]byte) string

func NewMD5() InstantHasher {
	return func(elems ...[]byte) string {
		hash := md5.New()

		for _, elem := range elems {
			hash.Write(elem)
		}

		return hex.EncodeToString(hash.Sum(nil))
	}
}
