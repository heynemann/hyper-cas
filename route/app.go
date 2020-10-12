package route

import (
	"fmt"
	"os"

	"github.com/heynemann/hyper-cas/cache"
	"github.com/heynemann/hyper-cas/hash"
	"github.com/heynemann/hyper-cas/storage"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

type App struct {
	Port             int
	Storage          storage.Storage
	Cache            cache.Cache
	DistroExtractors []DistroExtractor
}

func getStorage(storageType storage.StorageType) (storage.Storage, error) {
	switch storageType {
	case storage.Memory:
		return storage.NewMemStorage()
	case storage.FileSystem:
		return storage.NewFSStorage(nil)
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

func getDistroExtractors() ([]DistroExtractor, error) {
	types := viper.GetStringSlice("distroExtractors")
	extractors := []DistroExtractor{}
	for _, extractorType := range types {
		switch extractorType {
		case "path":
			extractor := NewPathExtractor()
			extractors = append(extractors, extractor)
		case "subdomain":
			extractor := NewSubdomainExtractor()
			extractors = append(extractors, extractor)
		}
	}
	return extractors, nil
}

func NewApp(port int, hasherType hash.HasherType, storageType storage.StorageType, cacheType cache.CacheType) (*App, error) {
	storage, err := getStorage(storageType)
	if err != nil {
		return nil, err
	}

	cache, err := getCache(cacheType)
	if err != nil {
		return nil, err
	}

	distroExtractors, err := getDistroExtractors()
	if err != nil {
		return nil, err
	}
	return &App{Port: port, Storage: storage, Cache: cache, DistroExtractors: distroExtractors}, nil
}

func (app *App) GetDistroAndPath(host, path string, header func(string) []byte) (string, string, error) {
	for _, extractor := range app.DistroExtractors {
		distro, path, err := extractor.ExtractDistroAndPath(host, path, header)
		if err == nil {
			return distro, path, nil
		}
	}
	return "", "", fmt.Errorf("No extractors could extract a distribution.")
}

func (app *App) ListenAndServe() {
	fileHandler := NewFileHandler(app)

	fmt.Printf("Running hyper-cas router in http://0.0.0.0:%d...\n", app.Port)
	err := fasthttp.ListenAndServe(fmt.Sprintf(":%d", app.Port), fileHandler.handleGet)
	if err != nil {
		fmt.Printf("Running hyper-cas router failed with %v", err)
		os.Exit(1)
	}
}
