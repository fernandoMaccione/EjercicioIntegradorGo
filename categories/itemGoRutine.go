package categories

import (
	"EjercicioIntegradorGo/config"
	"strconv"
	"EjercicioIntegradorGo/library"
	"errors"
	"log"
)

var FillPriceItemsGoRutine FillPrice = func (category string)([][]Item, error){

	pageTotales, err := getTotalPage(category)
	if (err!= nil){
		return nil, err
	}

	mItem := make([][]Item, pageTotales +2)
	conf := config.GetInstance()
	
	nCham := pageTotales
	if (nCham>conf.MaxGoRutine){
		nCham = conf.MaxGoRutine
	}

	var chans =  []chan ResponseSearch{}
	for i:=0; i<nCham + 2; i++{
		tmp := make(chan ResponseSearch)
		chans = append(chans, tmp)
		switch i {
		case nCham : go findItemsGoRutine(category, 0, 1, "price_asc", tmp)
		case nCham +1: go findItemsGoRutine(category, 0, 1, "price_desc", tmp)
		default: go findItemsGoRutine(category, conf.Limit*i, conf.Limit, "relevance", tmp)
		}
	}

	for i,ch := range chans{
		resp, ok := <- ch
		if ok{
			mItem[i] = resp.Result
		}
	}
	return mItem, nil
}

func getTotalPage(category string)(int, error){
	conf := config.GetInstance()
	url := conf.UrlCategory + category + "?attributes=total_items_in_this_category"
	res := struct {Total_items int `json:"total_items_in_this_category"` }{}
	err := library.DoRequest(url, "GET", &res)
	if err != nil {return 0, err}
	if res.Total_items == 0 {
		return 0, errors.New("No hay registro en la categoria")
	}

	pageTotales := int(float32(res.Total_items/conf.Limit) * conf.PorcentItems / 100)
	if (pageTotales<1){
		return 1, nil
	}
	return pageTotales, nil
}

func findItemsGoRutine(category string, offset int, limit int, orden string, ch chan ResponseSearch){
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Error al buscar los item de la categoría: ", err)
		}
		close(ch)
	}()
	conf := config.GetInstance()
	url := conf.UrlSearch + category + "&offset=" + strconv.Itoa(offset)+ "&limit=" + strconv.Itoa(limit) + "&sort=" + orden
	res := &ResponseSearch{}
	err := library.DoRequest(url, "GET", &res)
	if err == nil {
		ch <- *res
	}
}