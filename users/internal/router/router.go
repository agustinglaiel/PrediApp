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
    engine.POST("/signup", userController.SignUp)
    engine.POST("/login", userController.Login)

    // Grupo de rutas protegidas con JWT
    protected := engine.Group("/")
    protected.Use(jwt.JWTAuthMiddleware()) // Aplicar el middleware JWT a las rutas protegidas
    {
        protected.GET("/:id", userController.GetUserByID)
        protected.GET("/username/:username", userController.GetUserByUsername)
        protected.GET("", userController.GetUsers)
        protected.PUT("/:id", userController.UpdateUserByID)
        protected.PUT("/username/:username", userController.UpdateUserByUsername)
        protected.DELETE("/:id", userController.DeleteUserByID)
        protected.DELETE("/username/:username", userController.DeleteUserByUsername)
    }
	
    fmt.Println("Finishing mappings configurations")
}
