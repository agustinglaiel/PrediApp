package router

import (
	"fmt"
	"users/internal/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func MapUrls(engine *gin.Engine, userController *api.UserController) {
    // Use CORS middleware
    engine.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

    // Rutas relacionadas con usuarios
    engine.POST("/users/signup", userController.SignUp)
    engine.POST("/users/login", userController.Login)
    engine.POST("/oauth/signin", userController.OAuthSignIn)

    // Nuevas rutas para administración de usuarios
    engine.GET("/users/:id", userController.GetUserByID)
    engine.GET("/users/username/:username", userController.GetUserByUsername)
    engine.GET("/users", userController.GetUsers)
    engine.PUT("/users/:id", userController.UpdateUserByID)
    engine.PUT("/users/username/:username", userController.UpdateUserByUsername)
    engine.DELETE("/users/:id", userController.DeleteUserByID)
    engine.DELETE("/users/username/:username", userController.DeleteUserByUsername)
	
    fmt.Println("Finishing mappings configurations")
}