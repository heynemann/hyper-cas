package serve

import (
	"fmt"
	"os"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

type ServeHandler struct {
	App *App
}

func (handler *ServeHandler) handleGet(ctx *routing.Context) error {
	fmt.Fprintf(ctx, "OK\n")
	return nil
}

type App struct {
	Port int
}

func (app *App) ListenAndServe() {
	handler := ServeHandler{
		App: app,
	}

	router := routing.New()
	router.Get("/", handler.handleGet)

	fmt.Printf("Running hyper-cas API in http://0.0.0.0:%d...\n", app.Port)
	err := fasthttp.ListenAndServe(fmt.Sprintf(":%d", app.Port), router.HandleRequest)
	if err != nil {
		fmt.Printf("Running hyper-cas API failed with %v", err)
		os.Exit(1)
	}
}
