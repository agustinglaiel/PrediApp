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

	"github.com/gin-gonic/gin"
)

func main() {
	// Obtener el puerto de la variable de entorno PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8060" // Valor por defecto en caso de que no esté configurado
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

	// Inicializar repositorio y servicio
	sessionRepo := repository.NewSessionRepository(db)
	sessionService := service.NewSessionService(sessionRepo, externalAPIClient)
	sessionController := api.NewSessionController(sessionService)

	// Configurar router
	ginRouter := gin.Default()

	// Mapear URLs
	router.MapUrls(ginRouter, sessionController)

	// Iniciar servidor usando el puerto obtenido de la variable de entorno
	if err := ginRouter.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server on port %s: %v", port, err)
	}
}
