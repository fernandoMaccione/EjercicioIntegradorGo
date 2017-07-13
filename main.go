package main
import "github.com/gin-gonic/gin"
import "net/http"
import (
	"encoding/json"
	//"fmt"
	"fmt"
	"log"
	"io/ioutil"
)

type Prices struct {
	Max string `json:"max"`
	Suggested string `json:"suggested"`
	Min string `json:"min"`
}

func main() {
	router := gin.Default()

	router.GET("/categories/:categorie/price", ejecutar)

	router.Run(":9080")
}

func ejecutar (c *gin.Context) {

	name := c.Param("categorie")
	p := &Prices{"100", "2", "0"}

	res1B, _ := json.Marshal(p)

	url := "https://api.mercadolibre.com/sites/MLA/categories"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		fmt.Println (string(b))
	}

	c.String(http.StatusOK, "Categorias " + string(b))
	c.String(http.StatusOK, "Hello %s" + "Respuesta: " + string(res1B), name)


}