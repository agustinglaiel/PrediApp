package main

import (
	"fmt"
	"log"
	"predicciones/internal/api"
	"predicciones/internal/repository"
	"predicciones/internal/router"
	"predicciones/internal/service"
	"predicciones/pkg/utils"

	"github.com/gin-gonic/gin"
)

func main(){
    // Inicializar la base de datos
    db, err := utils.InitDB()
    if err != nil {
        fmt.Println("Error al conectar con la Base de Datos")
        panic(err)
    }
    defer utils.DisconnectDB()

    // Iniciar el motor de la base de datos para migrar las tablas
    utils.StartDbEngine()

    // Inicializar repositorio y servicio
    prodeRepo := repository.NewProdeRepository(db)
    prediService := service.NewPrediService(prodeRepo)

    // Inicializar el controlador
    prediController := api.NewPrediController(prediService)

    // Configurar el router
    ginRouter := gin.Default()

    // Mapear las URLs
    router.MapUrls(ginRouter, prediController)

    // Iniciar el servidor
    if err := ginRouter.Run(":8070"); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}
