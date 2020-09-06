package serve

import (
	"fmt"
	"os"

	"github.com/heynemann/hyper-cas/storage"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

type App struct {
	Port    int
	Storage storage.Storage
}

func getStorage(storageType storage.StorageType, rootPath string) (storage.Storage, error) {
	switch storageType {
	case storage.Memory:
		return storage.NewMemStorage()
	case storage.FileSystem:
		return storage.NewFSStorage(rootPath)
	}

	return nil, fmt.Errorf("No storage could be found for storage type %v", storageType)
}

func NewApp(port int, rootPath string, storageType storage.StorageType) (*App, error) {
	storage, err := getStorage(storageType, rootPath)
	if err != nil {
		return nil, err
	}
	return &App{Port: port, Storage: storage}, nil
}

func (app *App) ListenAndServe() {
	router := routing.New()
	fileHandler := NewFileHandler(app)
	router.Put("/", fileHandler.handlePut)
	router.Get("/<hash>", fileHandler.handleGet)

	fmt.Printf("Running hyper-cas API in http://0.0.0.0:%d...\n", app.Port)
	err := fasthttp.ListenAndServe(fmt.Sprintf(":%d", app.Port), router.HandleRequest)
	if err != nil {
		fmt.Printf("Running hyper-cas API failed with %v", err)
		os.Exit(1)
	}
}
