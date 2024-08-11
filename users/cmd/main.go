package main

import (
	"fmt"
	"log"
	"users/internal/api"
	"users/internal/repository"
	"users/internal/router"
	"users/internal/service"
	"users/pkg/utils"

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

    // Initialize repository and service
    userRepo := repository.NewUserRepository(db)
    userService := service.NewUserService(userRepo)
    // Initialize controller
    userController := api.NewUserController(userService)

    // Set up router
    ginRouter := gin.Default()

    // Map URLs
    router.MapUrls(ginRouter, userController)

    // Start server
    if err := ginRouter.Run(":8080"); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}
