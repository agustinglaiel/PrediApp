package main

import (
	"admin/internal/api"
	"admin/internal/repository"
	"admin/internal/service"
	"admin/pkg/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inicializar la base de datos
	db, err := utils.InitDB()
	if err != nil {
		panic(err)
	}

	// Crear el cliente HTTP
	httpClient := utils.NewHttpClient("http://admin-service-url")

	// Crear repositorios y servicios
	eventRepo := repository.NewEventRepository(db)
	eventService := service.NewEventService(eventRepo, httpClient)

	// Crear el controlador
	eventController := api.NewEventController(eventService)

	// Configurar el router
	r := gin.Default()
	r.GET("/events/:id", eventController.GetEventByID)

	// Iniciar el servidor
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
