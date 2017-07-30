package categories

import (
	"time"
	"errors"
	"EjercicioIntegradorGo/config"
)

func  fillAllPrice( c *Category)(err error) {
	conf, err := config.GetInstance()
	if err != nil {
		return err
	}
	var fillPrice FillPrice
	if conf.MethodFill == 1 {
		fillPrice = FillPriceByRelevance
	} else {
		fillPrice = FillPriceTotalItems
	}
	c.items, err = fillPrice(c.Id)
	if err != nil {
		return
	}
	c.lastUpdatePartial = time.Now()
	c.lastUpdateTotal = time.Now()
	return calculateTotalPrice(c)
}

func calculateTotalPrice (c *Category) error{
	c.prices = &Prices{}
	var priceT float64
	for _, vItem := range c.items {
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

	if c.cant > 0 {
		c.prices.Suggested = priceT / float64(c.cant)
	}else{
		return errors.New("La categoria no tiene articulos a la venta")
	}
	return nil
}