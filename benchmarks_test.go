package main

import (
	"testing"
	"EjercicioIntegradorGo/config"
	"EjercicioIntegradorGo/categories"
)

func BenchmarkOutCache(b *testing.B) {

	conf := config.GetInstance()
	conf.UrlSearch = "http://localhost:9090/search/"
	conf.UrlCategory = "http://localhost:9090/categories/"
	conf.Limit = 5
	conf.MethodFill = 2
	conf.MaxGoRutine = 50
	conf.PorcentItems = 100
	cat := &categories.Category{Id:"MLA1000"}
	b.ResetTimer()
	for i := 0; i < b.N*1000; i++ {
		cat.GetPrices()
	}
}

func BenchmarkInCache(b *testing.B) {

	conf := config.GetInstance()
	conf.UrlSearch = "http://localhost:9090/search/"
	conf.UrlCategory = "http://localhost:9090/categories/"
	conf.Limit = 5
	conf.MethodFill = 2
	conf.MaxGoRutine = 50
	conf.PorcentItems = 100
	cat := &categories.Category{Id:"MLA1000"}
	cat.GetPrices()
	b.ResetTimer()
	for i := 0; i < b.N*1000; i++ {
		cat.GetPrices()
	}
}