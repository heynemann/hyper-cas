package serve

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"github.com/vtex/hyper-cas/utils"
)

type FileHandler struct {
	App *App
}

func NewFileHandler(app *App) *FileHandler {
	viper.SetDefault("gzipSourceFiles", false)
	return &FileHandler{App: app}
}

func (handler *FileHandler) handlePut(ctx *fasthttp.RequestCtx) error {
	value := ctx.Request.Body()
	hash := utils.HashBytes(value)
	strHash := fmt.Sprintf("%x", hash)
	if viper.GetBool("gzipSourceFiles") {
		var err error
		value, err = utils.Zip(value)
		if err != nil {
			return err
		}
	}
	err := handler.App.Storage.Store(strHash, value)
	if err != nil {
		return err
	}
	ctx.SetBodyString(strHash)
	return nil
}

func (handler *FileHandler) handleGet(ctx *fasthttp.RequestCtx) error {
	hash := ctx.UserValue("hash").(string)
	contents, err := handler.App.Storage.Get(hash)
	if contents == nil {
		ctx.SetStatusCode(404)
		return nil
	}
	if err != nil {
		return err
	}
	ctx.SetBody(contents)
	return nil
}

func (handler *FileHandler) handleHead(ctx *fasthttp.RequestCtx) error {
	hash := ctx.UserValue("hash").(string)
	if has := handler.App.Storage.Has(hash); has {
		ctx.SetStatusCode(200)
	} else {
		ctx.SetStatusCode(404)
	}
	return nil
}
