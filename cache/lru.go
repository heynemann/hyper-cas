package cache

import (
	"time"

	"github.com/karlseguin/ccache/v2"
	"github.com/spf13/viper"
)

type LRUCache struct {
	maxCacheSize int64
	cache        *ccache.Cache
	distroCache  *ccache.Cache
}

type CacheDistro struct {
	items map[string]string
}

func (cd *CacheDistro) Size() int64 {
	return int64(len(cd.items))
}

func newCacheDistro(value map[string]string) *CacheDistro {
	return &CacheDistro{
		items: value,
	}
}

type CacheItem struct {
	value []byte
}

func (ci *CacheItem) Size() int64 {
	return int64(len(ci.value))
}

func newCacheItem(value []byte) *CacheItem {
	return &CacheItem{
		value: value,
	}
}

func NewLRUCache() (*LRUCache, error) {
	maxCacheSize := viper.GetInt64("cache.maxSizeBytes")
	maxDistros := viper.GetInt64("cache.maxDistros")
	cache := ccache.New(ccache.Configure().MaxSize(maxCacheSize).ItemsToPrune(100))
	distroCache := ccache.New(ccache.Configure().MaxSize(maxDistros).ItemsToPrune(100))
	return &LRUCache{maxCacheSize: maxCacheSize, cache: cache, distroCache: distroCache}, nil
}

func (c *LRUCache) Get(hash string) ([]byte, error) {
	item := c.cache.Get(hash)
	if item == nil {
		return nil, nil
	}
	val := item.Value().(*CacheItem).value
	return val, nil
}

func (c *LRUCache) Set(hash string, value []byte) error {
	item := newCacheItem(value)
	c.cache.Set(hash, item, time.Hour*1000000)
	return nil
}

func (c *LRUCache) GetDistro(hash string) (map[string]string, error) {
	item := c.cache.Get(hash)
	if item == nil {
		return nil, nil
	}
	val := item.Value().(*CacheDistro).items
	return val, nil
}

func (c *LRUCache) SetDistro(hash string, hashes map[string]string) error {
	item := newCacheDistro(hashes)
	c.cache.Set(hash, item, time.Hour*1000000)
	return nil
}
