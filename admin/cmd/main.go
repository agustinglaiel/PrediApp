package main

import (
	controllerS "admin/internal/api/sessions"
	repoS "admin/internal/repository/sessions"
	"admin/internal/router"
	serviceS "admin/internal/service/sessions"
	"admin/pkg/utils"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inicializar la base de datos
	db, err := utils.InitDB()
	if err != nil {
		fmt.Println("Error al conectar con la Base de Datos")
		panic(err)
	}
	defer utils.DisconnectDB()

	// Iniciar el motor de la base de datos para migrar tablas
	utils.StartDbEngine()

	// Inicializar repositorio y servicio
	sessionRepo := repoS.NewSessionRepository(db)
	sessionService := serviceS.NewSessionService(sessionRepo)

	// Inicializar controlador
	sessionController := controllerS.NewSessionController(sessionService)

	// Configurar el router
	ginRouter := gin.Default()

	// Mapear URLs
	router.MapUrls(ginRouter, sessionController)

	// Iniciar el servidor
	if err := ginRouter.Run(":8060"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
