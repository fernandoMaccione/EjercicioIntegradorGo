package categories

import "EjercicioIntegradorGo/config"

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

