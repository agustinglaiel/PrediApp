package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"prediapp.local/db"
	"prediapp.local/users/internal/api"
	"prediapp.local/users/internal/repository"
	"prediapp.local/users/internal/router"
	"prediapp.local/users/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1) Validar variables de entorno críticas
	requiredEnvVars := []string{
		"PORT",
		"JWT_SECRET",
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASS", "DB_NAME",
	}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("%s is not set in the environment", envVar)
		}
	}

	// 2) Inicializar conexión (sin migrar)
	if err := db.Init(); err != nil {
		log.Fatalf("db.Init failed: %v", err)
	}

	// 3) Inicializar repositorio, servicio y controlador
	userRepo := repository.NewUserRepository(db.DB)
	userService := service.NewUserService(userRepo)
	userController := api.NewUserController(userService)

	// 4) Configurar router
	ginRouter := gin.Default()
	router.MapUrls(ginRouter, userController)

	// 5) Iniciar servidor
	port := os.Getenv("PORT")
	go func() {
		fmt.Printf("Users service listening on port %s...\n", port)
		if err := ginRouter.Run(":" + port); err != nil {
			log.Fatalf("Failed to run server on port %s: %v", port, err)
		}
	}()

	// 6) Esperar señal de interrupción
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	log.Println("Stopping users service...")
}
