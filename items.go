package main
import (
	"strconv"
	"time"
)

type Item struct {
	Id string `json:"id"`
	Price float64 `json:"price"`
	Last_Update time.Time
}

type Paging struct{
	Total int `json:"total"`
}

type Respuesta struct{
	Paging Paging `json:"paging"`
	Site_id string `json:"site_id"`
	Result []Item `json:"results"`
}

type calcularOffset func(int, int, int, float32)(int, int)
type fillPrice func(string)([][]Item, error)

var fillPreciosPorMuestraTotal fillPrice = func (categoria string)([][]Item, error){
	var calcularOffsetMT  calcularOffset = func (regTotales int, limit int, offset int, porcentajeMuestreo float32) (offsetR int, pageTotales int){

		regTotales = regTotales - limit
		pageTotales = int(float32(regTotales/limit) * porcentajeMuestreo / 100)
		if pageTotales != 0{
			offsetR = regTotales / pageTotales + 1 + offset
		}
		return
	}
	var porcentajeMuestra float32 = 5 //sera un parametro
	return fillPrecios(categoria,0,100,nil, "relevance", calcularOffsetMT, 0, porcentajeMuestra)
}

var fillPreciosPorRelevancia fillPrice = func(categoria string)([][]Item, error){
	var calcularOffsetMAXMIN  calcularOffset = func (regTotales int, limit int, offset int, porcentajeMuestreo float32) (offsetR int, pageTotales int){
		return 1,1
	}
	mItem := make([][]Item, 2)
	mItem, err :=  fillPrecios(categoria,0,1,mItem, "price_asc", calcularOffsetMAXMIN, 0, 100) //Busco el maximo
	if (err!=nil) {return nil, err}
	mItem, err =  fillPrecios(categoria,0,1,mItem, "price_desc", calcularOffsetMAXMIN, 1, 100) //Busco el maximo
	if (err!=nil) {return nil, err}

	var calcularOffsetREL  calcularOffset = func (regTotales int, limit int, offset int, porcentajeMuestreo float32) (offsetR int, pageTotales int){
		pageTotales = int(float32(regTotales/limit) * porcentajeMuestreo / 100)
		offsetR = limit + offset
		return
	}
	var mItemRL [][]Item
	mItemRL, err =  fillPrecios(categoria,0,100,nil, "relevance", calcularOffsetREL, 0, .1) //Busco los mas relevantes
	mItemRL = append(mItemRL, mItem[0], mItem[1])
	return mItemRL, err
}



func fillPrecios(categoria string, offset int, limit int, mItem[][] Item, orden string, f calcularOffset ,page int, porcentajeMuestra float32)([][]Item, error){
	url := "https://api.mercadolibre.com/sites/MLA/search?category=" + categoria + "&offset=" + strconv.Itoa(offset)+ "&limit=" + strconv.Itoa(limit) + "&sort=" + orden

	res := &Respuesta{}
	err := doRequest(url, "GET", &res)
	if err != nil {
		return mItem, err
	}

	var pageTotales int
	offset, pageTotales = f (res.Paging.Total, limit, offset, porcentajeMuestra)
	if mItem == nil{
		mItem = make([][]Item, pageTotales)
	}

	mItem[page] = res.Result
	page++
	if (len(res.Result)> 0 &&  page < pageTotales){
		return fillPrecios(categoria, offset, limit, mItem,orden, f, page, porcentajeMuestra)
	}
	return mItem, nil
}