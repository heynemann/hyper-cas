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

func NewApp(port int) (*App, error) {
	storage, err := storage.NewMemStorage()
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
