package hash

type HasherType int

const (
	SHA1 HasherType = iota
)

type Hasher interface {
	Calc(value []byte) ([]byte, error)
}
