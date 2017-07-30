package main
import "github.com/gin-gonic/gin"
import (
	"net/http"
	/*"errors"
	"fmt"*/
	"EjercicioIntegradorGo/categories/cache"
)

func main() {
	router := gin.Default()

	router.GET("/categories/:categories/price", getPrice)
	router.GET("/name/:name", hola)
	router.GET("/categories", consult)

	router.Run(":9080")
}


func hola (c *gin.Context){
	name := c.Param("name")
	c.String(http.StatusOK, "Hola!!!  " + name)
}


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