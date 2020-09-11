package hash

type HasherType int

const (
	SHA1 HasherType = iota
	SHA256
)

type Hasher interface {
	Calc(value []byte) ([]byte, error)
}
