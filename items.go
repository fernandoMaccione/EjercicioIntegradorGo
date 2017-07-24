package main
import "net/http"
import (
	"log"
	"encoding/json"
	"strconv"
)


type Item struct {
	Id string `json:"id"`
	Price float64 `json:"price"`
}

type Paging struct{
	Total int32 `json:"total"`
}

type Respuesta struct{
	Paging Paging `json:"paging"`
	Site_id string `json:"site_id"`
	Result []Item `json:"results"`
}

func fillPrecios (categoria string, offset int, limit int, mItem[][] Item)([][]Item, error){
	url := "https://api.mercadolibre.com/sites/MLA/search?category=" + categoria + "&offset=" + strconv.Itoa( offset )+ "&limit=" + strconv.Itoa(limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return mItem, err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return mItem, err
	}

	defer resp.Body.Close()

	res := &Respuesta{}
	err = json.NewDecoder(resp.Body).Decode(&res)

	if mItem == nil{
		mItem = make([][]Item, res.Paging.Total/100)
	}

	mItem[offset/100] = res.Result

	if (len(res.Result)> 0 && int32(offset) < res.Paging.Total){
		return fillPrecios(categoria, offset + 100, limit, mItem)
	}
	return mItem, nil
}