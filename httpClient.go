package main
import (
	"net/http"
	"fmt"
	"log"
)

func doRequest(url string, method string) (*http.Response, error){
	fmt.Println(method + url )
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return nil , err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return nil, err
	}
	return resp, nil
}
