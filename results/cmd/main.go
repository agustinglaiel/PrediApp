package main

import (
	"fmt"
	"log"
	"os"

	"prediapp.local/db"
	"prediapp.local/results/internal/api"
	"prediapp.local/results/internal/client"
	"prediapp.local/results/internal/repository"
	"prediapp.local/results/internal/router"
	"prediapp.local/results/internal/service"

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
	externalAPIClient := client.NewHttpClient("http://localhost:8080")

	// Inicializar repositorio y servicio
	resultRepo := repository.NewResultRepository(db.DB)
	resultService := service.NewResultService(resultRepo, externalAPIClient)
	resultController := api.NewResultController(resultService)

	// Configurar router
	ginRouter := gin.Default()

	// Mapear URLs
	router.MapUrls(ginRouter, resultController)

	// Iniciar servidor usando el puerto obtenido de la variable de entorno
	fmt.Printf("Results service listening on port %s...\n", port)
	if err := ginRouter.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server on port %s: %v", port, err)
	}
}
