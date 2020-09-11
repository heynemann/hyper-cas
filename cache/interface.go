package cache

type CacheType int

const (
	LRU CacheType = iota
)

type Cache interface {
	Get(hash string) ([]byte, error)
	Set(hash string, value []byte) error

	GetDistro(hash string) (map[string]string, error)
	SetDistro(hash string, hashes map[string]string) error
}
