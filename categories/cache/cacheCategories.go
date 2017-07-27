package cache
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

func (c *CacheCategories) Add (category *categories.Category){
	if c.cache == nil{
		c.cache = make(map[string]*categories.Category)
	}
	c.cache[category.Id] = category
}

func (c *CacheCategories) Remove (category *categories.Category){
	delete(c.cache,category.Id)
}

func (c *CacheCategories) GetCategory(key string)(*categories.Category) {

	cat,exist := c.cache[key]
	if exist {
		return cat
	}else{
		mu.Lock()
		defer mu.Unlock()
		cat,exist := c.cache[key]
		if !exist { //is ya existe, es porque justo la petición anterior que bloqueó el proceso la creo entonces devuevlo esa
			cat = &categories.Category{Id: key}
			c.Add(cat)
			return cat
		}else {
			return cat
		}
	}
}

func (c *CacheCategories) GetCategories()map[string]*categories.Category{
	return c.cache
}

func (c *CacheCategories) Contains (key string) bool{
	_, exist := c.cache[key]
	if exist {
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
		cache.Add(&vecCat[i])
	}

	return true
}