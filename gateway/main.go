package main

import (
	"fmt"
	"gateway/handlers"
	"gateway/middleware"
	"gateway/proxy"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// El punto de entrada principal del gateway, donde se configura el enrutador, los middlewares y el proxy inverso.

func main() {
    // cargar variables de entorno y obtener puerto y secret key 
    err := godotenv.Load()
    if err != nil{
        log.Fatalf("Error loading .env file: %v", err)
    }

    // Obtener el puerto de la variable de entorno PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Valor por defecto
	}

    // Obtener la Secret Key de la variable de entorno
    secretKey := os.Getenv("JWT_SECRET")
    if secretKey == "" {
        log.Fatal("JWT_SECRET is not set in the environment")
    }
    // crear instancia del router gin
    router := gin.Default()

    // configurar middlewares CORS
    router.Use(middleware.CorsMiddleware())

    // configurar ruta para el login
    router.POST("/api/login", handlers.LoginHandler)

    // Configurar proxy inverso para las demas rutas
    router.Any("/api/*proxyPath", proxy.ReverseProxy())

    // Grupo de rutas que requieren autenticación
    protected := router.Group("/api")
    protected.Use(middleware.JwtAuthentication("")) // Aplicar el middleware JWT a las rutas protegidas
    {
        // Aquí van las rutas protegidas
    }

    // Iniciar el servidor HTTP
	fmt.Printf("Gateway listening on port %s...\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server on port %s: %v", port, err)
	}
}