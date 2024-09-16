package main

import (
	"prodes/internal/api"
	client "prodes/internal/client"
	"prodes/internal/repository"
	"prodes/internal/router"
	"prodes/internal/service"
	"prodes/pkg/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inicializar la base de datos
	db, err := utils.InitDB()
	if err != nil {
		panic(err)
	}

	// Inicializar el cliente HTTP para comunicarte con el microservicio de sessions
	httpClient := client.NewHttpClient("http://sessions-service-url")

	// Inicializar repositorios y servicios
	prodeRepo := repository.NewProdeRepository(db)
	prodeService := service.NewPrediService(prodeRepo, httpClient)

	// Inicializar el controlador
	prodeController := api.NewProdeController(prodeService)

	// Configurar el router
	r := gin.Default()

	// Llamar a MapUrls para configurar las rutas
	router.MapUrls(r, prodeController)  // Llama a la función de enrutamiento

	// Iniciar el servidor
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
