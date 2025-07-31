package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"prediapp.local/db"
	"prediapp.local/sessions/internal/api"
	"prediapp.local/sessions/internal/client"
	"prediapp.local/sessions/internal/repository"
	"prediapp.local/sessions/internal/router"
	"prediapp.local/sessions/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Usar UTC
	time.Local = time.UTC

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

	// 3) Construir cliente y caché
	externalAPI := client.NewHttpClient("https://api.openf1.org/v1/")
	// cache := utils.NewCache(30*time.Minute, 100)

	// 4) Repositorio, servicio y controlador
	sRepo := repository.NewSessionRepository(db.DB)
	sService := service.NewSessionService(sRepo, externalAPI)
	sController := api.NewSessionController(sService)

	// 5) Router
	r := gin.Default()
	router.MapUrls(r, sController)

	// 6) Servir
	port := os.Getenv("PORT")
	go func() {
		fmt.Printf("Sessions service en puerto %s...\n", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatalf("Fallo al correr en %s: %v", port, err)
		}
	}()

	// 7) Esperar señal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Deteniendo sessions service...")

	// // Print cache
	// entries := cache.ListEntries()
	// log.Println("Contenido de la caché de Sessions al cerrar:")
	// for _, entry := range entries {
	// 	log.Printf("Clave: %s, Expiración: %s, Valor: %+v\n", entry.Key, entry.Expiration.Format(time.RFC3339), entry.Value)
	// }
}
