package categories

import "time"

func calcatePricePartial(c *Category)  error{

	for i,_ := range c.items {
		for x, item := range c.items[i] {
			if newItem, err := findItem(item.Id); err== nil{
				if newItem.Last_Update.After(item.Last_Update) {
					recalculatePrice(c, newItem, &c.items[i][x])
					c.items[i][x] = *newItem
				}
			}else{
				return err
			}
		}
	}
	c.lastUpdatePartial = time.Now()
	return nil
}

func recalculatePrice(c *Category, newItem *Item, item *Item){
	if c.prices.Min>newItem.Price{
		c.prices.Min = newItem.Price
	}else if c.prices.Max < item.Price{
		c.prices.Max = newItem.Price
	}
	if c.cant > 0 {
		total := float64(c.cant) * c.prices.Suggested
		c.prices.Suggested = (total - item.Price + newItem.Price) / float64(c.cant)
	}
}
