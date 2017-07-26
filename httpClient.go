package main
import (
	"net/http"
	"fmt"
	"log"
	"encoding/json"
)

func doRequest(url string, method string, v interface{}) (error){
	fmt.Println(method + url )
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return  err
	}

	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return  err
	}
	resp.Body.Close()
	return  nil
}
