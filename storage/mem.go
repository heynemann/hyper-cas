package storage

import (
	"fmt"
	"sync"
)

type MemStorage struct {
	sync.RWMutex
	storage map[string][]byte
	distros map[string][]string
}

func NewMemStorage() (*MemStorage, error) {
	return &MemStorage{storage: map[string][]byte{}}, nil
}

func (st *MemStorage) Store(hash string, value []byte) error {
	st.Lock()
	defer st.Unlock()
	st.storage[hash] = value
	return nil
}

func (st *MemStorage) Get(hash string) ([]byte, error) {
	st.RLock()
	defer st.RUnlock()
	if val, ok := st.storage[hash]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("Hash %s was not found in storage.", hash)
}

func (st *MemStorage) Has(hash string) bool {
	st.RLock()
	defer st.RUnlock()
	if _, ok := st.storage[hash]; ok {
		return true
	}
	return false
}

func (st *MemStorage) StoreDistro(root string, hashes []string) error {
	st.Lock()
	defer st.Unlock()
	st.distros[root] = hashes
	return nil
}

func (st *MemStorage) GetDistro(root string) ([]string, error) {
	st.RLock()
	defer st.RUnlock()
	if val, ok := st.distros[root]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("Distro %s was not found in storage.", root)
}

func (st *MemStorage) HasDistro(hash string) bool {
	st.RLock()
	defer st.RUnlock()
	if _, ok := st.distros[hash]; ok {
		return true
	}
	return false
}
