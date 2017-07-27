package main
import "sync"
import (
	"sync/atomic"
	"EjercicioIntegradorGo/categories"
	"EjercicioIntegradorGo/library"
)

var initialized uint32
var mu sync.Mutex


type CacheCategories struct{
	cache map[string]*categories.Category
}

func (c *CacheCategories) add (categoria *categories.Category){
	if c.cache == nil{
		c.cache = make(map[string]*categories.Category)
	}
	c.cache[categoria.Id] = categoria
}

func (c *CacheCategories) remove (categoria *categories.Category){
	if c.cache != nil{
		delete(c.cache,categoria.Id)
	}
}

func (c *CacheCategories) getCategory(key string)(*categories.Category, error) {

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

func (c *CacheCategories) getCategories()map[string]*categories.Category{
	return c.cache
}

func (c *CacheCategories) contains (key string) bool{
	_,existe := c.cache[key]
	if existe {
		return true
	}
	return false
}

var cache *CacheCategories
func GetInstanceCache() *CacheCategories {

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

	cache = &CacheCategories{}

	//for _, v := range vecCat {
	for i:=0; i<len(vecCat); i++{
		cache.add(&vecCat[i])
	}

	return true
}