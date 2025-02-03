package main

import (
	"fmt"
	"gateway/handlers"
	"gateway/middleware"
	"gateway/proxy"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// El punto de entrada principal del gateway, donde se configura el enrutador, los middlewares y el proxy inverso.

func main() {
    // Construir la ruta al archivo .env
    currentDir, err := os.Getwd()
    if err != nil {
        fmt.Println("Error al obtener el directorio actual:", err)
        os.Exit(1)
    }
    envPath := filepath.Join(filepath.Dir(currentDir), ".env")

    // cargar variables de entorno y obtener puerto y secret key
    err = godotenv.Load(envPath)
    if err != nil {
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

    // Configurar ruta para el login (Fuera del grupo /api)
    router.POST("/api/login", handlers.LoginHandler)
    router.POST("/api/signup", handlers.SignupHandler)

    // Grupo de rutas que requieren autenticación
    protected := router.Group("/")
    protected.Use(middleware.JwtAuthentication(""))
    {
        protected.Any("/users/*proxyPath", proxy.ReverseProxy())
        protected.Any("/drivers/*proxyPath", proxy.ReverseProxy())
        protected.Any("/prodes/*proxyPath", proxy.ReverseProxy())
        protected.Any("/results/*proxyPath", proxy.ReverseProxy())
        protected.Any("/sessions/*proxyPath", proxy.ReverseProxy())
    }

    // Iniciar el servidor HTTP
    fmt.Printf("Gateway listening on port %s...\n", port)
    if err := router.Run(":" + port); err != nil {
        log.Fatalf("Failed to run server on port %s: %v", port, err)
    }
}