package drivers

import (
	"fmt"
	"github.com/volix-dev/leopard/caching"
	"sync"
	"time"
)

type MemoryDriver struct {
	cache   map[string]any
	lock    sync.RWMutex
	expires map[string]time.Time
}

func init() {
	caching.Register("memory", func(config any) (caching.Driver, error) {
		timer := time.NewTimer(time.Minute)

		driver := &MemoryDriver{
			cache:   make(map[string]any),
			lock:    sync.RWMutex{},
			expires: make(map[string]time.Time),
		}

		go func() {
			for {
				<-timer.C

				for key, expireTime := range driver.expires {
					if expireTime.Before(time.Now()) {
						driver.lock.Lock()
						delete(driver.cache, key)
						delete(driver.expires, key)
						driver.lock.Unlock()
					}
				}
			}
		}()

		return driver, nil
	})
}

func (m *MemoryDriver) Get(key string, target any) (bool, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if value, ok := m.cache[key]; ok {
		switch t := target.(type) {
		case *string:
			*t = value.(string)
		}
		fmt.Println(target)
		return true, nil
	}

	return false, nil
}

func (m *MemoryDriver) Set(key string, value any) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.cache[key] = value
	return nil
}

func (m *MemoryDriver) SetTTL(key string, value any, ttl int) error {
	err := m.Set(key, value)
	m.expires[key] = time.Now().Add(time.Duration(ttl) * time.Second)

	return err
}

func (m *MemoryDriver) Delete(key string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.cache, key)
	return nil
}

func (m *MemoryDriver) Close() error {
	return nil
}

func (m *MemoryDriver) Open() error {
	fmt.Println(`Started memory driver
This is not intended to be used in production.
You can change it by setting CACHE_DRIVER to "redis" in the .env file or in the environment variables.`)

	return nil
}
