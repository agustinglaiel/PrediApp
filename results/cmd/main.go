package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"prediapp.local/db"
	"prediapp.local/results/internal/api"
	"prediapp.local/results/internal/client"
	"prediapp.local/results/internal/repository"
	"prediapp.local/results/internal/router"
	"prediapp.local/results/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1) Validar variables de entorno
	required := []string{
		"PORT", "JWT_SECRET",
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASS", "DB_NAME",
		"DRIVERS_SERVICE_URL",
	}
	for _, v := range required {
		if os.Getenv(v) == "" {
			log.Fatalf("%s no está definida", v)
		}
	}

	// 2) Inicializar DB
	if err := db.Init(); err != nil {
		log.Fatalf("db.Init failed: %v", err)
	}

	// 3) Crear clientes HTTP para cada microservicio
	usersClient := client.NewHttpClient(os.Getenv("USERS_SERVICE_URL"))
	driversClient := client.NewHttpClient(os.Getenv("DRIVERS_SERVICE_URL"))
	sessionsClient := client.NewHttpClient(os.Getenv("SESSIONS_SERVICE_URL"))
	externalClient := client.NewHttpClient(os.Getenv("OPEN_F1_API_URL"))

	// …otros como groupsClient, prodesClient…

	// 4) Repositorio, servicio y controlador
	rRepo := repository.NewResultRepository(db.DB)
	rService := service.NewResultService(rRepo, driversClient, sessionsClient, usersClient, externalClient)
	rController := api.NewResultController(rService)

	// 5) Router
	r := gin.Default()
	router.MapUrls(r, rController)

	// 6) Servir
	port := os.Getenv("PORT")
	go func() {
		fmt.Printf("Results service en puerto %s...\n", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatalf("Fallo al correr en %s: %v", port, err)
		}
	}()

	// 7) Esperar señal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Deteniendo results service...")
}
