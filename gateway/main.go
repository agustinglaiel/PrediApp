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

	// Cargar el archivo .env para obtener el valor de ENV
	err = godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Obtener el valor de ENV desde el archivo .env (o desde la variable de entorno si está definida)
	env := os.Getenv("ENV")
	if env == "" {
		log.Println("ENV not set in .env file, defaulting to 'stage'")
		env = "stage"
	}

	// Construir la ruta al archivo .env.stage o .env.prod según el valor de ENV
	envFile := fmt.Sprintf(".env.%s", env)
	envSpecificPath := filepath.Join(filepath.Dir(currentDir), envFile)

	// Cargar las variables de entorno específicas del ambiente
	err = godotenv.Load(envSpecificPath)
	if err != nil {
		log.Printf("Error loading %s file: %v", envFile, err)
		log.Println("Continuing with variables from .env or system environment")
	}

	// Usar el puerto 8080 siempre para el gateway
	port := "8080"

	// Obtener la Secret Key de la variable de entorno
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		log.Fatal("JWT_SECRET is not set in the environment")
	}

	// Crear instancia del router gin
	router := gin.Default()

	// Configurar middlewares CORS
	router.Use(middleware.CorsMiddleware())

	// Rutas públicas
	router.POST("/api/login", handlers.LoginHandler)
	router.POST("/api/signup", handlers.SignupHandler)
	router.POST("/api/refresh", handlers.RefreshTokenHandler)
	router.POST("/api/signout", handlers.SignOutHandler)

	// Rutas proxy
	router.Any("/users", proxy.ReverseProxy())
	router.Any("/users/*proxyPath", proxy.ReverseProxy())
	router.Any("/drivers", proxy.ReverseProxy())
	router.Any("/drivers/*proxyPath", proxy.ReverseProxy())
	router.Any("/prodes", proxy.ReverseProxy())
	router.Any("/prodes/*proxyPath", proxy.ReverseProxy())
	router.Any("/results", proxy.ReverseProxy())
	router.Any("/results/*proxyPath", proxy.ReverseProxy())
	router.Any("/sessions", proxy.ReverseProxy())
	router.Any("/sessions/*proxyPath", proxy.ReverseProxy())
	router.Any("/groups", proxy.ReverseProxy())
	router.Any("/groups/*proxyPath", proxy.ReverseProxy())
	router.Any("/posts", proxy.ReverseProxy())
	router.Any("/posts/*proxyPath", proxy.ReverseProxy())

	// Iniciar el servidor HTTP
	fmt.Printf("Gateway (%s) listening on port %s...\n", env, port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server on port %s: %v", port, err)
	}
}