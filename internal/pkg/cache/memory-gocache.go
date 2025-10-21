package cache

import (
	"fiber-boilerplate/internal/defs"
	logging "fiber-boilerplate/internal/pkg/logging"
	"io"
	"time"

	goCache "github.com/patrickmn/go-cache"
)

// MemoryGoCacheBlock :
type MemoryGoCacheBlock struct {
	provider *goCache.Cache
}

func newMemoryGoCache(ttlSec int) *MemoryGoCacheBlock {
	ttl := time.Duration(ttlSec) * time.Second

	c := MemoryGoCacheBlock{
		provider: goCache.New(ttl, ttl*2),
	}

	return &c
}

// Set :
func (c *MemoryGoCacheBlock) Set(key string, value interface{}) error {
	c.provider.Set(key, value, goCache.DefaultExpiration)
	return nil
}

// Lookup :
func (c *MemoryGoCacheBlock) Lookup(key string) (found bool) {
	_, found = c.provider.Get(key)
	return
}

// Get :
func (c *MemoryGoCacheBlock) Get(key string) (value interface{}, found bool, err error) {
	value, found = c.provider.Get(key)
	err = nil
	return
}

// Del :
func (c *MemoryGoCacheBlock) Del(key string) error {
	c.provider.Delete(key)
	return nil
}

// Clear :
func (c *MemoryGoCacheBlock) Clear() error {
	c.provider.Flush()
	return nil
}

// Reader :
func (c *MemoryGoCacheBlock) Reader(key string) (io.Reader, error) {
	err := defs.ErrInvalid
	logging.Warn(err, "key: %s", key)
	return nil, err
}

// FilePath :
func (c *MemoryGoCacheBlock) FilePath(key string) string {
	err := defs.ErrInvalid
	logging.Warn(err, "key: %s", key)
	return ""
}

// List :
func (c *MemoryGoCacheBlock) List(prefix string) ([]string, error) {
	err := defs.ErrInvalid
	logging.Warn(err, "prefix: %s", prefix)
	return nil, err
}
