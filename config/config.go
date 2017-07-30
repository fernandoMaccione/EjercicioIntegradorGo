package config

import (
	"sync"
	"sync/atomic"
	"time"
)

type Config struct {
	MethodFill       int
	PorcentItems     float32
	Limit            int
	MinUpdatePartial time.Duration
	HourUpdateTotal  time.Duration
	MinRefreshCache	time.Duration
	MinOldEntry time.Duration
}
var initialized uint32
var mu sync.Mutex
var conf *Config
func GetInstance() (*Config,error) {

	if atomic.LoadUint32(&initialized) == 1 {
		return conf, nil
	}

	mu.Lock()
	defer mu.Unlock()

	if initialized == 0 {
		if err:=fillCache(); err ==nil {
			atomic.StoreUint32(&initialized, 1)
		}else{
			return nil,err
		}
	}

	return conf, nil
}

func fillCache()(error){
	//En la versión 2 hago que esto se levante de un archivo de configuración.
	conf = &Config{MethodFill: 1, PorcentItems: 5, Limit: 100, MinUpdatePartial: 2, HourUpdateTotal: 12, MinOldEntry:120, MinRefreshCache:1}
	return  nil
 }