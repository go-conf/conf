package conf

import (
	"sync"
)

// The interface that Loaders must implement
type Loader interface {
	Load(c *Conf) error
}

type LoaderFactory func() Loader

var registry = make(map[string]LoaderFactory)
var registryLock sync.Mutex

// RegisterLoader allows one to register a new loader to be loaded by string.
func RegisterLoader(key string, factory LoaderFactory) {
	registryLock.Lock()
	registry[key] = factory
	registryLock.Unlock()
}
