package categories

import "EjercicioIntegradorGo/config"

var FillPriceTotalItems FillPrice = func (categoria string)([][]Item, error){
	var calcularOffsetMT calculateOffset = func (regTotales int, limit int, offset int, porcenSample float32) (offsetR int, pageTotales int){

		regTotales = regTotales
		pageTotales = int(float32(regTotales/limit) * porcenSample / 100)
		if pageTotales != 0{
			offsetR = regTotales / pageTotales  + offset
		}
		return
	}
	conf := config.GetInstance()

	return findItems(categoria,0,conf.Limit,nil, "relevance", calcularOffsetMT, 0, conf.PorcentItems)
}
