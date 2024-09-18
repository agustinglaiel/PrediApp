package main

import (
	"fmt"
	"prodes/internal/api"
	client "prodes/internal/client"
	"prodes/internal/repository"
	"prodes/internal/router"
	"prodes/internal/service"
	"prodes/pkg/utils"

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

	// Inicializar el cliente HTTP para comunicarte con el microservicio de sessions
	httpClient := client.NewHttpClient("http://localhost:8060/sessions")

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
	if err := r.Run(":8081"); err != nil {
		panic(err)
	}
}
