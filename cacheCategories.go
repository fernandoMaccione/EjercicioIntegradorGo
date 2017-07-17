package main
import "sync"
import "net/http"
import (
	"sync/atomic"
	"log"
	"encoding/json"
	"fmt"
)

var initialized uint32
var mu sync.Mutex

type Category struct {
	Id string
	Name string
}

type cacheCategories struct{
	cache map[string]Category
}

func (c *cacheCategories) add (categoria Category){
	if c.cache == nil{
		c.cache = make(map[string]Category)
	}
	c.cache[categoria.Id] = categoria
}

func (c *cacheCategories) remove (categoria Category){
	if c.cache != nil{
		delete(c.cache,categoria.Id)
	}
}

func (c *cacheCategories) getCategories()map[string]Category{
	return c.cache
}

func (c *cacheCategories) contains (key string) bool{
	if c.cache != nil{
		_,existe := c.cache[key]
		if existe {
			return true
		}
	}
	return false
}

var cache *cacheCategories
func GetInstance() *cacheCategories {

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

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return false
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return false
	}

	defer resp.Body.Close()

	vecCat := make([]Category, 0, 31)

	err = json.NewDecoder(resp.Body).Decode(&vecCat)
	if (err!= nil){
		return false
	}
	fmt.Println(vecCat)

	cache = &cacheCategories{}

	for _, v := range vecCat {
		cache.add(v)
	}

	return true
}