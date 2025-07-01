package router

import (
	"fmt"

	"prediapp.local/users/internal/api"

	"github.com/gin-gonic/gin"
)

func MapUrls(engine *gin.Engine, userController *api.UserController) {
	usersGroup := engine.Group("/users")
	{
		// Rutas abiertas (sin autenticación)
		usersGroup.POST("/signup", userController.SignUp)
		usersGroup.POST("/login", userController.Login)

		// Rutas para obtener, actualizar o eliminar usuarios
		usersGroup.GET("/:id", userController.GetUserByID)
		usersGroup.GET("/username/:username", userController.GetUserByUsername)
		usersGroup.GET("", userController.GetUsers)
		usersGroup.PUT("/:id", userController.UpdateUserByID)
		usersGroup.PUT("/username/:username", userController.UpdateUserByUsername)
		usersGroup.PUT("/role/:id", userController.UpdateRoleByUserId)
		usersGroup.DELETE("/:id", userController.DeleteUserByID)
		usersGroup.DELETE("/username/:username", userController.DeleteUserByUsername)
	}

	fmt.Println("Finishing mappings configurations")
}
