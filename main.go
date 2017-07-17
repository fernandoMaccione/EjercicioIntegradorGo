package main
import "github.com/gin-gonic/gin"
import "net/http"

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

func ejecutar (c *gin.Context) {

	//name := c.Param("categorie")
	//p := &Prices{"100", "2", "0"}

	cacheCategories := GetInstance()
	//c.String(http.StatusOK, "Categorias " + name)
	c.JSON(http.StatusOK, cacheCategories.getCategories())
}