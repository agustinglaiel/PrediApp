package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	// 3) HTTP client y caché
	httpClient := client.NewHttpClient("http://localhost:")
	cache := utils.NewCache()

	// 4) Repos, servicio y controlador
	pRepo := repository.NewProdeRepository(db.DB)
	pService := service.NewPrediService(pRepo, httpClient, cache)
	pController := api.NewProdeController(pService)

	// 5) Router
	r := gin.Default()
	router.MapUrls(r, pController)

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
}
