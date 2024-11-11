package main

import (
	"fmt"
	"log"
	"os"
	"results/internal/api"
	"results/internal/client"
	"results/internal/repository"
	"results/internal/router"
	"results/internal/service"
	"results/pkg/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Obtener el puerto de la variable de entorno PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8055" // Valor por defecto en caso de que no esté configurado
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
	externalAPIClient := client.NewHttpClient("")

	// Inicializar repositorio y servicio
	resultRepo := repository.NewResultRepository(db)
	resultService := service.NewResultService(resultRepo, externalAPIClient)
	resultController := api.NewResultController(resultService)

	// Configurar router
	ginRouter := gin.Default()

	// Mapear URLs
	router.MapUrls(ginRouter, resultController)

	// Iniciar servidor usando el puerto obtenido de la variable de entorno
	if err := ginRouter.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server on port %s: %v", port, err)
	}
}
