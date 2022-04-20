package leopard

import (
	"fmt"
	"github.com/volix-dev/leopard/caching"
)

type Caching struct {
	Driver caching.Driver
}

func newCaching(driver string, config any) (*Caching, error) {
	d, err := caching.New(driver, config)

	if err != nil {
		return nil, err
	}

	return &Caching{
		Driver: d,
	}, nil
}

// Caching functions

// Get retrieves data from the cache.
func (c Caching) Get(key string, target any) (bool, error) {
	get, err := c.Driver.Get("cache:"+key, target)

	if err != nil {
		return false, err
	}

	if get {
		fmt.Println(target)

	}

	return get, nil
}

// Set stores data in the cache.
func (c Caching) Set(key string, value any) error {
	return c.Driver.Set("cache:"+key, value)
}

// SetTTL stores data in the cache with a TTL.
func (c Caching) SetTTL(key string, value any, ttl int) error {
	return c.Driver.SetTTL("cache:"+key, value, ttl)
}

// Delete removes data from the cache.
func (c Caching) Delete(key string) error {
	return c.Driver.Delete("cache:" + key)
}

// Connection functions

func (c Caching) open() error {
	return c.Driver.Open()
}

func (c Caching) close() error {
	return c.Driver.Close()
}
