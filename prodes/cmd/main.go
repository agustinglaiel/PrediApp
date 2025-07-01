package main

import (
	"fmt"
	"log"
	"os"

	"prediapp.local/db"
	"prediapp.local/prodes/internal/api"
	client "prediapp.local/prodes/internal/client"
	"prediapp.local/prodes/internal/repository"
	"prediapp.local/prodes/internal/router"
	"prediapp.local/prodes/internal/service"
	"prediapp.local/prodes/pkg/utils"

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
	cache := utils.NewCache()

	// Inicializar el cliente HTTP para comunicarte con el microservicio de sessions
	httpClient := client.NewHttpClient("http://localhost:")

	// Inicializar repositorios, servicios y controlador
	prodeRepo := repository.NewProdeRepository(db.DB)
	prodeService := service.NewPrediService(prodeRepo, httpClient, cache)
	prodeController := api.NewProdeController(prodeService)

	// Configurar el router
	ginRouter := gin.Default()

	// Llamar a MapUrls para configurar las rutas
	router.MapUrls(ginRouter, prodeController)

	// Iniciar servidor usando el puerto obtenido de la variable de entorno
	fmt.Printf("Prodes service listening on port %s...\n", port)
	if err := ginRouter.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server on port %s: %v", port, err)
	}
}
