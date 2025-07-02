package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"prediapp.local/db"
	"prediapp.local/drivers/internal/api"
	"prediapp.local/drivers/internal/client"
	"prediapp.local/drivers/internal/repository"
	"prediapp.local/drivers/internal/router"
	"prediapp.local/drivers/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1) Validar variables de entorno
	required := []string{
		"PORT", "JWT_SECRET",
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASS", "DB_NAME",
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

	// 3) Cliente HTTP externo
	externalAPI := client.NewHttpClient("https://api.openf1.org/v1/")

	// 4) Repos, servicio y controlador
	dRepo := repository.NewDriverRepository(db.DB)
	dService := service.NewDriverService(dRepo, externalAPI)
	dController := api.NewDriverController(dService)

	// 5) Router
	r := gin.Default()
	router.MapUrls(r, dController)

	// 6) Servir
	port := os.Getenv("PORT")
	go func() {
		fmt.Printf("Drivers service en puerto %s...\n", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatalf("Fallo al correr en %s: %v", port, err)
		}
	}()

	// 7) Esperar señal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Deteniendo drivers service...")
}
