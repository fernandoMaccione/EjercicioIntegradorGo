package config

import (
	"time"
	"log"
	"github.com/gin-gonic/gin"
)

type Config struct {
	MethodFill       int
	PorcentItems     float32
	Limit            int
	MinUpdatePartial time.Duration
	HourUpdateTotal  time.Duration
	MinRefreshCache	time.Duration
	MinOldEntry time.Duration
	UrlSearch string
	UrlItem string
	UrlCategory string
	MaxGoRutine int
	GinMode string
}

var conf *Config

func init() {
	log.Println("Leyendo configuracion...")
	if err:=fillConfig(); err!= nil {
		log.Fatalf("No fue posible leventar el servicio. Ocurrieron errores al leer la configuracion. ", err)
	}
}
func GetInstance() (*Config) {
	return conf
}

func fillConfig()(error){
	//En la versión 2 hago que esto se levante de un archivo de configuración.
	conf = &Config{MethodFill: 2, PorcentItems: 50, Limit: 100,
		MinUpdatePartial: 60, HourUpdateTotal: 12, MinOldEntry:120,
		MinRefreshCache:10,
		UrlSearch:"https://api.mercadolibre.com/sites/MLA/search?category=",
		UrlItem:"https://api.mercadolibre.com/items/",
		UrlCategory:"https://api.mercadolibre.com/categories/",
		MaxGoRutine:100,
		GinMode:gin.ReleaseMode}
	return  nil
 }