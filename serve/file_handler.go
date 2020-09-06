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
	value := string(ctx.Request.Body())
	hash, err := handler.App.Storage.Store(string(value))
	if err != nil {
		return err
	}
	fmt.Fprintf(ctx, hash)
	return nil
}

func (handler *FileHandler) handleGet(ctx *routing.Context) error {
	hash := ctx.Param("hash")
	contents, err := handler.App.Storage.Get(hash)
	if err != nil {
		return err
	}
	fmt.Fprintf(ctx, contents)
	return nil
}
