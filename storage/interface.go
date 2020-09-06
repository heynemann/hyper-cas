package storage

type Storage interface {
	Store(value string) (string, error)
	Get(hash string) (string, error)
}
