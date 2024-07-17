package router

import (
	"fmt"
	"users/internal/api"

	"github.com/gin-gonic/gin"
)

func MapUrls(engine *gin.Engine, userController *api.UserController) {
    // Rutas relacionadas con usuarios
    engine.POST("/signup", userController.SignUp)
    engine.POST("/login", userController.Login)
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