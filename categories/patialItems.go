package categories

import (
	"time"
	"log"
)

func calcatePricePartial(c *Category)(err error){
	var chans =  []chan error{}
	for i,_ := range c.items {
		tmp := make(chan error)
		chans = append(chans, tmp)
		chans[i] = tmp
		go goCalculatePrice(c, i, tmp)
	}
	for _, ch := range chans{
		cErr, ok := <- ch
		if ok && cErr != nil{
			err = cErr
		}
	}
	c.lastUpdatePartial = time.Now()
	return err
}

func goCalculatePrice(c *Category, i int, ch chan error){
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Error al ejecutar refresco parcial de items: ", err)

		}
		close(ch)
	}()
	var err error
	var newItem *Item
	for x, item := range c.items[i] {
		if newItem, err = findItem(item.Id); err== nil{
			if newItem.Last_Update.After(item.Last_Update) {
				recalculatePrice(c, newItem, &c.items[i][x])
				c.items[i][x] = *newItem
			}
		}else{
			break
		}
	}
	ch <- err
}

func recalculatePrice(c *Category, newItem *Item, item *Item){
	if c.prices.Min>newItem.Price{
		c.prices.Min = newItem.Price
	}else if c.prices.Max < newItem.Price{
		c.prices.Max = newItem.Price
	}
	if c.cant > 0 {
		total := float64(c.cant) * c.prices.Suggested
		c.prices.Suggested = (total - item.Price + newItem.Price) / float64(c.cant)
	}
}
