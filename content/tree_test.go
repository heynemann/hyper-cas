package content

import (
	"crypto"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTreeBuild(t *testing.T) {
	var sb strings.Builder
	for i := 0; i < 2048; i++ {
		sb.WriteString(".")
	}
	content := sb.String()

	tree, err := NewTree([]byte(content), 1024, crypto.SHA256)

	assert.Nil(t, err)
	assert.NotNil(t, tree)
	assert.Equal(t, 3, len(tree.Nodes))
	assert.Equal(t, 2, tree.Depth)
	assert.NotNil(t, tree.Nodes[0])
	assert.NotNil(t, tree.Nodes[1])
	assert.NotNil(t, tree.Nodes[2])
	hash := crypto.SHA256.New()
	hash.Write([]byte(tree.Nodes[0].Hash + tree.Nodes[1].Hash))
	h := fmt.Sprintf("%x", hash.Sum(nil))
	assert.Equal(t, h, tree.Nodes[2].Hash)
}

func TestTreeBuildIrregular(t *testing.T) {
	var sb strings.Builder
	for i := 0; i < 2300; i++ {
		sb.WriteString(".")
	}
	content := sb.String()

	tree, err := NewTree([]byte(content), 1024, crypto.SHA256)

	assert.Nil(t, err)
	assert.NotNil(t, tree)
	assert.Equal(t, 7, len(tree.Nodes))
	assert.Equal(t, 3, tree.Depth)
	assert.NotNil(t, tree.Nodes[0])
	assert.NotNil(t, tree.Nodes[1])
	hash := crypto.SHA256.New()
	hash.Write([]byte(tree.Nodes[0].Hash + tree.Nodes[1].Hash))
	h := fmt.Sprintf("%x", hash.Sum(nil))
	assert.Equal(t, h, tree.Nodes[4].Hash)
	hash = crypto.SHA256.New()
	hash.Write([]byte(tree.Nodes[2].Hash + tree.Nodes[3].Hash))
	h = fmt.Sprintf("%x", hash.Sum(nil))
	assert.Equal(t, h, tree.Nodes[5].Hash)
	hash = crypto.SHA256.New()
	hash.Write([]byte(tree.Nodes[4].Hash + tree.Nodes[5].Hash))
	h = fmt.Sprintf("%x", hash.Sum(nil))
	assert.Equal(t, h, tree.Nodes[6].Hash)
}
