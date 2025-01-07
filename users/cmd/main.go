package main

import (
	"fmt"
	"log"
	"os"
	"users/pkg/utils"
	"users/shared/api"
	"users/shared/repository"
	"users/shared/router"
	"users/shared/service"

	_ "users/docs" // Importa el paquete generado por swag init

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // Swagger UI
	ginSwagger "github.com/swaggo/gin-swagger" // Swagger handler
)

// @title Users Microservice API
// @version 1.0
// @description API documentation for the Users microservice.
// @host localhost:8057
// @BasePath /api
func main() {
	// Obtener el puerto de la variable de entorno PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8057" // Valor por defecto en caso de que no esté configurado
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

	// Inicializar repositorio y servicio
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	// Inicializar controlador
	userController := api.NewUserController(userService)

	// Configurar router
	ginRouter := gin.Default()

	// Mapear URLs
	router.MapUrls(ginRouter, userController)

	// Configurar Swagger UI
	ginRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Iniciar servidor usando el puerto obtenido de la variable de entorno
	fmt.Printf("Users service listening on port %s...\n", port)
	if err := ginRouter.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server on port %s: %v", port, err)
	}
}
