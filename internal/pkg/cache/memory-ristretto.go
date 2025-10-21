package cache

import (
	"io"
	"time"

	"github.com/dgraph-io/ristretto"

	"fiber-boilerplate/internal/defs"
	logging "fiber-boilerplate/internal/pkg/logging"
)

// MemoryRistrettoBlock :
type MemoryRistrettoBlock struct {
	provider *ristretto.Cache
	ttl      time.Duration
}

func newMemoryRistretto(ttlSec int) *MemoryRistrettoBlock {
	c := MemoryRistrettoBlock{
		ttl: time.Duration(ttlSec) * time.Second,
	}

	var maxCost int64 = 100 * 1000 * 1000 // 100MB
	provider, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: maxCost * 10, // recommended
		MaxCost:     maxCost,
		BufferItems: 64, // recommended
	})
	if err != nil {
		panic(err)
	}

	c.provider = provider
	return &c
}

// Set :
func (c *MemoryRistrettoBlock) Set(key string, value interface{}) (err error) {
	if c.provider.SetWithTTL(key, value, 1, c.ttl) {
		err = nil
	} else {
		err = defs.ErrFault
		logging.Warn(err, "key:%s", key)
	}

	return
}

// Lookup :
func (c *MemoryRistrettoBlock) Lookup(key string) (found bool) {
	_, found = c.provider.Get(key)
	return
}

// Get :
func (c *MemoryRistrettoBlock) Get(key string) (value interface{}, found bool, err error) {
	value, found = c.provider.Get(key)
	err = nil
	return
}

// Del :
func (c *MemoryRistrettoBlock) Del(key string) error {
	c.provider.Del(key)
	return nil
}

// Clear :
func (c *MemoryRistrettoBlock) Clear() error {
	c.provider.Clear()
	// ristretto clear cleanup the cache store, policy, and metrics. so, re-create provider
	c.provider = newMemoryRistretto(int(c.ttl.Seconds())).provider
	return nil
}

// Reader :
func (c *MemoryRistrettoBlock) Reader(key string) (io.Reader, error) {
	err := defs.ErrInvalid
	logging.Warn(err, "key: %s", key)
	return nil, err
}

// FilePath :
func (c *MemoryRistrettoBlock) FilePath(key string) string {
	err := defs.ErrInvalid
	logging.Warn(err, "key: %s", key)
	return ""
}

// List :
func (c *MemoryRistrettoBlock) List(prefix string) ([]string, error) {
	err := defs.ErrInvalid
	logging.Warn(err, "prefix: %s", prefix)
	return nil, err
}
