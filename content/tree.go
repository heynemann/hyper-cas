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
	Extra      map[string]string
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
		Extra:    map[string]string{},
		IsLeaf:   isLeaf,
		IsRepeat: false,
		Hash:     fmt.Sprintf("%x", hash.Sum(nil)),
	}
}

func NewRepeatNode() *Node {
	return &Node{
		Content:  []byte(REPEATCONTENT),
		Extra:    map[string]string{},
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

func (t *Tree) RebuildFromData(data []byte, leafSize int, hasher crypto.Hash) error {
	size := math.Ceil(float64(len(data)) / float64(leafSize))
	nextnum := math.Ceil(math.Log2(size))
	leafCount := int(math.Pow(2.0, nextnum))
	nodeCount := leafCount*2 - 1
	depth := int(math.Log(float64(nodeCount))/math.Log(2)) + 1
	nodes := make([]*Node, nodeCount)
	t.Nodes = nodes
	t.Depth = depth
	t.LeafCount = leafCount
	t.Hasher = hasher

	err := t.buildLeavesWithData(data, leafSize, hasher)
	if err != nil {
		return err
	}
	err = t.buildParents()
	if err != nil {
		return err
	}
	return nil
}

func (t *Tree) RebuildFromHashes(data map[string]string, hasher crypto.Hash) error {
	size := len(data)
	nextnum := math.Ceil(math.Log2(float64(size)))
	leafCount := int(math.Pow(2.0, nextnum))
	nodeCount := leafCount*2 - 1
	depth := int(math.Log(float64(nodeCount))/math.Log(2)) + 1
	nodes := make([]*Node, nodeCount)
	t.Nodes = nodes
	t.Depth = depth
	t.LeafCount = leafCount
	t.Hasher = hasher

	t.buildLeavesWithMap(data, hasher)
	t.buildParents()
	return nil
}

func (t *Tree) Hash(data []byte) ([]byte, error) {
	h := t.Hasher.New()
	h.Write(data)
	return h.Sum(nil), nil
}

func (t *Tree) buildLeavesWithData(data []byte, leafSize int, hasher crypto.Hash) error {
	for i := 0; i < t.LeafCount; i++ {
		if i*leafSize >= len(data) {
			t.Nodes[i] = NewRepeatNode()
			continue
		}

		begin := i * 1024
		end := (i + 1) * 1024
		if end > len(data) {
			end = len(data)
		}
		t.Nodes[i] = NewNode(
			data[begin:end], true, hasher,
		)
	}
	return nil
}

func (t *Tree) buildLeavesWithMap(data map[string]string, hasher crypto.Hash) error {
	itemCount := 0
	for path, item := range data {
		t.Nodes[itemCount] = &Node{
			Hash:     item,
			Content:  nil,
			IsLeaf:   true,
			IsRepeat: false,
			Extra:    map[string]string{"path": path},
		}
		itemCount += 1
	}
	for i := len(data); i < t.LeafCount; i++ {
		t.Nodes[i] = NewRepeatNode()
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

func NewTreeWithData(data []byte, leafSize int, hasher crypto.Hash) (*Tree, error) {
	t := &Tree{}
	t.RebuildFromData(data, leafSize, hasher)
	return t, nil
}

func NewTreeWithHashes(data map[string]string, hasher crypto.Hash) (*Tree, error) {
	t := &Tree{}
	t.RebuildFromHashes(data, hasher)
	return t, nil
}
