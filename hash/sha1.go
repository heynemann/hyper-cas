package hash

import (
	"crypto/sha1"
)

type SHA1Hasher struct {
}

func NewSHA1Hasher() (*SHA1Hasher, error) {
	return &SHA1Hasher{}, nil
}

func (hasher *SHA1Hasher) Calc(value []byte) ([]byte, error) {
	h := sha1.New()
	h.Write(value)
	return h.Sum(nil), nil
}
