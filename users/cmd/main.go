package main

import (
	"users/internal/api"
	"users/internal/repository"
	"users/internal/service"
	"users/pkg/utils"

	"github.com/gin-gonic/gin"
)

func main() {
    // Initialize database
    if err := utils.InitDB(); err != nil {
        panic(err)
    }
    defer utils.DisconnectDB()

    // Initialize repository and service
    userRepo := repository.NewUserRepository()
    userService := service.NewUserService(userRepo)

    // Initialize controller
    userController := api.NewUserController(userService)

    // Set up router
    router := gin.Default()
    router.POST("/signup", userController.SignUp)
    router.POST("/login", userController.Login)
    router.POST("/oauth/signin", userController.OAuthSignIn)

    // Start server
    router.Run(":8080")
}