package serve

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

type LabelHandler struct {
	App *App
}

func NewLabelHandler(app *App) *LabelHandler {
	return &LabelHandler{App: app}
}

func (handler *LabelHandler) handlePut(ctx *fasthttp.RequestCtx) error {
	label := string(ctx.PostArgs().Peek("label"))
	hash := string(ctx.PostArgs().Peek("hash"))
	if label == "" || hash == "" {
		return fmt.Errorf("Both label and hash must be set (label: '%s', hash: '%s')", label, hash)
	}
	err := handler.App.Storage.StoreLabel(label, hash)
	if err != nil {
		return err
	}
	return nil
}

func (handler *LabelHandler) handleGet(ctx *fasthttp.RequestCtx) error {
	label := ctx.UserValue("label").(string)
	if !handler.App.Storage.HasLabel(label) {
		ctx.SetStatusCode(404)
		return nil
	}
	contents, err := handler.App.Storage.GetLabel(label)
	if err != nil {
		return err
	}
	ctx.SetBodyString(contents)
	return nil
}

func (handler *LabelHandler) handleHead(ctx *fasthttp.RequestCtx) error {
	label := ctx.UserValue("label").(string)
	if has := handler.App.Storage.HasLabel(label); has {
		ctx.SetStatusCode(200)
	} else {
		ctx.SetStatusCode(404)
	}
	return nil
}
