package cache

import (
	"sync"
)

type (
	Cache interface {
		Get() (float64, error)
		Update(price float64)
	}
	cache struct {
		latestPrice float64
		lock        *sync.RWMutex
	}
)

func New(initialPrice float64) Cache {
	return &cache{
		latestPrice: initialPrice,
		lock:        &sync.RWMutex{},
	}
}
func (c cache) Get() (float64, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.latestPrice, nil
}

func (c *cache) Update(newValue float64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.latestPrice = newValue
}
