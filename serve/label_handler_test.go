package serve

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"path"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/vtex/hyper-cas/storage"
	"github.com/vtex/hyper-cas/utils"
)

func TestLabelHandlerPut(t *testing.T) {
	app, err := NewApp(200, storage.FileSystem)
	assert.Nil(t, err)
	label := "test"
	hash := fmt.Sprintf("%x", utils.Hash("qwe"))

	form := url.Values{}
	form.Add("label", label)
	form.Add("hash", hash)
	_, status, body, err := utils.DoRequest(app, "PUT", "/label", form.Encode())

	assert.NoError(t, err)
	assert.Equal(t, status, 200)
	assert.Empty(t, body)
	if status == 200 && err == nil {
		filePath := path.Join(viper.GetString("storage.rootPath"), "labels", label)
		assert.True(t, utils.FileExists(filePath), "Should exist: %s", filePath)
		dat, err := ioutil.ReadFile(filePath)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, hash, string(dat))
		confPath := path.Join(viper.GetString("storage.sitesPath"), fmt.Sprintf("%s.conf", label))
		assert.True(t, utils.FileExists(confPath), "Should exist: %s", confPath)
		dat, err = ioutil.ReadFile(confPath)
		if err != nil {
			t.Error(err)
		}
		assert.NotEmpty(t, string(dat))
	}
}

func TestLabelHandlerGet(t *testing.T) {
	app, err := NewApp(200, storage.FileSystem)
	assert.Nil(t, err)
	label := "test"
	hash := fmt.Sprintf("%x", utils.Hash("qwe"))
	form := url.Values{}
	form.Add("label", label)
	form.Add("hash", hash)
	_, status, body, err := utils.DoRequest(app, "PUT", "/label", form.Encode())
	assert.NoError(t, err)
	assert.Equal(t, 200, status)

	_, status, body, err = utils.DoRequest(app, "GET", fmt.Sprintf("/label/%s", label), "")

	assert.NoError(t, err)
	assert.Equal(t, status, 200)
	assert.Equal(t, hash, body)
}

func TestLabelHandlerGetNotFound(t *testing.T) {
	app, err := NewApp(200, storage.FileSystem)
	assert.Nil(t, err)

	_, status, body, err := utils.DoRequest(app, "GET", "/label/invalidhash", "")

	assert.NoError(t, err)
	assert.Equal(t, 404, status)
	assert.Equal(t, "", body)
}

func TestLabelHandlerHead(t *testing.T) {
	app, err := NewApp(200, storage.FileSystem)
	assert.Nil(t, err)
	label := "test"
	hash := fmt.Sprintf("%x", utils.Hash("qwe"))
	form := url.Values{}
	form.Add("label", label)
	form.Add("hash", hash)
	_, status, body, err := utils.DoRequest(app, "PUT", "/label", form.Encode())
	assert.NoError(t, err)
	assert.Equal(t, 200, status)

	_, status, body, err = utils.DoRequest(app, "HEAD", fmt.Sprintf("/label/%s", label), "")

	assert.Equal(t, 200, status)
	assert.Equal(t, "", body)
}

func TestLabelHandlerHeadNotFound(t *testing.T) {
	app, err := NewApp(200, storage.FileSystem)
	assert.Nil(t, err)

	_, status, body, err := utils.DoRequest(app, "HEAD", "/label/invalidhash", "")

	assert.NoError(t, err)
	assert.Equal(t, 404, status)
	assert.Equal(t, "", body)
}
