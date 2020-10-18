package serve

import (
	"fmt"
	"strings"
	"time"

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
	ctx.SetBody([]byte(strHash))
	return nil
}

func writeContents(ctx *fasthttp.RequestCtx, contents []byte) error {
	gzipEnabled := strings.Contains(string(ctx.Request.Header.Peek("Accept-Encoding")), "gzip")
	ctx.Response.Header.Add("content-type", "text/html; charset=utf-8")
	ctx.Response.Header.Add("date", time.Now().Format("RFC1123"))
	ctx.Response.Header.Set("server", "hyper-cas")
	if gzipEnabled {
		ctx.Response.Header.Add("content-encoding", "gzip")
		ctx.SetBody(contents)
	} else {
		res, err := utils.Unzip(contents)
		if err != nil {
			return err
		}
		ctx.SetBody(res)
	}

	return nil
}

func (handler *FileHandler) handleGet(ctx *fasthttp.RequestCtx) error {
	hash := ctx.UserValue("hash").(string)
	contents, err := handler.App.Storage.Get(hash)
	if err != nil {
		return err
	}
	return writeContents(ctx, contents)
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
