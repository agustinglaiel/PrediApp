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
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
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

	// Inicializar repositorios y servicios
	driverRepo := repository.NewDriverRepository(db)
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