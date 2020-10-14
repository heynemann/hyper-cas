package content

import (
	"github.com/heynemann/hyper-cas/utils"
	"math"
)

const REPEATHASH string = "6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d"
const REPEATCONTENT string = "\x00"

type NodeItem = struct {
	Key  string
	Hash []byte
}

type Node struct {
	Hash       []byte
	Content    []byte
	IsLeaf     bool
	IsRepeat   bool
	Parent     *Node
	LeftChild  *Node
	RightChild *Node
}

type Tree struct {
	Nodes     []*Node
	Depth     int
	Size      int
	LeafCount int
}

func (t *Tree) newNode(content []byte, isLeaf bool) *Node {
	hash := utils.HashBytes(content)
	return &Node{
		Content:  content,
		IsLeaf:   isLeaf,
		IsRepeat: false,
		Hash:     hash,
	}
}

func (t *Tree) newRepeatNode() *Node {
	hash := utils.Hash(REPEATCONTENT)
	return &Node{
		Content:  []byte(REPEATCONTENT),
		IsLeaf:   true,
		IsRepeat: true,
		Hash:     hash,
	}
}

func (t *Tree) calculateTreeDimensions(size int) (int, int, int) {
	nextnum := math.Ceil(math.Log2(float64(size)))
	leafCount := int(math.Pow(2.0, nextnum))
	nodeCount := leafCount*2 - 1
	depth := int(math.Log(float64(nodeCount))/math.Log(2)) + 1

	return depth, nodeCount, leafCount
}

func (t *Tree) initTree(depth, size, nodeCount, leafCount int) error {
	nodes := make([]*Node, nodeCount)
	t.Nodes = nodes
	t.Depth = depth
	t.LeafCount = leafCount
	t.Size = size
	return nil
}

func (t *Tree) rebuildFromData(data []byte, leafSize int) error {
	size := int(math.Ceil(float64(len(data)) / float64(leafSize)))
	depth, nodeCount, leafCount := t.calculateTreeDimensions(size)
	t.initTree(depth, size, nodeCount, leafCount)
	err := t.buildLeavesWithData(data, leafSize)
	if err != nil {
		return err
	}
	err = t.buildParents()
	if err != nil {
		return err
	}
	return nil
}

func (t *Tree) rebuildFromHashes(data []NodeItem) error {
	size := len(data)
	depth, nodeCount, leafCount := t.calculateTreeDimensions(size)
	t.initTree(depth, size, nodeCount, leafCount)
	err := t.buildLeavesWithMap(data)
	if err != nil {
		return err
	}
	err = t.buildParents()
	if err != nil {
		return err
	}

	return nil
}

func (t *Tree) buildLeavesWithData(data []byte, leafSize int) error {
	for i := 0; i < t.LeafCount; i++ {
		if i*leafSize >= len(data) {
			t.Nodes[i] = t.newRepeatNode()
			continue
		}

		begin := i * 1024
		end := (i + 1) * 1024
		if end > len(data) {
			end = len(data)
		}
		t.Nodes[i] = t.newNode(data[begin:end], true)
	}
	return nil
}

func (t *Tree) buildLeavesWithMap(data []NodeItem) error {
	itemCount := 0
	for _, item := range data {
		hash := utils.HashBytes(
			[]byte(item.Key),
			[]byte(":"),
			item.Hash,
		)
		t.Nodes[itemCount] = &Node{
			Hash:     hash,
			Content:  nil,
			IsLeaf:   true,
			IsRepeat: false,
		}
		itemCount += 1
	}
	for i := len(data); i < t.LeafCount; i++ {
		t.Nodes[i] = t.newRepeatNode()
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
		hash := utils.HashBytes(left.Hash, right.Hash)
		t.Nodes[i] = &Node{
			IsLeaf:   false,
			IsRepeat: false,
			Hash:     hash,
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
	return t.Nodes[0:t.Size]
}

func NewTreeWithData(data []byte, leafSize int) (*Tree, error) {
	t := &Tree{}
	t.rebuildFromData(data, leafSize)
	return t, nil
}

func NewTreeWithHashes(data []NodeItem) (*Tree, error) {
	t := &Tree{}
	t.rebuildFromHashes(data)
	return t, nil
}
