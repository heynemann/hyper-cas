package distro

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestParsingEmptyDistro(t *testing.T) {
// data := ``

// distro, err := ParseDistro([]byte(data))

// assert.Nil(t, err)
// assert.NotNil(t, distro)
// }

func TestParsingDistroWithPaths(t *testing.T) {
	d := NewDistro()
	d.AppendPath("some/path/to/file.txt", "72de8eb2853d5a5c89f256c96966f101987b2596")
	d.AppendPath("other/path/to/file.txt", "72de8eb2853d5a5c89f256c96966f101987b2596")

	data, err := json.Marshal(d)
	fmt.Println(string(data))

	distro, err := ParseDistro(data)

	assert.Nil(t, err)
	assert.NotNil(t, distro)
	assert.Equal(t, 2, len(distro.Paths))
	assert.Equal(t, "72de8eb2853d5a5c89f256c96966f101987b2596", distro.Hashes[0])
	assert.Equal(t, "72de8eb2853d5a5c89f256c96966f101987b2596", distro.Hashes[1])
}
