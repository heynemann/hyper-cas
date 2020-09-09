package content

import (
	"crypto"
	"crypto/sha256"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func sha(content ...[]byte) []byte {
	h := sha256.New()
	for _, hash := range content {
		h.Write(hash)
	}
	return h.Sum(nil)
}

func fileHash(name string, hash []byte) []byte {
	return sha([]byte(name), []byte(":"), hash)
}

func assertHash(t *testing.T, h1 []byte, h2 []byte) {
	assert.Equal(t, fmt.Sprintf("%x", h1), fmt.Sprintf("%x", h2))
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
	h := sha(tree.Nodes[0].Hash, tree.Nodes[1].Hash)
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
	h := sha(tree.Nodes[0].Hash, tree.Nodes[1].Hash)
	assertHash(t, h, tree.Nodes[4].Hash)
	h = sha(tree.Nodes[2].Hash, tree.Nodes[3].Hash)
	assertHash(t, h, tree.Nodes[5].Hash)
	h = sha(tree.Nodes[4].Hash, tree.Nodes[5].Hash)
	assertHash(t, h, tree.Nodes[6].Hash)
}

func TestTreeWithHashes(t *testing.T) {
	hash1 := sha([]byte("test1"))
	hash2 := sha([]byte("test2"))
	hash3 := sha([]byte("test3"))
	hash4 := sha([]byte("test4"))

	tree, err := NewTreeWithHashes([]struct {
		key   string
		value []byte
	}{
		{"test1", hash1},
		{"test2", hash2},
		{"test3", hash3},
		{"test4", hash4},
	}, crypto.SHA256)

	assert.Nil(t, err)
	assert.NotNil(t, tree)
	assert.Equal(t, 7, len(tree.Nodes))
	assert.Equal(t, 3, tree.Depth)
	// Assert leaf hashes
	assertHash(t, tree.Nodes[0].Hash, fileHash("test1", hash1))
	assertHash(t, tree.Nodes[1].Hash, fileHash("test2", hash2))
	assertHash(t, tree.Nodes[2].Hash, fileHash("test3", hash3))
	assertHash(t, tree.Nodes[3].Hash, fileHash("test4", hash4))
	// Assert parent node hashes
	h := sha(tree.Nodes[0].Hash, tree.Nodes[1].Hash)
	assertHash(t, h, tree.Nodes[4].Hash)
	h = sha(tree.Nodes[2].Hash, tree.Nodes[3].Hash)
	assertHash(t, h, tree.Nodes[5].Hash)
	h = sha(tree.Nodes[4].Hash, tree.Nodes[5].Hash)
	assertHash(t, h, tree.Nodes[6].Hash)
}
