package content

import (
	"crypto"
	"crypto/sha256"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func sha(content string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(content)))
}

func TestTreeWithData(t *testing.T) {
	var sb strings.Builder
	for i := 0; i < 2048; i++ {
		sb.WriteString(".")
	}
	content := sb.String()

	tree, err := NewTreeWithData([]byte(content), 1024, crypto.SHA256)

	assert.Nil(t, err)
	assert.NotNil(t, tree)
	assert.Equal(t, 3, len(tree.Nodes))
	assert.Equal(t, 2, tree.Depth)
	assert.NotNil(t, tree.Nodes[0])
	assert.NotNil(t, tree.Nodes[1])
	assert.NotNil(t, tree.Nodes[2])
	h := sha(tree.Nodes[0].Hash + tree.Nodes[1].Hash)
	assert.Equal(t, h, tree.Nodes[2].Hash)
}

func TestTreeWithDataIrregular(t *testing.T) {
	var sb strings.Builder
	for i := 0; i < 2300; i++ {
		sb.WriteString(".")
	}
	content := sb.String()

	tree, err := NewTreeWithData([]byte(content), 1024, crypto.SHA256)

	assert.Nil(t, err)
	assert.NotNil(t, tree)
	assert.Equal(t, 7, len(tree.Nodes))
	assert.Equal(t, 3, tree.Depth)
	assert.NotNil(t, tree.Nodes[0])
	assert.NotNil(t, tree.Nodes[1])
	h := sha(tree.Nodes[0].Hash + tree.Nodes[1].Hash)
	assert.Equal(t, h, tree.Nodes[4].Hash)
	h = sha(tree.Nodes[2].Hash + tree.Nodes[3].Hash)
	assert.Equal(t, h, tree.Nodes[5].Hash)
	h = sha(tree.Nodes[4].Hash + tree.Nodes[5].Hash)
	assert.Equal(t, h, tree.Nodes[6].Hash)
}

func TestTreeWithHashes(t *testing.T) {
	hash1 := fmt.Sprintf("%x", sha256.Sum256([]byte("test1")))
	hash2 := fmt.Sprintf("%x", sha256.Sum256([]byte("test2")))
	hash3 := fmt.Sprintf("%x", sha256.Sum256([]byte("test3")))
	hash4 := fmt.Sprintf("%x", sha256.Sum256([]byte("test4")))

	tree, err := NewTreeWithHashes(map[string]string{
		"test1": hash1,
		"test2": hash2,
		"test3": hash3,
		"test4": hash4,
	}, crypto.SHA256)

	assert.Nil(t, err)
	assert.NotNil(t, tree)
	assert.Equal(t, 7, len(tree.Nodes))
	assert.Equal(t, 3, tree.Depth)
	assert.Equal(t, tree.Nodes[0].Hash, hash1)
	assert.Equal(t, tree.Nodes[1].Hash, hash2)
	assert.Equal(t, tree.Nodes[2].Hash, hash3)
	assert.Equal(t, tree.Nodes[3].Hash, hash4)
	h := sha(tree.Nodes[0].Hash + tree.Nodes[1].Hash)
	assert.Equal(t, h, tree.Nodes[4].Hash)
	h = sha(tree.Nodes[2].Hash + tree.Nodes[3].Hash)
	assert.Equal(t, h, tree.Nodes[5].Hash)
	h = sha(tree.Nodes[4].Hash + tree.Nodes[5].Hash)
	assert.Equal(t, h, tree.Nodes[6].Hash)
}
