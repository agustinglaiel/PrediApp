package main

import (
	"fmt"
	"log"
	"os"
	"sessions/internal/api"
	"sessions/internal/client"
	"sessions/internal/repository"
	"sessions/internal/router"
	"sessions/internal/service"
	"sessions/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	time.Local = time.UTC // Establecer la zona horaria a UTC
	// Obtener el puerto de la variable de entorno PORT
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set in the environment")
	}

	// Inicializar la base de datos
	db, err := utils.InitDB()
	if err != nil {
		fmt.Println("Error al conectar con la Base de Datos")
		panic(err)
	}
	defer utils.DisconnectDB()

	// Iniciar el motor de la base de datos y migrar tablas
	utils.StartDbEngine()

	// Crear el cliente HTTP para interactuar con la API externa
	externalAPIClient := client.NewHttpClient("https://api.openf1.org/v1/")

	// Crear la instancia de caché con expiración de 30 minutos y tamaño máximo de 100 entradas
	cache := utils.NewCache(30*time.Minute, 100)

	// Inicializar repositorio y servicio
	sessionRepo := repository.NewSessionRepository(db)
	sessionService := service.NewSessionService(sessionRepo, externalAPIClient, cache)
	sessionController := api.NewSessionController(sessionService)

	// Configurar router
	ginRouter := gin.Default()

	// Mapear URLs
	router.MapUrls(ginRouter, sessionController)

	// Iniciar servidor usando el puerto obtenido de la variable de entorno
	fmt.Printf("Sessions service listening on port %s...\n", port)
	if err := ginRouter.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server on port %s: %v", port, err)
	}
}