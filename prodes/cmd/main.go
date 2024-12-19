package main

import (
	"fmt"
	"log"
	"os"
	"prodes/internal/api"
	client "prodes/internal/client"
	"prodes/internal/repository"
	"prodes/internal/router"
	"prodes/internal/service"
	"prodes/pkg/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Obtener el puerto de la variable de entorno PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // Valor por defecto en caso de que no esté configurado
	}

	// Inicializar la base de datos
	db, err := utils.InitDB()
	if err != nil {
		fmt.Println("Error al conectar con la Base de Datos")
		panic(err)
	}
	defer utils.DisconnectDB()

	// Iniciar el motor de la base de datos y migrar tablas y la caché
	utils.StartDbEngine()
	cache := utils.NewCache()

	// Inicializar el cliente HTTP para comunicarte con el microservicio de sessions
	httpClient := client.NewHttpClient("http://localhost:")

	// Inicializar repositorios, servicios y controlador
	prodeRepo := repository.NewProdeRepository(db)
	prodeService := service.NewPrediService(prodeRepo, httpClient, cache)
	prodeController := api.NewProdeController(prodeService)

	// Configurar el router
	ginRouter := gin.Default()

	// Configurar RabbitMQ y comenzar a consumir mensajes
	conn, ch, err := utils.SetupRabbitMQ()
	if err != nil {
		log.Fatalf("Error configurando RabbitMQ: %v", err)
	}
	defer utils.CloseRabbitMQ(conn, ch)

	// Iniciar la escucha de mensajes en segundo plano
	msgs, err := utils.ConsumeMessages(ch)
	if err != nil {
		log.Fatalf("Error al consumir mensajes de RabbitMQ: %v", err)
	}
	go utils.HandleMessage(msgs)

	// Llamar a MapUrls para configurar las rutas
	router.MapUrls(ginRouter, prodeController)

	// Iniciar el servidor usando el puerto obtenido de la variable de entorno
	// if err := ginRouter.Run(":" + port); err != nil {
	// 	log.Fatalf("Failed to run server on port %s: %v", port, err)
	// }
}
