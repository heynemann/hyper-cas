package serve

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"path"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/vtex/hyper-cas/storage"
	"github.com/vtex/hyper-cas/utils"
)

func newHash() string {
	test := fmt.Sprintf("test%d", rand.Intn(10000))
	return fmt.Sprintf(
		"%s:%s",
		test,
		fmt.Sprintf("%x", utils.Hash(test)),
	)
}

func TestDistroHandlerPut(t *testing.T) {
	hash1 := newHash()
	hash2 := newHash()
	hash3 := newHash()
	hashes := []string{
		hash1, hash2, hash3,
	}
	app, err := NewApp(200, storage.FileSystem)
	assert.Nil(t, err)

	_, status, body, err := utils.DoRequest(app, "PUT", "/distro", strings.Join(hashes, "\n"))

	assert.NoError(t, err)
	assert.Equal(t, status, 200)
	assert.NotEmpty(t, body)
	if status == 200 && err == nil {
		filePath := path.Join(viper.GetString("storage.rootPath"), "distros", body)
		assert.True(t, utils.FileExists(filePath), "Should exist: %s", filePath)
		dat, err := ioutil.ReadFile(filePath)
		if err != nil {
			t.Error(err)
		}
		assert.NotEmpty(t, string(dat))
	}
}

func TestDistroHandlerPutWithWrongBody(t *testing.T) {
	app, err := NewApp(200, storage.FileSystem)
	assert.Nil(t, err)

	_, status, body, err := utils.DoRequest(app, "PUT", "/distro", "qwe")

	assert.NoError(t, err)
	assert.Equal(t, 500, status)
	assert.Equal(t, "Error: The body should be composed of lines with {filepath}:{content hash} only.\n", body)
}

func TestDistroHandlerGet(t *testing.T) {
	hash1 := newHash()
	hash2 := newHash()
	hash3 := newHash()
	hashes := []string{
		hash1, hash2, hash3,
	}
	app, err := NewApp(200, storage.FileSystem)
	assert.Nil(t, err)
	_, status, body, err := utils.DoRequest(app, "PUT", "/distro", strings.Join(hashes, "\n"))
	assert.NoError(t, err)
	assert.Equal(t, status, 200)

	_, status, body, err = utils.DoRequest(app, "GET", fmt.Sprintf("/distro/%s", body), "")

	assert.NoError(t, err)
	assert.Equal(t, 200, status)
	assert.NotEmpty(t, body)
	var content []string
	err = json.Unmarshal([]byte(body), &content)
	assert.NoError(t, err)
	assert.Len(t, content, 3)
	assert.Equal(t, content[0], hash1)
	assert.Equal(t, content[1], hash2)
	assert.Equal(t, content[2], hash3)
}

func TestDistroHandlerGetNotFound(t *testing.T) {
	app, err := NewApp(200, storage.FileSystem)
	assert.Nil(t, err)

	_, status, body, err := utils.DoRequest(app, "GET", "/distro/invalidhash", "")

	assert.NoError(t, err)
	assert.Equal(t, 404, status)
	assert.Equal(t, "", body)
}

func TestDistroHandlerHead(t *testing.T) {
	hash1 := newHash()
	hash2 := newHash()
	hash3 := newHash()
	hashes := []string{
		hash1, hash2, hash3,
	}
	app, err := NewApp(200, storage.FileSystem)
	assert.Nil(t, err)
	_, status, body, err := utils.DoRequest(app, "PUT", "/distro", strings.Join(hashes, "\n"))
	assert.NoError(t, err)
	assert.Equal(t, status, 200)

	_, status, body, err = utils.DoRequest(app, "HEAD", fmt.Sprintf("/distro/%s", body), "")

	assert.Equal(t, 200, status)
	assert.Equal(t, "", body)
}

func TestDistroHandlerHeadNotFound(t *testing.T) {
	app, err := NewApp(200, storage.FileSystem)
	assert.Nil(t, err)

	_, status, body, err := utils.DoRequest(app, "HEAD", "/distro/invalidhash", "")

	assert.NoError(t, err)
	assert.Equal(t, 404, status)
	assert.Equal(t, "", body)
}
