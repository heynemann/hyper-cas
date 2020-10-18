package serve

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"path"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/vtex/hyper-cas/storage"
	"github.com/vtex/hyper-cas/utils"
)

func TestFileHandlerPut(t *testing.T) {
	app, err := NewApp(200, storage.FileSystem)
	assert.Nil(t, err)
	text := fmt.Sprintf("some random text: %d", rand.Intn(100))

	_, status, body, err := utils.DoRequest(app, "PUT", "/file", text)

	assert.NoError(t, err)
	assert.Equal(t, status, 200)
	assert.NotEmpty(t, body)
	if status == 200 && err == nil {
		filePath := path.Join(viper.GetString("storage.rootPath"), "files", body[0:2], body[2:4], body)
		assert.True(t, utils.FileExists(filePath), "Should exist: %s", filePath)
		dat, err := ioutil.ReadFile(filePath)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, text, string(dat))
	}
}

func TestFileHandlerGet(t *testing.T) {
	app, err := NewApp(200, storage.FileSystem)
	assert.NoError(t, err)
	text := fmt.Sprintf("some random text: %d", rand.Intn(100))
	_, status, body, err := utils.DoRequest(app, "PUT", "/file", text)
	assert.NoError(t, err)
	hash := body

	_, status, body, err = utils.DoRequest(app, "GET", fmt.Sprintf("/file/%s", hash), "")

	assert.NoError(t, err)
	assert.Equal(t, 200, status)
	assert.Equal(t, text, body)
}

func TestFileHandlerGetNotFound(t *testing.T) {
	app, err := NewApp(200, storage.FileSystem)
	assert.Nil(t, err)

	_, status, body, err := utils.DoRequest(app, "GET", "/file/invalidhash", "")

	assert.NoError(t, err)
	assert.Equal(t, 404, status)
	assert.Equal(t, "", body)
}

func TestFileHandlerHead(t *testing.T) {
	app, err := NewApp(200, storage.FileSystem)
	assert.Nil(t, err)
	text := fmt.Sprintf("some random text: %d", rand.Intn(100))
	_, status, body, err := utils.DoRequest(app, "PUT", "/file", text)
	if err != nil {
		t.Error(err)
	}
	hash := body

	_, status, body, err = utils.DoRequest(app, "HEAD", fmt.Sprintf("/file/%s", hash), "")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 200, status)
	assert.Equal(t, "", body)
}

func TestFileHandlerHeadNotFound(t *testing.T) {
	app, err := NewApp(200, storage.FileSystem)
	assert.Nil(t, err)

	_, status, body, err := utils.DoRequest(app, "HEAD", "/file/invalidhash", "")

	assert.NoError(t, err)
	assert.Equal(t, 404, status)
	assert.Equal(t, "", body)
}
