package main
/*Este test esta recontra verde... pero bueno... al menos para jugar un rato sirvio
Lo que todavía no me queda claro es porque solo me calcula la covertura sobre la clase main... ¿?¿?¿?¿?¿?¿?¿?¿?¿
Falta terminar
 */
import (
	"os"
	"testing"
	"EjercicioIntegradorGo/config"
	"EjercicioIntegradorGo/categories/cache"
	"github.com/gin-gonic/gin"
	"time"
	"strings"
	"EjercicioIntegradorGo/categories"
	"strconv"
	"net/http"
	"EjercicioIntegradorGo/library"
)

func TestMain(m *testing.M) {
	/*Para testear bien esto no me queda mas remedio que pegarle a un api de test...
		UrlSearch:"https://api.mercadolibre.com/sites/MLA/search?categories=",
		UrlItem:"https://api.mercadolibre.com/items/"}
	 */
	router := gin.Default()

	router.GET("/search/:categories", getSearch)
	router.GET("/item/:item", hola)
	router.GET("/categories", consult)

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

	conf := config.GetInstance()
	if conf == nil{
		t.Fatalf("NO funciona la configración del api")
	}else{
		t.Logf("Se cargó exitosamente la configuración %d", conf)
	}

}

func TestCache (t *testing.T){

	cac := cache.GetInstanceCache()
	if cac == nil{
		t.Fatalf("No se puede inicilaizar el cache")
	}else{
		t.Logf("Se inicializo exitosamente el cache %d", cac)
	}

}

func TestPriceByRelevance (t *testing.T){

	conf := config.GetInstance()
	conf.UrlSearch = "http://localhost:9090/search/"
	conf.Limit = 5
	conf.PorcentItems = 100

	var result *categories.Prices
	err := library.DoRequest("http://localhost:9080/categories/MLA1000/price", "GET", result)

	if err!=nil{
		t.Logf("Resultado con error %v", err)
	}else {
		t.Logf("Resultado ok %v", result)
	}

}

func TestPriceByTotal (t *testing.T){

	cacheCategories := cache.GetInstanceCache()
	conf := config.GetInstance()
	conf.UrlSearch = "http://localhost:9090/search/"
	conf.Limit = 5
	conf.MethodFill = 0
	conf.PorcentItems = 100
	cat, err := cacheCategories.GetCategory("MLA1004")
	if err!=nil{
		t.Logf("Resultado con error %v", err)
	}else if result, err := cat.GetPrices(); err!=nil{
		t.Logf("Resultado con error %v", err)
	}else {
		t.Logf("Resultado ok %v", result)
	}

}

func TestPriceByRelevanceCatCacheada (t *testing.T){

	conf := config.GetInstance()
	conf.UrlSearch = "http://localhost:9090/search/"
	conf.Limit = 5
	conf.PorcentItems = 100

	var result *categories.Prices
	err := library.DoRequest("http://localhost:9080/categories/MLA1000/price", "GET", result)

	if err!=nil{
		t.Logf("Resultado con error %v", err)
	}else {
		t.Logf("Resultado ok %v", result)
	}

}

func TestSinCategoria (t *testing.T){

	var result *categories.Prices
	err := library.DoRequest("http://localhost:9080/categories/MLA100443430/price", "GET", result)

	if err!=nil{
		t.Logf("Resultado con error %v", err)
	}else {
		t.Logf("Resultado ok %v", result)
	}
}
/*
func TestRefrescarCategoria (t *testing.T){

	var result *categories.Prices
	err := library.DoRequest("http://localhost:9080/categories/MLA100443430/price", "GET", result)

	if err!=nil{
		t.Logf("Resultado con error %v", err)
	}else {
		t.Logf("Resultado ok %v", result)
	}
}
*/
func TestHolaApi (t *testing.T){

	var result *categories.Prices
	err := library.DoRequest("http://localhost:9080/name/Api", "GET", result)

	if err!=nil{
		t.Logf("Resultado con error %v", err)
	}else {
		t.Logf("Resultado ok %v", result)
	}
}

func TestConsultarCache (t *testing.T){

	var result *categories.Prices
	err := library.DoRequest("http://localhost:9080/categories", "GET", result)

	if err!=nil{
		t.Logf("Resultado con error %v", err)
	}else {
		t.Logf("Resultado ok %v", result)
	}
}

func TestRefreshCache (t *testing.T){
	t.Logf("mapa %d", cache.GetInstanceCache().GetCategories())
	err := cache.Refresh(config.GetInstance())
	if err != nil {
		t.Fatalf("Ocurrieron errores al refrescar el cache %v", err)
	}else{
		t.Logf("Se refrescó el cache exitosamente")
	}
}