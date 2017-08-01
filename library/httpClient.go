package library
import (
	"net/http"
	"fmt"
	"log"
	"encoding/json"
	"EjercicioIntegradorGo/config"
	"github.com/gin-gonic/gin"
)

func DoRequest(url string, method string, v interface{}) (error){
	if config.GetInstance().GinMode == gin.DebugMode {
		fmt.Println(method + url)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("NewRequest: ", err)
		return err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Do: ", err)
		return  err
	}

	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return  err
	}
	resp.Body.Close()
	return  nil
}
