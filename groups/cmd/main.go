package main

import (
	"fmt"
	"groups/internal/api"
	"groups/internal/repository"
	"groups/internal/router"
	"groups/internal/service"
	"groups/pkg/utils"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Obtener el puerto de la variable de entorno PORT
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set in the environment")
	}

	// Inicializar la base de datos
	db, err := utils.InitDB()
	if err != nil {
		fmt.Println("Error al conectar con la Base de Datos")
		panic(err)
	}
	defer utils.DisconnectDB()

	utils.StartDbEngine()

	groupRepo := repository.NewGroupRepository(db)
	groupService := service.NewGroupService(groupRepo)
	groupController := api.NewGroupController(groupService)

	ginRouter := gin.Default()

	router.MapUrls(ginRouter, groupController)

	// Iniciar servidor usando el puerto obtenido de la variable de entorno
	fmt.Printf("Groups service listening on port %s...\n", port)
	if err := ginRouter.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server on port %s: %v", port, err)
	}
}