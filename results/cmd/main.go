package main

import (
	"fmt"
	"log"
	"results/internal/api"
	"results/internal/client"
	"results/internal/repository"
	"results/internal/router"
	"results/internal/service"
	"results/pkg/utils"

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

	// Initialize repositories and services
	sessionRepo := repository.NewResultRepository(db)
	sessionService := service.NewResultService(sessionRepo, externalAPIClient)
	sessionController := api.NewResultController(sessionService)

	// Set up router
    ginRouter := gin.Default()

    // Map URLs
    router.MapUrls(ginRouter, sessionController)

	// Start server
    if err := ginRouter.Run(":8071"); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}
