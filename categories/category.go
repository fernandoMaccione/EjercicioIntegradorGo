package categories

import (
	"sync"
	"sync/atomic"
)

type Prices struct {
	Max string `json:"max"`
	Suggested string `json:"suggested"`
	Min string `json:"min"`
}

type Category struct {
	Id string
	items []Item
	prices Prices
	initialized uint32
	mu sync.Mutex
}

func (c *Category) getPrices()*Prices{
	if atomic.LoadUint32(&c.initialized) == 1 {
		return &c.prices
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.initialized == 0 {
		if calculatePrice() {
			atomic.StoreUint32(&c.initialized, 1)
		}
	}

	return &c.prices
}

func 
