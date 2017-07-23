package main
import "net/http"
import (
	"log"
	"encoding/json"
)


type Item struct {
	Id string `json:"id"`
	Price float64 `json:"price"`
}

type Paging struct{
	Total int32 `json:"total"`
	Offset int32 `json:"offset"`
	Limit int32 `json:"limit"`
}

type Respuesta struct{
	Paging Paging `json:"paging"`
	Site_id string `json:"site_id"`
	Result []Item `json:"results"`
}

func listPrecios (categoria string)([]Item){

	url := "https://api.mercadolibre.com/sites/MLA/search?category=" + categoria

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return nil
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return nil
	}

	defer resp.Body.Close()

	res := &Respuesta{}
	res.Result = make([]Item, 0, 50)

	err = json.NewDecoder(resp.Body).Decode(&res)

	return res.Result
}