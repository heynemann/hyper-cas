package serve

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

type ServeHandler struct {
	App *App
}

func (handler *ServeHandler) handleRequest(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/":
	// case "/foo":
	// fooHandler(ctx)
	// case "/bar":
	// barHandler(ctx)
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}

func (handler *ServeHandler) handleGet(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("OK")
}

type App struct {
	Port int
}

func (app *App) Run() {
	handler := ServeHandler{
		App: app,
	}
	fasthttp.ListenAndServe(fmt.Sprintf("%d", app.Port), handler.handleRequest)
}
