package categories
import (
	"strconv"
	"time"
	"EjercicioIntegradorGo/library"
	"EjercicioIntegradorGo/config"
	"errors"
)

type Item struct {
	Id string `json:"id"`
	Price float64 `json:"price"`
	Last_Update time.Time `json:"last_updated"`
}

type Paging struct{
	Total int `json:"total"`
}

type ResponseSearch struct{
	Paging Paging `json:"paging"`
	Site_id string `json:"site_id"`
	Result []Item `json:"results"`
}

type calculateOffset func(int, int, int, float32)(int, int)
type FillPrice func(string)([][]Item, error)


func findItems(category string, offset int, limit int, mItem[][] Item, orden string, f calculateOffset,page int, porcenSample float32)([][]Item, error){
	conf := config.GetInstance()
	url := conf.UrlSearch + category + "&offset=" + strconv.Itoa(offset)+ "&limit=" + strconv.Itoa(limit) + "&sort=" + orden
	res := &ResponseSearch{}
	err := library.DoRequest(url, "GET", &res)
	if err != nil {return mItem, err}
	if res.Paging.Total == 0 {
		return mItem, errors.New("No hay registro en la categoria")
	}
	var pageTotales int
	offset, pageTotales = f (res.Paging.Total, limit, offset, porcenSample)
	if mItem == nil{
		mItem = make([][]Item, pageTotales +1)
	}
	mItem[page] = res.Result
	page++
	if len(res.Result)> 0 &&  page < pageTotales{
		return findItems(category, offset, limit, mItem,orden, f, page, porcenSample)
	}
	return mItem, nil
}

func findItem(id string)(*Item, error){
	url := config.GetInstance().UrlItem+id+"?attributes=id,last_updated,price"
	res := &Item{}
	return res, library.DoRequest(url, "GET", &res)
}