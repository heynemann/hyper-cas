package serve

import (
	"fmt"
	"os"

	router "github.com/fasthttp/router"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"github.com/vtex/hyper-cas/sitebuilder"
	"github.com/vtex/hyper-cas/storage"
	"github.com/vtex/hyper-cas/utils"
	"go.uber.org/zap"
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
	viper.SetDefault("serve.maxRequestBodySize", 4*1024*1024*1024)
	viper.SetDefault("serve.TCPKeepaliveEnabled", true)

	siteBuilder, err := getSiteBuilder()
	if err != nil {
		utils.LogError("Could not create site builder.", zap.Error(err))
		return nil, err
	}

	storage, err := getStorage(storageType, siteBuilder)
	if err != nil {
		utils.LogError("Could not create storage.", zap.Error(err))
		return nil, err
	}

	return &App{Port: port, Storage: storage, SiteBuilder: siteBuilder}, nil
}

func (app *App) HandleError(handler func(ctx *fasthttp.RequestCtx) error) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		err := handler(ctx)
		if err != nil {
			ctx.SetBodyString(fmt.Sprintf("Error: %v\n", err))
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
	logger := utils.LoggerWith(
		zap.String("ip", "0.0.0.0"),
		zap.Int("port", app.Port),
	)
	logger.Info("hyper-cas API running successfully.")
	s := &fasthttp.Server{
		Handler: router.Handler,
		Name:    "hyper-cas",

		MaxRequestBodySize: viper.GetInt("serve.maxRequestBodySize"),
		DisableKeepalive:   false,
		TCPKeepalive:       viper.GetBool("serve.TCPKeepaliveEnabled"),
	}
	err := s.ListenAndServe(
		fmt.Sprintf(":%d", app.Port),
	)
	if err != nil {
		logger.Error("Running hyper-cas API failed.", zap.Error(err))
		os.Exit(1)
	}
}
