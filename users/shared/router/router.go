package router

import (
	"fmt"
	"users/shared/api"

	"github.com/gin-gonic/gin"
)

func MapUrls(routerGroup *gin.Engine, userController *api.UserController) {
    // Rutas abiertas (sin autenticación)
    routerGroup.POST("/signup", userController.SignUp)
    routerGroup.POST("/login", userController.Login)

    // Rutas para obtener, actualizar o eliminar usuarios
    routerGroup.GET("/:id", userController.GetUserByID)
    routerGroup.GET("/username/:username", userController.GetUserByUsername)
    routerGroup.GET("", userController.GetUsers)
    routerGroup.PUT("/:id", userController.UpdateUserByID)
    routerGroup.PUT("/username/:username", userController.UpdateUserByUsername)
    routerGroup.DELETE("/:id", userController.DeleteUserByID)
    routerGroup.DELETE("/username/:username", userController.DeleteUserByUsername)
	
    fmt.Println("Finishing mappings configurations")
}