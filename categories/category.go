package categories

import (
	"sync"
	"sync/atomic"
	config2 "EjercicioIntegradorGo/config"
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
		if calculatePrice(c.Id) {
			atomic.StoreUint32(&c.initialized, 1)
		}
	}

	return &c.prices
}

func calculatePrice(category string){
	config := config2.GetInstance()
	var fillPrice FillPrice = config.MethodFill
	mItem, _ := fillPrice(category)

	for _, vItem := range mItem {
		for _, item := range vItem {
			
		}
	}
}
