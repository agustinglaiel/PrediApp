package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"results/internal/api"
	"results/internal/client"
	"results/internal/repository"
	"results/internal/router"
	"results/internal/service"
	"results/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Función para buscar el archivo .env en el directorio raíz
func loadEnv() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error al obtener el directorio actual: %v", err)
	}

	for {
		envPath := filepath.Join(currentDir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			// Si encontramos el archivo .env, lo cargamos
			err = godotenv.Load(envPath)
			if err != nil {
				log.Fatalf("Error al cargar el archivo .env: %v", err)
			}
			fmt.Printf("Archivo .env cargado desde: %s\n", envPath)
			return
		}

		// Si llegamos al directorio raíz del sistema y no encontramos el archivo, salimos
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			log.Fatalf("Archivo .env no encontrado en la jerarquía de directorios")
		}

		// Continuamos buscando en el directorio padre
		currentDir = parentDir
	}
}

func main() {
	loadEnv()

	// Obtener el puerto de la variable de entorno PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8055" // Valor por defecto en caso de que no esté configurado
	}

	// Inicializar la base de datos
	db, err := utils.InitDB()
	if err != nil {
		fmt.Println("Error al conectar con la Base de Datos")
		panic(err)
	}
	defer utils.DisconnectDB()

	// Iniciar el motor de la base de datos y migrar tablas
	utils.StartDbEngine()

	// Crear el cliente HTTP para interactuar con la API externa
	externalAPIClient := client.NewHttpClient("http://localhost:8080")

	// Inicializar repositorio y servicio
	resultRepo := repository.NewResultRepository(db)
	resultService := service.NewResultService(resultRepo, externalAPIClient)
	resultController := api.NewResultController(resultService)

	// Configurar router
	ginRouter := gin.Default()

	// Mapear URLs
	router.MapUrls(ginRouter, resultController)

	// Iniciar servidor usando el puerto obtenido de la variable de entorno
	if err := ginRouter.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server on port %s: %v", port, err)
	}
}
