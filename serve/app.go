package serve

import (
	"fmt"
	"os"

	"github.com/heynemann/hyper-cas/cache"
	"github.com/heynemann/hyper-cas/hash"
	"github.com/heynemann/hyper-cas/storage"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

type App struct {
	Port    int
	Hasher  hash.Hasher
	Storage storage.Storage
	Cache   cache.Cache
}

func getStorage(storageType storage.StorageType) (storage.Storage, error) {
	switch storageType {
	case storage.Memory:
		return storage.NewMemStorage()
	case storage.FileSystem:
		return storage.NewFSStorage()
	}

	return nil, fmt.Errorf("No storage could be found for storage type %v", storageType)
}

func getCache(cacheType cache.CacheType) (cache.Cache, error) {
	switch cacheType {
	case cache.LRU:
		return cache.NewLRUCache()
	}

	return nil, fmt.Errorf("No cache could be found for cache type %v", cacheType)
}

func getHasher(hasherType hash.HasherType) (hash.Hasher, error) {
	switch hasherType {
	case hash.SHA1:
		return hash.NewSHA1Hasher()
	}

	return nil, fmt.Errorf("No cache could be found for cache type %v", hasherType)
}

func NewApp(port int, hasherType hash.HasherType, storageType storage.StorageType, cacheType cache.CacheType) (*App, error) {
	hasher, err := getHasher(hasherType)
	if err != nil {
		return nil, err
	}

	storage, err := getStorage(storageType)
	if err != nil {
		return nil, err
	}

	cache, err := getCache(cacheType)
	if err != nil {
		return nil, err
	}
	return &App{Port: port, Hasher: hasher, Storage: storage, Cache: cache}, nil
}

func (app *App) ListenAndServe() {
	router := routing.New()
	fileHandler := NewFileHandler(app)
	distroHandler := NewDistroHandler(app)
	router.Put("/file", fileHandler.handlePut)
	router.Get("/file/<hash>", fileHandler.handleGet)
	router.Put("/distro", distroHandler.handlePut)
	router.Get("/distro/<hash>", distroHandler.handleGet)

	fmt.Printf("Running hyper-cas API in http://0.0.0.0:%d...\n", app.Port)
	err := fasthttp.ListenAndServe(fmt.Sprintf(":%d", app.Port), router.HandleRequest)
	if err != nil {
		fmt.Printf("Running hyper-cas API failed with %v", err)
		os.Exit(1)
	}
}
