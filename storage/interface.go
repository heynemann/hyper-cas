package storage

type StorageType int

const (
	Memory StorageType = iota
	FileSystem
)

type Storage interface {
	Store(key string, value []byte) error
	Get(hash string) ([]byte, error)
}
