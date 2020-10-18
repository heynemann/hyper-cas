package serve

import (
	"fmt"
	"os"

	router "github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"github.com/vtex/hyper-cas/sitebuilder"
	"github.com/vtex/hyper-cas/storage"
)

type App struct {
	Port        int
	Storage     storage.Storage
	SiteBuilder sitebuilder.SiteBuilder
}

func getStorage(storageType storage.StorageType, siteBuilder sitebuilder.SiteBuilder) (storage.Storage, error) {
	switch storageType {
	case storage.FileSystem:
		return storage.NewFSStorage(siteBuilder)
	}

	return nil, fmt.Errorf("No storage could be found for storage type %v", storageType)
}

func getSiteBuilder() (sitebuilder.SiteBuilder, error) {
	siteBuilder, err := sitebuilder.NewNginxSiteBuilder()
	return siteBuilder, err
}

func NewApp(port int, storageType storage.StorageType) (*App, error) {
	siteBuilder, err := getSiteBuilder()
	if err != nil {
		return nil, err
	}

	storage, err := getStorage(storageType, siteBuilder)
	if err != nil {
		return nil, err
	}

	return &App{Port: port, Storage: storage, SiteBuilder: siteBuilder}, nil
}

func (app *App) HandleError(handler func(ctx *fasthttp.RequestCtx) error) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		err := handler(ctx)
		if err != nil {
			fmt.Fprintf(ctx, "Error: %v\n", err)
			ctx.SetStatusCode(500)
		}
	}
}

func (app *App) GetRouter() *router.Router {
	router := router.New()
	healthcheckHandler := NewHealthcheckHandler(app)
	fileHandler := NewFileHandler(app)
	distroHandler := NewDistroHandler(app)
	labelHandler := NewLabelHandler(app)

	router.GET("/healthcheck", app.HandleError(healthcheckHandler.handleGet))

	router.PUT("/file", app.HandleError(fileHandler.handlePut))
	router.GET("/file/{hash}", app.HandleError(fileHandler.handleGet))
	router.HEAD("/file/{hash}", app.HandleError(fileHandler.handleHead))

	router.PUT("/distro", app.HandleError(distroHandler.handlePut))
	router.GET("/distro/{distro}", app.HandleError(distroHandler.handleGet))
	router.HEAD("/distro/{distro}", app.HandleError(distroHandler.handleHead))

	router.PUT("/label", app.HandleError(labelHandler.handlePut))
	router.GET("/label/{label}", app.HandleError(labelHandler.handleGet))
	router.HEAD("/label/{label}", app.HandleError(labelHandler.handleHead))

	return router
}

func (app *App) ListenAndServe() {
	router := app.GetRouter()
	fmt.Printf("Running hyper-cas API in http://0.0.0.0:%d...\n", app.Port)
	err := fasthttp.ListenAndServe(fmt.Sprintf(":%d", app.Port), router.Handler)
	if err != nil {
		fmt.Printf("Running hyper-cas API failed with %v", err)
		os.Exit(1)
	}
}
