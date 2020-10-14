package storage

type StorageType int

const (
	FileSystem StorageType = iota
)

type Storage interface {
	Store(key string, value []byte) error
	Get(hash string) ([]byte, error)
	Has(hash string) bool

	StoreDistro(hash string, contents []string) error
	GetDistro(root string) ([]string, error)
	HasDistro(hash string) bool

	StoreLabel(hash string, label string) error
	GetLabel(label string) (string, error)
	HasLabel(label string) bool
}
