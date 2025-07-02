package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"prediapp.local/db"
	"prediapp.local/posts/internal/api"
	"prediapp.local/posts/internal/repository"
	"prediapp.local/posts/internal/router"
	"prediapp.local/posts/internal/service"

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

	// 3) Repositorio, servicio y controlador
	pRepo := repository.NewPostRepository(db.DB)
	pService := service.NewPostService(pRepo)
	pController := api.NewPostController(pService)

	// 4) Router
	r := gin.Default()
	router.MapUrls(r, pController)

	// 5) Servir
	port := os.Getenv("PORT")
	go func() {
		fmt.Printf("Posts service en puerto %s...\n", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatalf("Fallo al correr en %s: %v", port, err)
		}
	}()

	// 6) Esperar señal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Deteniendo posts service...")
}
