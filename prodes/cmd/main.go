package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	// 1) Validar variables de entorno
	required := []string{
		"PORT", "JWT_SECRET",
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASS", "DB_NAME",
		"SESSIONS_SERVICE_URL",
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

	sessionsURL := os.Getenv("SESSIONS_SERVICE_URL")
	usersURL := os.Getenv("USERS_SERVICE_URL")
	driversURL := os.Getenv("DRIVERS_SERVICE_URL")
	resultsURL := os.Getenv("RESULTS_SERVICE_URL")

	sessionClient := client.NewHttpClient(sessionsURL)
	userClient := client.NewHttpClient(usersURL)
	driverClient := client.NewHttpClient(driversURL)
	resultsClient := client.NewHttpClient(resultsURL)
	cache := utils.NewCache(30*time.Minute, 100)

	// 4) Repos, servicio y controlador
	pRepo := repository.NewProdeRepository(db.DB)
	pService := service.NewPrediService(pRepo, sessionClient, userClient, driverClient, resultsClient, cache)
	pCtrl := api.NewProdeController(pService)

	// 5) Router
	r := gin.Default()
	router.MapUrls(r, pCtrl)

	// 6) Servir
	port := os.Getenv("PORT")
	go func() {
		fmt.Printf("Prodes service en puerto %s...\n", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatalf("Fallo al correr en %s: %v", port, err)
		}
	}()

	// 7) Esperar señal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Deteniendo prodes service...")

	entries := cache.ListEntries()
	log.Println("Contenido de la caché de prodes al cerrar:")
	for _, entry := range entries {
		log.Printf("Clave: %s, Expiración: %s, Valor: %+v\n", entry.Key, entry.Expiration.Format(time.RFC3339), entry.Value)
	}
}
