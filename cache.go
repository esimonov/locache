package locache

import (
	"sync"
	"time"
)

var singleton = newCache()

// LoadLocation returns the Location with the given name.
// It is calling time.LoadLocation under the hood,
// but caches the locations that were retrieved previously in order to speed up subsequent lookups.
func LoadLocation(name string) (*time.Location, error) {
	return singleton.LoadLocation(name)
}

type cache struct {
	*sync.RWMutex
	locs map[string]*time.Location
}

func newCache() *cache {
	return &cache{
		RWMutex: new(sync.RWMutex),
		locs:    make(map[string]*time.Location),
	}
}

func (c *cache) LoadLocation(name string) (*time.Location, error) {
	c.RLock()
	loc, ok := c.locs[name]
	c.RUnlock()

	if ok {
		return loc, nil
	}

	c.Lock()
	defer c.Unlock()

	loc, ok = c.locs[name]
	if ok {
		return loc, nil
	}

	loc, err := time.LoadLocation(name)
	if err != nil {
		return nil, err
	}

	c.locs[name] = loc

	return loc, nil
}
