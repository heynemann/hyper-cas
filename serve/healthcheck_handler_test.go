package serve

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vtex/hyper-cas/storage"
	"github.com/vtex/hyper-cas/utils"
)

func TestHandler(t *testing.T) {
	app, err := NewApp(200, storage.FileSystem)
	assert.Nil(t, err)
	handler := NewHealthcheckHandler(app)

	_, status, body, err := utils.DoRequest(app, handler.handleGet, "GET", "http://healthcheck/", "")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, status, 200)
	assert.Equal(t, body, "OK")
}
