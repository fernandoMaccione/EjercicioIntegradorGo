package main
import "github.com/gin-gonic/gin"
import (
	"net/http"
	/*"errors"
	"fmt"*/
	"EjercicioIntegradorGo/categories/cache"
	"EjercicioIntegradorGo/config"
)

func main() {
	router := gin.Default()
	gin.SetMode(config.GetInstance().GinMode)
	router.GET("/categories/:categories/price", getPrice)
	router.GET("/name/:name", hola)
	router.GET("/categories", consult)
//	router.GET("/refreshCache", refreshCache)

	router.Run(":80")
}


func hola (c *gin.Context){
	name := c.Param("name")
	c.String(http.StatusOK, "Hola!!!  " + name)
}
/*
func refreshCache (c *gin.Context){
	err := cache.Refresh(config.GetInstance())
	if (err!=nil) {
		c.JSON(http.StatusNotFound,  gin.H{"status": err.Error()})
	}else{
		c.String(http.StatusOK, "Refresco Exitoso")
	}
}
*/
func consult(c *gin.Context) {
	cacheCategories := cache.GetInstanceCache()
	c.JSON(http.StatusOK, cacheCategories.GetCategories())
}

func getPrice(c *gin.Context) {
	name := c.Param("categories")
	cacheCategories := cache.GetInstanceCache()
	cat, err := cacheCategories.GetCategory(name)
	if err!=nil{
		errorCategory(err, c, name)
	}else if result, err := cat.GetPrices(); err!=nil{
		errorCategory(err, c, name)
	}else {
		c.JSON(http.StatusOK, result)
	}
}
func errorCategory(err error, c *gin.Context, name string){
	c.JSON(http.StatusNotFound,  gin.H{"categories": name, "status": err.Error()})
}