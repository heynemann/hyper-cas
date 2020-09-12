package serve

import (
	"fmt"
	"strings"
	"time"

	"github.com/heynemann/hyper-cas/utils"
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
	zipped, err := utils.Zip(value)
	if err != nil {
		return err
	}
	err = handler.App.Storage.Store(strHash, zipped)
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

func writeContents(ctx *routing.Context, contents []byte) error {
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

func (handler *FileHandler) handleGet(ctx *routing.Context) error {
	hash := ctx.Param("hash")
	cached, err := handler.App.Cache.Get(hash)
	if err != nil {
		return err
	}
	if cached != nil {
		return writeContents(ctx, cached)
	}
	contents, err := handler.App.Storage.Get(hash)
	if err != nil {
		return err
	}
	err = handler.App.Cache.Set(hash, contents)
	if err != nil {
		return err
	}
	return writeContents(ctx, contents)
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
