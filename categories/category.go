package categories

import (
	"sync"
	"sync/atomic"
	"time"
	"EjercicioIntegradorGo/config"
	"log"
)

type Prices struct {
	Max float64 `json:"max"`
	Suggested float64 `json:"suggested"`
	Min float64 `json:"min"`
}

type Category struct {
	Id string `json:"id"`
	items [][]Item
	cant int
	prices *Prices
	lastUpdatePartial time.Time
	lastUpdateTotal time.Time
	LastEntry time.Time
	initialized uint32
	mu sync.Mutex
}

func (c *Category) GetPrices()(*Prices, error){

	if atomic.LoadUint32(&c.initialized) == 1 {
		return c.prices, nil
	}
	//Tengo que probar de sincronizar esto con una gorutina en vez de Mutex
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.initialized == 0 {
		if err := fillAllPrice(c); err == nil {
			atomic.StoreUint32(&c.initialized, 1)
		}else{
			return nil, err
		}
	}

	return c.prices, nil
}

func (c *Category) validateState() (err error){
	conf, _ := config.GetInstance()

	if c.lastUpdateTotal.Add(time.Hour *conf.HourUpdateTotal).Before(time.Now()){
		err = fillAllPrice(c)
	}else if c.lastUpdatePartial.Add(time.Minute *conf.MinUpdatePartial).Before(time.Now()){
		err = calcatePricePartial(c) //Pregunto item x item cual fue el que tuvo novedad y actualizo las precios solo en base a ese.
	}
	if err != nil{
		log.Fatal("ValidateState: ", err)
	}
	return
}



