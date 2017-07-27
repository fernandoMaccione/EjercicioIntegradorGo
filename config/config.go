package config

import (
	"sync"
	"EjercicioIntegradorGo/categories"
	"sync/atomic"
)

type Config struct {
	MethodFill categories.FillPrice
	PorcentItems float32
	Limit int
}
var initialized uint32
var mu sync.Mutex
var config *Config
func GetInstance() *Config {

	if atomic.LoadUint32(&initialized) == 1 {
		return config
	}

	mu.Lock()
	defer mu.Unlock()

	if initialized == 0 {
		if fillCache() {
			atomic.StoreUint32(&initialized, 1)
		}
	}

	return config
}

func fillCache(){
	//En la versión 2 hago que esto se levante de un archivo de configuración.
	config = &Config{MethodFill:categories.FillPreciosByRelevancia, PorcentItems:5, Limit:100}
 }