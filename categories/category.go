package categories

import (
	"sync"
	"sync/atomic"
	"time"
	"errors"
	"EjercicioIntegradorGo/config"
)

type Prices struct {
	Max float64 `json:"max"`
	Suggested float64 `json:"suggested"`
	Min float64 `json:"min"`
}

type Category struct {
	Id string
	items [][]Item
	cant int
	prices *Prices
	lastUpdatePartial time.Time
	lastUpdateTotal time.Time
	initialized uint32
	mu sync.Mutex
}

func (c *Category) GetPrices()(*Prices, error){
	if atomic.LoadUint32(&c.initialized) == 1 {
		return c.prices, nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.initialized == 0 {
		if err := calculatePrice(c.Id, c); err == nil {
			atomic.StoreUint32(&c.initialized, 1)
		}else{
			return nil, err
		}
	}

	return c.prices, nil
}

func  calculatePrice(category string, c *Category)(err error){
	conf, err := config.GetInstance()
	if err != nil {return err}
	var fillPrice FillPrice
	if conf.MethodFill == 1{
		fillPrice = FillPriceByRelevance
	}else{
		fillPrice = FillPriceTotalItems
	}
	mItem, err := fillPrice(category)
	if err !=nil{
		return
	}

	c.lastUpdatePartial = time.Now()
	c.lastUpdateTotal = time.Now()
	c.prices = &Prices{}
	var priceT float64
	for _, vItem := range mItem {
		c.cant += len(vItem)
		for _, item := range vItem {
			if item.Price < c.prices.Min{
				c.prices.Min = item.Price
			}else if item.Price > c.prices.Max{
				c.prices.Max = item.Price
			}
			priceT += item.Price
			item.Last_Update = c.lastUpdateTotal
		}
	}
	c.items = mItem
	if c.cant > 0 {
		c.prices.Suggested = priceT / float64(c.cant)
	}else{
		return errors.New("La categoria no existe")
	}
	return
}
