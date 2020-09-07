package content

import (
	"crypto"
	"fmt"
	"math"
)

const REPEATHASH string = "6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d"
const REPEATCONTENT string = "\x00"

type Node struct {
	Hash       string
	Content    []byte
	IsLeaf     bool
	IsRepeat   bool
	Parent     *Node
	LeftChild  *Node
	RightChild *Node
}

func NewNode(content []byte, isLeaf bool, hasher crypto.Hash) *Node {
	hash := hasher.New()
	hash.Write(content)
	return &Node{
		Content:  content,
		IsLeaf:   isLeaf,
		IsRepeat: false,
		Hash:     fmt.Sprintf("%x", hash.Sum(nil)),
	}
}

func NewRepeatNode() *Node {
	return &Node{
		Content:  []byte(REPEATCONTENT),
		IsLeaf:   true,
		IsRepeat: true,
		Hash:     REPEATHASH,
	}
}

type Tree struct {
	Nodes     []*Node
	Depth     int
	LeafCount int
	Hasher    crypto.Hash
}

func (t *Tree) Rebuild(data []byte, leafSize int, hasher crypto.Hash) error {
	size := math.Ceil(float64(len(data)) / float64(leafSize))
	depth := int(math.Log2(size)) + 1
	leafCount := int(math.Pow(2, float64(depth-1)))
	nodes := make([]*Node, int(leafCount*2-1))
	t.Nodes = nodes
	t.Depth = depth
	t.LeafCount = leafCount
	t.Hasher = hasher

	err := t.buildLeaves(data, leafSize, hasher)
	if err != nil {
		return err
	}
	err = t.buildParents()
	if err != nil {
		return err
	}
	return nil
}

func (t *Tree) Hash(data []byte) ([]byte, error) {
	h := t.Hasher.New()
	h.Write(data)
	return h.Sum(nil), nil
}

func (t *Tree) buildLeaves(data []byte, leafSize int, hasher crypto.Hash) error {
	for i := 0; i < t.LeafCount; i++ {
		if i*leafSize >= len(data) {
			t.Nodes[i] = NewRepeatNode()
			continue
		}

		t.Nodes[i] = NewNode(
			data[i*1024:(i+1)*1024], true, hasher,
		)
	}
	return nil
}

// 6 5 4 3 2 1 0
// 0 1 2 3 4 5 6
func (t *Tree) buildParents() error {
	for i := t.LeafCount; i < len(t.Nodes); i++ {
		depth := len(t.Nodes) - i - 1
		leftIndex := len(t.Nodes) - (depth*2 + 2) - 1
		left := t.Nodes[leftIndex]
		rightIndex := len(t.Nodes) - (depth*2 + 1) - 1
		right := t.Nodes[rightIndex]
		hash, err := t.Hash([]byte(left.Hash + right.Hash))
		if err != nil {
			return err
		}
		t.Nodes[i] = &Node{
			IsLeaf:   false,
			IsRepeat: false,
			Hash:     fmt.Sprintf("%x", hash),
			Content:  nil,
		}
	}
	return nil
}

func (t *Tree) Root() *Node {
	if len(t.Nodes) < 1 {
		return nil
	}
	return t.Nodes[len(t.Nodes)-1]
}

func (t *Tree) Leaves() []*Node {
	if len(t.Nodes) < 1 {
		return nil
	}
	return nil
}

func NewTree(data []byte, leafSize int, hasher crypto.Hash) (*Tree, error) {
	t := &Tree{}
	t.Rebuild(data, leafSize, hasher)
	return t, nil
}
