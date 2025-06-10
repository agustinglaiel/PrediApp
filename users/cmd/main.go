package main

import (
	"fmt"
	"log"
	"os"
	"users/internal/api"
	"users/internal/repository"
	"users/internal/router"
	"users/internal/service"
	"users/pkg/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Obtener el puerto de la variable de entorno PORT
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set in the environment")
	}

	// Inicializar la base de datos
	db, err := utils.InitDB()
	if err != nil {
		fmt.Println("Error al conectar con la Base de Datos")
		panic(err)
	}
	defer utils.DisconnectDB()
	utils.StartDbEngine()

	// Inicializar repositorio, servicio y controlador
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := api.NewUserController(userService)

	// Configurar router
	ginRouter := gin.Default()

	// Mapear URLs
	router.MapUrls(ginRouter, userController)

	// Iniciar servidor usando el puerto obtenido de la variable de entorno
	fmt.Printf("Users service listening on port %s...\n", port)
	if err := ginRouter.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server on port %s: %v", port, err)
	}
}