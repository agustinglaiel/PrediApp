package main

import (
	"fmt"
	"groups/internal/api"
	"groups/internal/repository"
	"groups/internal/router"
	"groups/internal/service"
	"groups/pkg/utils"
	"os"

	"github.com/gin-gonic/gin"
)


func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8053"
	}

	db, err := utils.InitDB()
	if err != nil {
		fmt.Println("Error al conectar la Base de Datos")
	}
	defer utils.DisconnectDB()

	utils.StartDbEngine()

	groupRepo := repository.NewGroupRepository(db)
	groupService := service.NewGroupService(groupRepo)
	groupController := api.NewGroupController(groupService)

	ginRouter := gin.Default()

	router.MapUrls(ginRouter, groupController)

	if err := ginRouter.Run(":" + port); err != nil {
		fmt.Printf("Failed to run server on port %s: %v", port, err)
	}
}
