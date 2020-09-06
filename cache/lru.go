package cache

import (
	"time"

	"github.com/karlseguin/ccache/v2"
	"github.com/spf13/viper"
)

type LRUCache struct {
	maxCacheSize int64
	cache        *ccache.Cache
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
	cache := ccache.New(ccache.Configure().MaxSize(maxCacheSize).ItemsToPrune(100))
	return &LRUCache{maxCacheSize: maxCacheSize, cache: cache}, nil
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
