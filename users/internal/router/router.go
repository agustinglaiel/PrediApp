package router

import (
	"fmt"
	"users/internal/api"

	"github.com/gin-gonic/gin"
)

func MapUrls(engine *gin.Engine, userController *api.UserController) {
    // Rutas abiertas (sin autenticación)
    engine.POST("/signup", userController.SignUp)
    engine.POST("/login", userController.Login)

    // Rutas para obtener, actualizar o eliminar usuarios
    engine.GET("/:id", userController.GetUserByID)
    engine.GET("/username/:username", userController.GetUserByUsername)
    engine.GET("", userController.GetUsers)
    engine.PUT("/:id", userController.UpdateUserByID)
    engine.PUT("/username/:username", userController.UpdateUserByUsername)
    engine.DELETE("/:id", userController.DeleteUserByID)
    engine.DELETE("/username/:username", userController.DeleteUserByUsername)
	
    fmt.Println("Finishing mappings configurations")
}