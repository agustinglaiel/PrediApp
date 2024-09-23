package main

import (
	"drivers/internal/api"
	"drivers/internal/client"
	"drivers/internal/repository"
	"drivers/internal/router"
	"drivers/internal/service"
	"drivers/pkg/utils"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
    db, err := utils.InitDB()
    if err != nil {
        fmt.Println("Error al conectar con la Base de Datos")
        panic(err)
    }
    defer utils.DisconnectDB()

    // Start the database engine to migrate tables
    utils.StartDbEngine()

	// Crear el cliente HTTP para interactuar con la API externa
	externalAPIClient := client.NewHttpClient("https://api.openf1.org/v1/")

	// Inicializar repositorios y servicios
	driverRepo := repository.NewDriverRepository(db)
	driverService := service.NewDriverService(driverRepo, externalAPIClient)
	driverController := api.NewDriverController(driverService)

	// Set up router
    ginRouter := gin.Default()

    // Map URLs
    router.MapUrls(ginRouter, driverController)

	// Start server
    if err := ginRouter.Run(":8070"); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}