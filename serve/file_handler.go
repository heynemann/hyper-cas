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
	hash, err := handler.App.Hasher.Calc(value)
	if err != nil {
		return err
	}
	strHash := fmt.Sprintf("%x", hash)
	err = handler.App.Storage.Store(strHash, value)
	if err != nil {
		return err
	}
	err = handler.App.Cache.Set(strHash, value)
	if err != nil {
		return err
	}
	ctx.SetBody([]byte(strHash))
	return nil
}

func (handler *FileHandler) handleGet(ctx *routing.Context) error {
	hash := ctx.Param("hash")
	cached, err := handler.App.Cache.Get(hash)
	if err != nil {
		return err
	}
	if cached != nil {
		ctx.SetBody(cached)
		return nil
	}
	contents, err := handler.App.Storage.Get(hash)
	if err != nil {
		return err
	}
	err = handler.App.Cache.Set(hash, contents)
	if err != nil {
		return err
	}
	ctx.SetBody(contents)
	return nil
}

func (handler *FileHandler) handleHead(ctx *routing.Context) error {
	hash := ctx.Param("hash")

	if has := handler.App.Storage.Has(hash); has {
		ctx.SetStatusCode(200)
	} else {
		ctx.SetStatusCode(404)
	}
	return nil
}
