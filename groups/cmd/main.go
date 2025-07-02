package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"prediapp.local/db"
	"prediapp.local/groups/internal/api"
	"prediapp.local/groups/internal/repository"
	"prediapp.local/groups/internal/router"
	"prediapp.local/groups/internal/service"

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

	// 3) Repos, servicio y controlador
	gRepo := repository.NewGroupRepository(db.DB)
	gService := service.NewGroupService(gRepo)
	gController := api.NewGroupController(gService)

	// 4) Router
	r := gin.Default()
	router.MapUrls(r, gController)

	// 5) Servir
	port := os.Getenv("PORT")
	go func() {
		fmt.Printf("Groups service en puerto %s...\n", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatalf("Fallo al correr en %s: %v", port, err)
		}
	}()

	// 6) Esperar señal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Deteniendo groups service...")
}
