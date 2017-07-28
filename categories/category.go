package categories

import (
	"sync"
	"sync/atomic"
	"time"
	"errors"
	"EjercicioIntegradorGo/config"
	"log"
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
		go validateState(c)
		return c.prices, nil
	}
	//Tengo que probar de sincronizar esto con una gorutina
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.initialized == 0 {
		if err := calculatePrice(c); err == nil {
			atomic.StoreUint32(&c.initialized, 1)
		}else{
			return nil, err
		}
	}

	return c.prices, nil
}

func validateState(c *Category)(err error){
	conf, _ := config.GetInstance()
	if time.Since(c.lastUpdateTotal).Hours() > conf.HourUpdateTotal.Hours(){
		err = calculatePrice(c)
	}else if time.Since(c.lastUpdatePartial).Minutes() > conf.MinUpdatePartial.Minutes(){
		//calcatePricePartial(c) Pregunto item x item cual fue el que tuvo novedad y actualizo las precios solo en base a ese.
	}
	if (err != nil){
		log.Fatal("ValidateState: ", err)
	}
	return
}

func  calculatePrice( c *Category)(err error){
	conf, err := config.GetInstance()
	if err != nil {return err}
	var fillPrice FillPrice
	if conf.MethodFill == 1{
		fillPrice = FillPriceByRelevance
	}else{
		fillPrice = FillPriceTotalItems
	}
	mItem, err := fillPrice(c.Id)
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
