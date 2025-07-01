package main

import (
	"fmt"
	"log"
	"os"

	"prediapp.local/db"
	"prediapp.local/drivers/internal/api"
	"prediapp.local/drivers/internal/client"
	"prediapp.local/drivers/internal/repository"
	"prediapp.local/drivers/internal/router"
	"prediapp.local/drivers/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Obtener el puerto de la variable de entorno PORT
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set in the environment")
	}

	// Inicializar la base de datos
	err := db.Init()
	if err != nil {
		fmt.Println("Error al conectar con la Base de Datos")
		panic(err)
	}
	defer db.DisconnectDB()

	// Crear el cliente HTTP para interactuar con la API externa
	externalAPIClient := client.NewHttpClient("https://api.openf1.org/v1/")

	// Inicializar repositorios y servicios
	driverRepo := repository.NewDriverRepository(db.DB)
	driverService := service.NewDriverService(driverRepo, externalAPIClient)
	driverController := api.NewDriverController(driverService)

	// Configurar el router
	ginRouter := gin.Default()

	// Mapear URLs
	router.MapUrls(ginRouter, driverController)

	// Iniciar servidor usando el puerto obtenido de la variable de entorno
	fmt.Printf("Drivers service listening on port %s...\n", port)
	if err := ginRouter.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server on port %s: %v", port, err)
	}
}
