package conf

import (
	"sync"
)

type Loader interface {
	Load(c *Conf) error
}

type LoaderFactory func(c *Conf) Loader

var registry = make(map[string]LoaderFactory)
var registryLock sync.Mutex

func RegisterLoader(key string, factory LoaderFactory) {
	registryLock.Lock()
	registry[key] = factory
	registryLock.Unlock()
}
