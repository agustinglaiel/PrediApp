package main

import (
	"fmt"
	"users/internal/api"
	"users/internal/repository"
	"users/internal/router"
	"users/internal/service"
	"users/pkg/utils"

	"github.com/gin-gonic/gin"
)

func main() {
    // Initialize database
    if err := utils.InitDB(); err != nil {
        fmt.Println("Error al conectar con la Base de Datos")
        panic(err)
    }
    defer utils.DisconnectDB()

    // Initialize repository and service
    userRepo := repository.NewUserRepository()
    userService := service.NewUserService(userRepo)
    // Initialize controller
    userController := api.NewUserController(userService)

    //Set up router
    ginRouter := gin.Default()

    //Map URLs
    router.MapUrls(ginRouter, userController)

    // Start server
    ginRouter.Run(":8080")
}