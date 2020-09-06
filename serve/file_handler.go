package serve

import (
	"fmt"

	routing "github.com/qiangxue/fasthttp-routing"
)

type FileHandler struct {
	App *App
}

func NewFileHandler(app *App) *FileHandler {
	return &FileHandler{App: app}
}

func (handler *FileHandler) handlePut(ctx *routing.Context) error {
	value := ctx.Request.Body()
	hash, err := handler.App.Storage.Store(value)
	if err != nil {
		return err
	}
	fmt.Fprintf(ctx, "%s", hash)
	return nil
}

func (handler *FileHandler) handleGet(ctx *routing.Context) error {
	hash := ctx.Param("hash")
	contents, err := handler.App.Storage.Get(hash)
	if err != nil {
		return err
	}
	ctx.SetBody(contents)
	return nil
}
