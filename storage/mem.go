package storage

import (
	"crypto/sha1"
	"fmt"
	"io"
	"sync"
)

type MemStorage struct {
	sync.RWMutex
	storage map[string]string
}

func NewMemStorage() (*MemStorage, error) {
	return &MemStorage{storage: map[string]string{}}, nil
}

func (st *MemStorage) Store(value string) (string, error) {
	h := sha1.New()
	io.WriteString(h, value)
	hash := fmt.Sprintf("%x", h.Sum(nil))
	st.Lock()
	defer st.Unlock()
	st.storage[hash] = value
	return hash, nil
}

func (st *MemStorage) Get(hash string) (string, error) {
	st.RLock()
	defer st.RUnlock()
	if val, ok := st.storage[hash]; ok {
		return val, nil
	}
	return "", fmt.Errorf("Hash %s was not found in storage.", hash)
}
