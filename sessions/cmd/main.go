package main

import (
	"fmt"
	"log"
	"sessions/internal/api"
	"sessions/internal/client"
	"sessions/internal/repository"
	"sessions/internal/router"
	"sessions/internal/service"
	"sessions/pkg/utils"

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
	sessionRepo := repository.NewSessionRepository(db)
	sessionService := service.NewSessionService(sessionRepo, externalAPIClient)
	sessionController := api.NewSessionController(sessionService)

	// Set up router
    ginRouter := gin.Default()

    // Map URLs
    router.MapUrls(ginRouter, sessionController)

	// Start server
    if err := ginRouter.Run(":8060"); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}
