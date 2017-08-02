package categories

import (
	"EjercicioIntegradorGo/config"
	"strconv"
	"EjercicioIntegradorGo/library"
	"errors"
)

var FillPriceTotalItemsGoRutine FillPrice = func (categoria string)([][]Item, error){
	var calcularOffsetMT calculateOffset = func (regTotales int, limit int, offset int, porcenSample float32) (offsetR int, pageTotales int){

		regTotales = regTotales
		pageTotales = int(float32(regTotales/limit) * porcenSample / 100)
		if pageTotales != 0{
			offsetR = regTotales / pageTotales  + offset
		}
		return
	}
	conf := config.GetInstance()

	return findByGoRutine(categoria,0,conf.Limit,nil, "relevance", calcularOffsetMT, 0, conf.PorcentItems)
}

func findByGoRutine(category string, offset int, limit int, mItem[][] Item, orden string, f calculateOffset,page int, porcenSample float32)([][]Item, error){
	conf := config.GetInstance()
	url := conf.UrlSearch + category + "&offset=" + strconv.Itoa(offset)+ "&limit=1&sort=" + orden
	res := &ResponseSearch{}
	err := library.DoRequest(url, "GET", &res)
	if err != nil {return mItem, err}
	if res.Paging.Total == 0 {
		return mItem, errors.New("No hay registro en la categoria")
	}
	var pageTotales int
	offset, pageTotales = f (res.Paging.Total, limit, offset, porcenSample)
	mItem = make([][]Item, pageTotales +1)
	//Declaro un channel por pagina
	nCham := pageTotales +1
	if (pageTotales>conf.MaxGoRutine){
		nCham = conf.MaxGoRutine
	}

	var chans =  []chan ResponseSearch{}
	for i:=0; i<nCham; i++{
		tmp := make(chan ResponseSearch)
		chans = append(chans, tmp)
		go findItemsGoRutine(category, limit * i , limit, orden, tmp)
	}

	for i,ch := range chans{
		resp := <- ch
		mItem[i] = resp.Result
	}
	return mItem, nil
}

func findItemsGoRutine(category string, offset int, limit int, orden string, ch chan ResponseSearch){
	conf := config.GetInstance()
	url := conf.UrlSearch + category + "&offset=" + strconv.Itoa(offset)+ "&limit=" + strconv.Itoa(limit) + "&sort=" + orden
	res := &ResponseSearch{}
	library.DoRequest(url, "GET", &res)
	//if err != nil {return mItem, err}

	ch <- *res
}