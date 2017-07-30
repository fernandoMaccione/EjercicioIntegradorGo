package cache

import (
	"time"
	"EjercicioIntegradorGo/config"
	"log"
)

func refreshCache(){
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Ocurrio un error inesperado al refrescar el cache. Se elimino la instancia y se reestablecera automaticaente con la proxima llamada al servicio. error: ", err)
			CleanCache()
		}
	}()

	if conf, err := config.GetInstance(); err == nil {
		log.Println("Se inicio la terea de refresco de cache exitosamente")
		for {
			time.Sleep(conf.MinRefreshCache * time.Minute)
			log.Println("Ejecutando refresco de cache...")
			if err := refresh(conf); err != nil{
				log.Printf("Ocurrieron error al ejecutar al refrescar el cache: ", err)
			}
		}
	}else{
		log.Printf("No pudo inicializarce la tarea de refresco y limpieza de cache, su aplicaciónn colapsara.........", err)
	}
}

func refresh(config *config.Config) error{
	cache := GetInstanceCache()
	for _,v := range cache.cache{
		if  v.LastEntry.Add(time.Minute *config.MinOldEntry).Before(time.Now()) || !v.Initialized(){ //si por alguna razón tampoco se pudo inicializar, tambien la vuelo.
			cache.Remove(v)
		}else{
			if err := v.ValidateState(); err != nil{ //Si por alguna razon, no pude actualizar los preciso de la categoria, la elimino. Significa que esta corrupta.
				log.Printf("Ocurrieron error al refrescar la categoria " + v.Id + ". Error: ", err)
				cache.Remove(v)
			}
		}
	}
	return nil
}