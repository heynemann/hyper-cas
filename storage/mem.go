package storage

import (
	"crypto/sha1"
	"fmt"
	"sync"
)

type MemStorage struct {
	sync.RWMutex
	storage map[string][]byte
}

func NewMemStorage() (*MemStorage, error) {
	return &MemStorage{storage: map[string][]byte{}}, nil
}

func (st *MemStorage) Store(value []byte) (string, error) {
	h := sha1.New()
	h.Write(value)
	hash := h.Sum(nil)
	st.Lock()
	defer st.Unlock()
	st.storage[string(hash)] = value
	return string(hash), nil
}

func (st *MemStorage) Get(hash string) ([]byte, error) {
	st.RLock()
	defer st.RUnlock()
	if val, ok := st.storage[hash]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("Hash %s was not found in storage.", hash)
}
