package hash

import (
	"crypto/sha1"
	"crypto/sha256"
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

type SHA256Hasher struct {
}

func NewSHA256Hasher() (*SHA256Hasher, error) {
	return &SHA256Hasher{}, nil
}

func (hasher *SHA256Hasher) Calc(value []byte) ([]byte, error) {
	h := sha256.New()
	h.Write(value)
	return h.Sum(nil), nil
}
