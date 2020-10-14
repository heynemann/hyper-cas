package serve

import (
	"fmt"
	"os"

	routing "github.com/qiangxue/fasthttp-routing"
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

func (app *App) ListenAndServe() {
	router := routing.New()
	fileHandler := NewFileHandler(app)
	distroHandler := NewDistroHandler(app)
	labelHandler := NewLabelHandler(app)
	router.Put("/file", fileHandler.handlePut)
	router.Get("/file/<hash>", fileHandler.handleGet)
	router.Head("/file/<hash>", fileHandler.handleHead)
	router.Put("/distro", distroHandler.handlePut)
	router.Put("/label", labelHandler.handlePut)
	router.Get("/label/<label>", labelHandler.handleGet)
	router.Head("/label/<label>", labelHandler.handleHead)

	fmt.Printf("Running hyper-cas API in http://0.0.0.0:%d...\n", app.Port)
	err := fasthttp.ListenAndServe(fmt.Sprintf(":%d", app.Port), router.HandleRequest)
	if err != nil {
		fmt.Printf("Running hyper-cas API failed with %v", err)
		os.Exit(1)
	}
}
