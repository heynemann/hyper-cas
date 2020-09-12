package serve

import (
	"fmt"

	routing "github.com/qiangxue/fasthttp-routing"
)

type LabelHandler struct {
	App *App
}

func NewLabelHandler(app *App) *LabelHandler {
	return &LabelHandler{App: app}
}

func (handler *LabelHandler) handlePut(ctx *routing.Context) error {
	label := string(ctx.FormValue("label"))
	hash := string(ctx.FormValue("hash"))
	if label == "" || hash == "" {
		return fmt.Errorf("Both label and hash must be set (label: %s, hash: %s)", label, hash)
	}
	err := handler.App.Storage.StoreLabel(label, hash)
	if err != nil {
		return err
	}
	return nil
}

func (handler *LabelHandler) handleGet(ctx *routing.Context) error {
	label := ctx.Param("label")
	cached, err := handler.App.Cache.Get(label)
	if err != nil {
		return err
	}
	if cached != nil {
		ctx.SetBody(cached)
		return nil
	}
	contents, err := handler.App.Storage.GetLabel(label)
	if err != nil {
		return err
	}
	err = handler.App.Cache.Set(label, []byte(contents))
	if err != nil {
		return err
	}
	ctx.SetBodyString(contents)
	return nil
}

func (handler *LabelHandler) handleHead(ctx *routing.Context) error {
	label := ctx.Param("label")
	if has := handler.App.Storage.HasLabel(label); has {
		ctx.SetStatusCode(200)
	} else {
		ctx.SetStatusCode(404)
	}
	return nil
}
