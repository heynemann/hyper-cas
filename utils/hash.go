package utils

import (
	"crypto"
)

func HashBytes(content ...[]byte) []byte {
	h := crypto.SHA1.New()
	for _, d := range content {
		h.Write(d)
	}
	return h.Sum(nil)
}

func Hash(content string) []byte {
	return HashBytes([]byte(content))
}
