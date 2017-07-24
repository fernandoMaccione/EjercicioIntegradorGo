package main
import "github.com/gin-gonic/gin"
import (
	"net/http"
	"errors"
	"fmt"
)

type Prices struct {
	Max string `json:"max"`
	Suggested string `json:"suggested"`
	Min string `json:"min"`
}


func main() {
	router := gin.Default()

	router.GET("/categories/:category/price", ejecutar)
	router.GET("/categories/:category/exist", checkCategory)
	router.GET("/name/:name", hola)
	router.GET("/categories", consult)

	router.Run(":9080")
}


func hola (c *gin.Context){
	name := c.Param("name")
	c.String(http.StatusOK, "Hola!!!  " + name)
}


func checkCategory (c *gin.Context) {

	cacheCategories := GetInstance()
	name := c.Param("category")
	if cacheCategories.contains(name){
		cat := cacheCategories.getCategories()[name]
		c.JSON(http.StatusOK, cat)
	}else {
		c.JSON(http.StatusNotFound,  gin.H{"category": name, "status": http.StatusNotFound})
	}
}

func consult(c *gin.Context) {
	cacheCategories := GetInstance()
	c.JSON(http.StatusOK, cacheCategories.getCategories())
}

func ejecutar (c *gin.Context) {

	name := c.Param("category")
	//p := &Prices{"100", "2", "0"}
	result, err := getPrice(name);

	if err !=  nil{
		c.JSON(http.StatusNotFound,  gin.H{"category": name, "status": err.Error()})
	}else {
		c.JSON(http.StatusOK, result)
	}
}

func getPrice(categoria string) (result *Prices, err error){
	cacheCategories := GetInstance()
	if cacheCategories.contains(categoria){
		result = &Prices{"100", "2", "0"}
		mItem, _ := fillPrecios(categoria, 0,2,nil)
		fmt.Printf("%+v\n", mItem)
	}else{
		err= errors.New("No exiset la categoria solicitada")
	}
	return
}