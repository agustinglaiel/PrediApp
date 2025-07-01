package main

import (
	"fmt"
	"log"
	"os"

	"prediapp.local/posts/internal/api"
	"prediapp.local/posts/internal/repository"
	"prediapp.local/posts/internal/router"
	"prediapp.local/posts/internal/service"
	"prediapp.local/posts/pkg/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set in the environment")
	}

	db, err := utils.InitDB()
	if err != nil {
		fmt.Println("Error al conectar con la Base de Datos")
		panic(err)
	}
	defer utils.DisconnectDB()

	utils.StartDbEngine()

	// Inicializar repositorio y servicio
	postRepo := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepo)

	// Inicializar controlador
	postController := api.NewPostController(postService)

	// Configurar router
	ginRouter := gin.Default()
	router.MapUrls(ginRouter, postController)

	// Iniciar servidor
	fmt.Printf("Posts service listening on port %s...\n", port)
	if err := ginRouter.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server on port %s: %v", port, err)
	}
}
