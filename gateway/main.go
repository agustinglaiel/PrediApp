package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"prediapp.local/gateway/handlers"
	"prediapp.local/gateway/middleware"
	"prediapp.local/gateway/proxy"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1) Carga de .env y .env.{stage|prod}
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error al obtener el directorio actual:", err)
		os.Exit(1)
	}
	baseEnvPath := filepath.Join(filepath.Dir(currentDir), ".env")
	if err := godotenv.Load(baseEnvPath); err != nil {
		log.Fatalf("Error loading %s: %v", baseEnvPath, err)
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "stage"
		log.Println("ENV not set, defaulting to 'stage'")
	}
	envPath := filepath.Join(filepath.Dir(currentDir), fmt.Sprintf(".env.%s", env))
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Warning: could not load %s: %v", envPath, err)
	}

	// 2) Leer configuración esencial
	port := os.Getenv("PORT_GATEWAY")
	if port == "" {
		log.Fatal("PORT_GATEWAY is not set")
	}
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	// 3) Inicializar Gin
	router := gin.Default()
	router.Use(middleware.CorsMiddleware())

	// 4) Grupo de rutas bajo /api
	api := router.Group("/api")
	{
		// 4.1) Endpoints públicos de auth
		api.POST("/login", handlers.LoginHandler)
		api.POST("/signup", handlers.SignupHandler)
		api.GET("/auth/me", middleware.JwtAuthentication(""), handlers.MeHandler)

		// 4.2) Rutas proxy para microservicios
		api.Any("/users", proxy.ReverseProxy())
		api.Any("/users/*proxyPath", proxy.ReverseProxy())

		api.Any("/drivers", proxy.ReverseProxy())
		api.Any("/drivers/*proxyPath", proxy.ReverseProxy())

		api.Any("/prodes", proxy.ReverseProxy())
		api.Any("/prodes/*proxyPath", proxy.ReverseProxy())

		api.Any("/results", proxy.ReverseProxy())
		api.Any("/results/*proxyPath", proxy.ReverseProxy())

		api.Any("/sessions", proxy.ReverseProxy())
		api.Any("/sessions/*proxyPath", proxy.ReverseProxy())

		api.Any("/groups", proxy.ReverseProxy())
		api.Any("/groups/*proxyPath", proxy.ReverseProxy())

		api.Any("/posts", proxy.ReverseProxy())
		api.Any("/posts/*proxyPath", proxy.ReverseProxy())
	}

	// 5) Arrancar el servidor
	fmt.Printf("Gateway (%s) listening on port %s...\n", env, port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server on port %s: %v", port, err)
	}
}
