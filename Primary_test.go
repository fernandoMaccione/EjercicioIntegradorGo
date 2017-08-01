package main

import (
	"os"
	"testing"
	"EjercicioIntegradorGo/config"
	"github.com/gin-gonic/gin"
	"time"
	"strings"
	"EjercicioIntegradorGo/categories"
	"strconv"
	"net/http"

	"EjercicioIntegradorGo/library"
	"EjercicioIntegradorGo/categories/cache"
)

func TestMain(m *testing.M) {
	/*Para testear bien esto no me queda mas remedio que pegarle a un api de test...
		UrlSearch:"https://api.mercadolibre.com/sites/MLA/search?categories=",
		UrlItem:"https://api.mercadolibre.com/items/"}
		UrlCategoria https://api.mercadolibre.com/categories/
	 */
	router := gin.Default()
	gin.SetMode(config.GetInstance().GinMode)
	router.GET("/search/:categories", getSearch)
	router.GET("/item/:item", getItem)
	router.GET("/categories/:category", getCategory)

	items = make([]categories.Item, 0, 20)
	items = append(items, categories.Item{Id:"MLA1232", Price:10})
	items = append(items, categories.Item{Id:"MLA1233", Price:1333})
	for i := 2; i<20; i++{
		items = append(items, categories.Item{Id:"MLA1232" + strconv.Itoa(i), Price:200})
	}

	go router.Run(":9090")
	go main()
	time.Sleep(2000)
	i := m.Run()
	os.Exit(i)
}
var items []categories.Item

func getItem(c *gin.Context){
	keyItem := c.Param("item")
	for _ ,v := range items{
		if v.Id == keyItem{
			c.JSON(http.StatusOK, v)
			break
		}
	}
}

func getCategory(c *gin.Context){
	name := c.Param("category")
	params := strings.Split(name, "?")
	category := params[0]

	if category == "MLA1000" || category == "MLA1004" || category == "MLA9999"{
		c.JSON(http.StatusOK, &categories.Category{Id:category})
	}else{
		c.JSON(http.StatusOK, &categories.Category{Id:""})
	}
}

func getSearch(c *gin.Context){
	name := c.Param("categories")
	params := strings.Split(name, "&")
	category := params[0]
	offset,_ := strconv.Atoi(strings.Split(params[1],"=")[1])

	limit,_ := strconv.Atoi(strings.Split(params[2],"=")[1])

	pag := &categories.Paging{len(items)}
	var cat1 *categories.ResponseSearch
	if limit > len(items){limit = len(items)}
	if category == "MLA1000" || category == "MLA1004" {
		cat1 = &categories.ResponseSearch{Paging: *pag, Site_id: "MLA", Result: items[offset: offset+limit]}
	}else{
		cat1 = &categories.ResponseSearch{Paging: categories.Paging{Total:0}, Site_id: "MLA"}
	}
	c.JSON(http.StatusOK, cat1)
}

func TestConfig (t *testing.T){
	t.Logf("----------------------Ejecutando testeo de configuracion...----------------------")
	conf := config.GetInstance()
	if conf == nil{
		t.Fatalf("NO funciona la configración del api")
	}else{
		conf.Limit = 5
		if conf.Limit == 5{
			t.Logf("Se cargó exitosamente la configuración %v", conf)
		}else{
			t.Fatalf("NO funciona la configración del api")
		}
	}
	t.Logf("----------------------Testeo de configuracion Ok----------------------")
}

func TestCache (t *testing.T){
	t.Logf("----------------------Iniciando testeo de cache...----------------------")
	cac := cache.GetInstanceCache()
	if cac == nil{
		t.Fatalf("No se puede inicilaizar el cache")
	}else{
		cac.Add(&categories.Category{Id:"MLA123456789"})
		aux, err := cac.GetCategory("MLA123456789")
		if err == nil && aux.Id == "MLA123456789"{
			t.Logf("Se inicializo exitosamente el cache %d", aux)
		}else{
			t.Fatalf("Ocurrieron errores al leer el cache. %v", err)
		}
	}
	t.Logf("----------------------Testeo de configuracion Ok----------------------")
}

func TestPriceByRelevance (t *testing.T){
	t.Logf("----------------------Iniciando testeo de calculo por relevancia...----------------------")
	conf := config.GetInstance()
	conf.UrlSearch = "http://localhost:9090/search/"
	conf.UrlCategory = "http://localhost:9090/categories/"
	conf.Limit = 5
	conf.MethodFill = 1
	conf.PorcentItems = 100

	cat := &categories.Category{Id:"MLA1000"}
	result, err := cat.GetPrices()

	if err!=nil{
		t.Fatalf("Resultado con error %v", err)
	}else {
		if (result.Max == 1333 && result.Min == 10 && int(result.Suggested) == 225){
			t.Logf("Resultado ok %v", result)
		}else{
			t.Fatalf("El resultado del caluclo no es el esperado %v", result)
		}
	}
	t.Logf("----------------------testeo de calculo por relevancia Ok----------------------")
}

func TestPriceByTotal (t *testing.T){
	t.Logf("----------------------Iniciando testeo de calculo por total...----------------------")
	conf := config.GetInstance()
	conf.UrlSearch = "http://localhost:9090/search/"
	conf.UrlCategory = "http://localhost:9090/categories/"
	conf.Limit = 5
	conf.MethodFill = 0
	conf.PorcentItems = 100

	cat := &categories.Category{Id:"MLA1000"}
	result, err := cat.GetPrices()

	if err!=nil{
		t.Fatalf("Resultado con error %v", err)
	}else {
		if (result.Max == 1333 && result.Min == 10 && int(result.Suggested) == 247){
			t.Logf("Resultado ok %v", result)
		}else{
			t.Fatalf("El resultado del caluclo no es el esperado %v", result)
		}
	}
	t.Logf("----------------------testeo de calculo por total Ok----------------------")
}

func TestRefreshPricePartial (t *testing.T){
	t.Logf("----------------------Iniciando testeo de Refresco parcial del precio de los items----------------------")
	conf := config.GetInstance()
	conf.UrlSearch = "http://localhost:9090/search/"
	conf.UrlItem = "http://localhost:9090/item/"
	conf.UrlCategory = "http://localhost:9090/categories/"
	conf.Limit = 5
	conf.MethodFill = 0
	conf.PorcentItems = 100
	conf.MinUpdatePartial = 0 // lo pongo en 0 para que asuma que el precio que está cacheado en el calculo ya es obsoleto
	cat := &categories.Category{Id:"MLA1000"}
	cat.GetPrices()
	auxI := &items[5]
	auxI.Price = 1500
	auxI.Last_Update = time.Now()
	auxI = &items[7]
	auxI.Price = 5
	auxI.Last_Update = time.Now()

	time.Sleep(100) // dejo pasar unos milisegundos para que se invalide el precio
	err :=cat.ValidateState()
	if err!=nil{
		t.Fatalf("Error al validad el precio %v", err)
	}else {
		r,err := cat.GetPrices()
		if err != nil{
			t.Fatalf("Error al devolver precio cacheado %v", err)
		}else {
			if (r.Min == 5 && r.Max == 1500 && int (r.Suggested) == 302) {
				t.Logf("Resultado ok %v", r)
			}
		}

	}
	t.Logf("---------------------testeo de Refresco parcial del precio de los items Ok----------------------")
}

func TestRefreshPriceTotal (t *testing.T){
	t.Logf("----------------------Iniciando testeo de Refresco total del precio de los items----------------------")
	conf := config.GetInstance()
	conf.UrlSearch = "http://localhost:9090/search/"
	conf.UrlItem = "http://localhost:9090/item/"
	conf.UrlCategory = "http://localhost:9090/categories/"
	conf.Limit = 5
	conf.MethodFill = 0
	conf.PorcentItems = 100
	conf.HourUpdateTotal = 0 // lo pongo en 0 para que asuma que el precio que está cacheado en el calculo ya es obsoleto
	cat := &categories.Category{Id:"MLA1000"}
	cat.GetPrices()
	auxI := &items[5]
	auxI.Price = 1500
	auxI.Last_Update = time.Now()
	auxI = &items[7]
	auxI.Price = 5
	auxI.Last_Update = time.Now()
	items = append(items, categories.Item{Id:"MLA6232", Price:2000})//al ser un refresco total, deberia enterarme si hay algun item nuevo
	items = append(items, categories.Item{Id:"MLA6233", Price:2000})//al ser un refresco total, deberia enterarme si hay algun item nuevo
	items = append(items, categories.Item{Id:"MLA6234", Price:2000})//al ser un refresco total, deberia enterarme si hay algun item nuevo
	items = append(items, categories.Item{Id:"MLA6235", Price:2000})//al ser un refresco total, deberia enterarme si hay algun item nuevo
	items = append(items, categories.Item{Id:"MLA6236", Price:2000})//al ser un refresco total, deberia enterarme si hay algun item nuevo

	time.Sleep(100) // dejo pasar unos milisegundos para que el precio quede invalido
	err :=cat.ValidateState()
	if err!=nil{
		t.Fatalf("Error al validad el precio %v", err)
	}else {
		r,err := cat.GetPrices()
		if err != nil{
			t.Fatalf("Error al devolver precio cacheado %v", err)
		}else {
			if (r.Min == 5 && r.Max == 2000 && int (r.Suggested) == 641) {
				t.Logf("Resultado ok %v", r)
			}else{
				t.Fatalf("El calculo no fue el esperado %v", r)
			}
		}

	}
	t.Logf("---------------------testeo de Refresco total del precio de los items Ok----------------------")
}

func TestRefreshCache (t *testing.T){
	t.Logf("----------------------Iniciando testeo de Refresco total del precio de los items----------------------")
	conf := config.GetInstance()
	conf.UrlSearch = "http://localhost:9090/search/"
	conf.UrlItem = "http://localhost:9090/item/"
	conf.UrlCategory = "http://localhost:9090/categories/"
	conf.Limit = 5
	conf.MethodFill = 0
	conf.PorcentItems = 100
	conf.MinOldEntry = 10

	ca := cache.GetInstanceCache()

	cat, err := ca.GetCategory("MLA1000")
	if (err != nil){
		t.Fatalf("Error al pedir la categoria %v", err)
	}
	_, err = cat.GetPrices()
	if (err != nil){
		t.Fatalf("Error al calcular el precio %v", err)
	}

	cat, err = ca.GetCategory("MLA1004")
	if (err != nil){
		t.Fatalf("Error al pedir la categoria %v", err)
	}
	_, err = cat.GetPrices()
	if (err != nil){
		t.Fatalf("Error al calcular el precio %v", err)
	}
	cat.LastEntry = time.Now().Add(time.Hour * -1)

	// tambien agrego bazura que deberia limpiarse
	ca.Add(&categories.Category{Id:"88888888"})

	t.Logf("elemantos cacheados ok %v", ca.GetCategories())
	err =cache.Refresh(config.GetInstance())
	if err!=nil{
		t.Fatalf("Error al refrescar el cache %v", err)
	}else {

		if ca.Contains("88888888"){
			t.Fatalf("El caché no está eliminando categorias que no se pueden calcular %v", err)
		}
		if ca.Contains("MLA1004"){
			t.Fatalf("El caché no está eliminando categorias que no tienen actividad %v", err)
		}
		if !ca.Contains("MLA1000"){
			t.Fatalf("El caché esta eliminando categorias que está vigentes %v", err)
		}
		t.Logf("Resultado ok %v", ca.GetCategories())
	}
	t.Logf("---------------------testeo de Refresco total del precio de los items Ok----------------------")
}

func TestCleanCache (t *testing.T){
	t.Logf("----------------------Iniciando testeo de vaciado total del cache----------------------")

	ca := cache.GetInstanceCache()

	_, err := ca.GetCategory("MLA1000")
	if (err != nil){
		t.Fatalf("Error al pedir la categoria %v", err)
	}
	_, err = ca.GetCategory("MLA1004")
	if (err != nil){
		t.Fatalf("Error al pedir la categoria %v", err)
	}
	ca.Add(&categories.Category{Id:"88888888"})

	t.Logf("elemantos cacheados ok %v", ca.GetCategories())

	cache.CleanCache()
	ca = cache.GetInstanceCache() //lo vuelvo a pedir y verifico que esté vacio.
	t.Logf("elemantos cacheados despues del vaciado total %v", ca.GetCategories())
	if len(ca.GetCategories()) >0 {
		t.Fatalf("Error al vaciar el cache %v")
	}
	t.Logf("---------------------testeo de vaciado total del cache Ok----------------------")
}

func TestBlackGetPrice (t *testing.T){

	t.Logf("----------------------Iniciando testeo de caja negra del metodo Price----------------------")
	var r = &categories.Prices{}
	err := library.DoRequest("http://localhost:9080/categories/MLA1000/price", "GET", r)

	if err != nil{
		t.Fatalf("El metodo respondio con error %v", err)
	}else {
		if (r.Min == 5 && r.Max == 2000 && int (r.Suggested) == 641) {
			t.Logf("Resultado ok %v", r)
		}else{
			t.Fatalf("El calculo no fue el esperado %v", r)
		}
	}
	t.Logf("----------------------testeo de caja negra del metodo Price Ok----------------------")
}

func TestBlackGetPriceCache (t *testing.T){
	t.Logf("----------------------Iniciando testeo de caja negra del metodo Price pidiendo un dato cacheado----------------------")
	var r = &categories.Prices{}
	err := library.DoRequest("http://localhost:9080/categories/MLA1000/price", "GET", r)

	if err != nil{
		t.Fatalf("El metodo respondio con error %v", err)
	}else {
		if (r.Min == 5 && r.Max == 2000 && int (r.Suggested) == 641) {
			t.Logf("Resultado ok %v", r)
		}else{
			t.Fatalf("El calculo no fue el esperado %v", r)
		}
	}
	t.Logf("----------------------testeo de caja negra del metodo Price pidiendo un dato cachead Ok----------------------")
}

func TestBlackOutCategory (t *testing.T){
	t.Logf("----------------------Iniciando testeo de caja negra del metodo Price pidiendo una categoria inexistente----------------------")
	var result = &categories.Prices{}
	err := library.DoRequest("http://localhost:9080/categories/MLA100443430/price", "GET", result)

	if err!=nil{
		t.Fatalf("Resultado con error %v", err)
	}else {
		t.Logf("Resultado ok %v", result)
	}
	t.Logf("----------------------Fin testeo de caja negra del metodo Price pidiendo una categoria inexistente----------------------")
}

func TestBlackCategoryOutItems (t *testing.T){
	t.Logf("----------------------Iniciando testeo de caja negra del metodo Price pidiendo una categoria sin items----------------------")
	var result = &categories.Prices{}
	err := library.DoRequest("http://localhost:9080/categories/MLA9999/price", "GET", result)

	if err!=nil{
		t.Fatalf("Resultado con error %v", err)
	}else {
		t.Logf("Resultado ok %v", result)
	}
	t.Logf("----------------------Fin testeo de caja negra del metodo Price pidiendo una categoria sin items----------------------")
}

func TestHolaApi (t *testing.T){

	var result = &categories.Prices{}
	library.DoRequest("http://localhost:9080/name/Api", "GET", result)

}

func TestConsult (t *testing.T){

	var result = &categories.Prices{}
	library.DoRequest("http://localhost:9080/categories", "GET", result)

}
