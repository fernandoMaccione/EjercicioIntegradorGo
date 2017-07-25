package main
import "net/http"
import (
	"log"
	"encoding/json"
	"strconv"
	"fmt"
)


type Item struct {
	Id string `json:"id"`
	Price float64 `json:"price"`
}

type Paging struct{
	Total int `json:"total"`
}

type Respuesta struct{
	Paging Paging `json:"paging"`
	Site_id string `json:"site_id"`
	Result []Item `json:"results"`
}

func fillPreciosPorMuestraTotal(categoria string, offset int, limit int, mItem[][] Item, page int)([][]Item, error){
	url := "https://api.mercadolibre.com/sites/MLA/search?category=" + categoria + "&offset=" + strconv.Itoa(offset * page)+ "&limit=" + strconv.Itoa(limit)
	fmt.Println("page: " + strconv.Itoa(page) + url )
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
	res := &Respuesta{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	resp.Body.Close()
	if mItem == nil{
		var pageTotales int
		offset, pageTotales = calcularOffset (res.Paging.Total - limit, limit)
		mItem = make([][]Item, pageTotales +1)
	}

	mItem[page - 1] = res.Result

	if (len(res.Result)> 0 && offset * page < res.Paging.Total - limit){
		return fillPreciosPorMuestraTotal(categoria, offset, limit, mItem, page + 1)
	}
	return mItem, nil
}

func calcularOffset(regTotales int, limit int) (offset int, pageTotales int){
	porcentajeMuestreo := 5 //va a ser configurable

	pageTotales = (regTotales/limit + 1) * porcentajeMuestreo / 100
	if pageTotales == 0{
		offset = regTotales
	}else {
		offset = regTotales / pageTotales + 1
	}

	return
}

func fillPreciosPorRelevancia(categoria string, offset int, limit int, mItem[][] Item, page int)([][]Item, error){
	return nil,nil
}