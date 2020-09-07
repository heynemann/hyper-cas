package distro

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsingEmptyDistro(t *testing.T) {
	data := `labelName`

	distro, err := ParseDistro(data)

	assert.Nil(t, err)
	assert.NotNil(t, distro)
	assert.Equal(t, distro.Label, "labelName")
}

func TestParsingDistroWithPaths(t *testing.T) {
	data := `labelName
some/path/to/file.txt@72de8eb2853d5a5c89f256c96966f101987b2596
other/path/to/file.txt@72de8eb2853d5a5c89f256c96966f101987b2596`

	distro, err := ParseDistro(data)

	assert.Nil(t, err)
	assert.NotNil(t, distro)
	assert.Equal(t, "labelName", distro.Label)
	assert.Equal(t, 2, len(distro.PathToHash))
	assert.Equal(t, "72de8eb2853d5a5c89f256c96966f101987b2596", distro.PathToHash["some/path/to/file.txt"])
	assert.Equal(t, "72de8eb2853d5a5c89f256c96966f101987b2596", distro.PathToHash["other/path/to/file.txt"])
}
