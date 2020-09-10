package storage

type StorageType int

const (
	Memory StorageType = iota
	FileSystem
)

type Storage interface {
	Store(key string, value []byte) error
	Get(hash string) ([]byte, error)
	Has(hash string) bool
	StoreDistro(hash string, contents []string) error
	GetDistro(root string) ([]string, error)
	HasDistro(hash string) bool
}
