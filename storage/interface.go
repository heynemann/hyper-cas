package storage

type StorageType int

const (
	Memory StorageType = iota
	FileSystem
)

type Storage interface {
	Store(value []byte) (string, error)
	Get(hash string) ([]byte, error)
}
