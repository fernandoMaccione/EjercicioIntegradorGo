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
	router.GET("/categories/:categories/exist", checkCategory)
	router.GET("/name/:name", hola)
	router.GET("/categories", consult)

	router.Run(":9080")
}


func hola (c *gin.Context){
	name := c.Param("name")
	c.String(http.StatusOK, "Hola!!!  " + name)
}


func checkCategory (c *gin.Context) {

	cacheCategories := cache.GetInstanceCache()
	name := c.Param("categories")
	if cacheCategories.Contains(name){
		cat := cacheCategories.GetCategories()[name]
		c.JSON(http.StatusOK, cat)
	}else {
		c.JSON(http.StatusNotFound,  gin.H{"categories": name, "status": http.StatusNotFound})
	}
}

func consult(c *gin.Context) {
	cacheCategories := cache.GetInstanceCache()
	c.JSON(http.StatusOK, cacheCategories.GetCategories())
}

func getPrice(c *gin.Context) {

	name := c.Param("categories")

	cacheCategories := cache.GetInstanceCache()

	cat:= cacheCategories.GetCategory(name)
	result, err := cat.GetPrices()

	if err !=  nil{
		c.JSON(http.StatusNotFound,  gin.H{"categories": name, "status": err.Error()})
	}else {
		c.JSON(http.StatusOK, result)
	}
}
