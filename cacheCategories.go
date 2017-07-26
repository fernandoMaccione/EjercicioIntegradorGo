package main
import "sync"
import (
	"sync/atomic"
	"proyecto1/categories"
	"proyecto1/library"
)

var initialized uint32
var mu sync.Mutex


type cacheCategories struct{
	cache map[string]*categories.Category
}

func (c *cacheCategories) add (categoria *categories.Category){
	if c.cache == nil{
		c.cache = make(map[string]*categories.Category)
	}
	c.cache[categoria.Id] = categoria
}

func (c *cacheCategories) remove (categoria *categories.Category){
	if c.cache != nil{
		delete(c.cache,categoria.Id)
	}
}

func (c *cacheCategories) getCategory(key string)(*categories.Category, error) {

	cat,exist := c.cache[key]
	if exist {
		return cat, nil
	}else{
		mu.Lock()
		defer mu.Unlock()
		cat,exist := c.cache[key]
		if (!exist) { //is ya existe, es porque justo la petición anterior que bloqueó el proceso la creo entonces devuevlo esa
			cat = &categories.Category{Id: key}
			c.add(cat)
			return cat, nil
		}else {
			return cat, nil
		}
	}
}

func (c *cacheCategories) getCategories()map[string]*categories.Category{
	return c.cache
}

func (c *cacheCategories) contains (key string) bool{
	_,existe := c.cache[key]
	if existe {
		return true
	}
	return false
}

var cache *cacheCategories
func GetInstanceCache() *cacheCategories {

	if atomic.LoadUint32(&initialized) == 1 {
		return cache
	}

	mu.Lock()
	defer mu.Unlock()

	if initialized == 0 {
		if fillCache() {
			atomic.StoreUint32(&initialized, 1)
		}
	}

	return cache
}

func fillCache() bool{

	url := "https://api.mercadolibre.com/sites/MLA/categories"
	vecCat := make([]categories.Category, 0, 31)
	err := library.DoRequest(url, "GET",&vecCat)
	if err != nil {
		return false
	}

	cache = &cacheCategories{}

	//for _, v := range vecCat {
	for i:=0; i<len(vecCat); i++{
		cache.add(&vecCat[i])
	}

	return true
}