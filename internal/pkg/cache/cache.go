package cache

import (
	"io"

	"fiber-boilerplate/internal/defs"
)

// StoreInterface : 캐쉬 스토어 인터페이스
type StoreInterface interface {
	Set(key string, value interface{}) error
	Lookup(key string) bool
	Get(key string) (interface{}, bool, error)
	Reader(key string) (io.Reader, error)
	Del(key string) error
	Clear() error
	FilePath(key string) string
	List(prefix string) ([]string, error)
}

// StoreEnum : 캐쉬 스토어 Enum
type StoreEnum int

// StoreEnums
const (
	StoreMemoryRistretto StoreEnum = iota // ristretto
	StoreMemoryGoCache                    // go-cache
	storeCount

	StoreMemoryDefault = StoreMemoryGoCache
)

// StoreEnumNames :
var StoreEnumNames = [storeCount]string{"ristretto", "go-cache"}

func (id StoreEnum) String() string {
	return StoreEnumNames[id]
}

// New :
func New(id StoreEnum, data interface{}) StoreInterface {
	switch id {
	case StoreMemoryRistretto:
		return newMemoryRistretto(data.(int))
	case StoreMemoryGoCache:
		return newMemoryGoCache(data.(int))
	default:
		panic(defs.ErrInvalid)
	}
}
