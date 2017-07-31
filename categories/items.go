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

var FillPriceTotalItems FillPrice = func (categoria string)([][]Item, error){
	var calcularOffsetMT calculateOffset = func (regTotales int, limit int, offset int, porcenSample float32) (offsetR int, pageTotales int){

		regTotales = regTotales - limit
		pageTotales = int(float32(regTotales/limit) * porcenSample / 100)
		if pageTotales != 0{
			offsetR = regTotales / pageTotales + 1 + offset
		}
		return
	}
	conf := config.GetInstance()

	return findItems(categoria,0,conf.Limit,nil, "relevance", calcularOffsetMT, 0, conf.PorcentItems)
}
var FillPriceByRelevance FillPrice = func(categoria string)([][]Item, error){
	var calculateOffsetMAXMIN calculateOffset = func (regTotales int, limit int, offset int, porcentajeMuestreo float32) (offsetR int, pageTotales int){
		return 1,1
	}
	mItem := make([][]Item, 2)
	mItem, err :=  findItems(categoria,0,1,mItem, "price_asc", calculateOffsetMAXMIN, 0, 100) //Busco el maximo
	if err!=nil {return nil, err}
	mItem, err =  findItems(categoria,0,1,mItem, "price_desc", calculateOffsetMAXMIN, 1, 100) //Busco el maximo
	if err!=nil{return nil, err}

	var calculateOffsetREL calculateOffset = func (regTotales int, limit int, offset int, porcenSample float32) (offsetR int, pageTotales int){
		pageTotales = int(float32(regTotales/limit) * porcenSample / 100)
		offsetR = limit + offset
		return
	}
	var mItemRL [][]Item
	conf := config.GetInstance()
	mItemRL, err =  findItems(categoria,0,conf.Limit,nil, "relevance", calculateOffsetREL, 0, conf.PorcentItems) //Busco los mas relevantes
	mItemRL = append(mItemRL, mItem[0], mItem[1])
	return mItemRL, err
}
func findItems(category string, offset int, limit int, mItem[][] Item, orden string, f calculateOffset,page int, porcenSample float32)([][]Item, error){
	conf := config.GetInstance()
	url := conf.UrlSearch + category + "&offset=" + strconv.Itoa(offset)+ "&limit=" + strconv.Itoa(limit) + "&sort=" + orden
	res := &ResponseSearch{}
	err := library.DoRequest(url, "GET", &res)
	if err != nil {return mItem, err}
	if res.Paging.Total == 0 {return mItem, errors.New("No hay registro en la categoria")}
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
	url := "https://api.mercadolibre.com/items/"+id+"?attributes=id,last_updated,price"
	res := &Item{}
	return res, library.DoRequest(url, "GET", &res)
}