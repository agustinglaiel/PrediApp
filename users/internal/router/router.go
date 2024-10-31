package router

import (
	"fmt"
	"users/internal/api"
	jwt "users/pkg/jwt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func MapUrls(engine *gin.Engine, userController *api.UserController) {
    // Configurar el middleware CORS
    engine.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

    // Rutas abiertas (sin autenticación)
    engine.POST("users/signup", userController.SignUp)
    engine.POST("users/login", userController.Login)

    // Grupo de rutas protegidas con JWT
    protected := engine.Group("/")
    protected.Use(jwt.JWTAuthMiddleware()) // Aplicar el middleware JWT a las rutas protegidas
    {
        protected.GET("/users/:id", userController.GetUserByID)
        protected.GET("/users/username/:username", userController.GetUserByUsername)
        protected.GET("/users", userController.GetUsers)
        protected.PUT("/users/:id", userController.UpdateUserByID)
        protected.PUT("/users/username/:username", userController.UpdateUserByUsername)
        protected.DELETE("/users/:id", userController.DeleteUserByID)
        protected.DELETE("/users/username/:username", userController.DeleteUserByUsername)
    }
	
    fmt.Println("Finishing mappings configurations")
}
