package serve

import "github.com/valyala/fasthttp"

type HealthcheckHandler struct {
	App *App
}

func NewHealthcheckHandler(app *App) *HealthcheckHandler {
	return &HealthcheckHandler{App: app}
}

func (handler *HealthcheckHandler) handleGet(ctx *fasthttp.RequestCtx) error {
	ctx.SetBodyString("OK")
	return nil
}
