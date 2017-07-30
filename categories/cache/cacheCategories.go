package cache
import "sync"
import (
	"sync/atomic"
	"EjercicioIntegradorGo/categories"
	"time"
	"EjercicioIntegradorGo/library"
	"errors"
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

func (c *CacheCategories) GetCategory(key string)(*categories.Category, error) {
	cat,exist := c.cache[key]
	if exist {
		cat.LastEntry = time.Now()
		return cat, nil
	}else{
		mu.Lock()
		defer mu.Unlock()
		cat,exist := c.cache[key]
		if !exist { //is ya existe, es porque justo la petición anterior que bloqueó el proceso la creo entonces devuevlo esa
			if err := verifyCategory(key); err == nil {
				cat = &categories.Category{Id: key, LastEntry:time.Now()}
				c.Add(cat)
				return cat, nil
			}else{
				return nil, err
			}
		}else {
			return cat, nil
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

func  CleanCache (){
	cache = nil
	atomic.StoreUint32(&initialized, 0)
}

func GetInstanceCache() *CacheCategories {
	if atomic.LoadUint32(&initialized) == 1 {
		return cache
	}
	mu.Lock()
	defer mu.Unlock()
	if initialized == 0 {
		cache = &CacheCategories{}
		atomic.StoreUint32(&initialized, 1)
		go refreshCache()
	}
	return cache
}

func verifyCategory (key string) (error){
	url := "https://api.mercadolibre.com/categories/"+key+"?attributes=id"
	res := &categories.Category{}
	if err := library.DoRequest(url, "GET", &res); err != nil || res.Id == ""{
		return errors.New("La categoria solicitada no existe")
	}
	return nil

}