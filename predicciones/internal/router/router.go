package router

import (
	"fmt"
	"predicciones/internal/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func MapUrls(engine *gin.Engine, prediController *api.PrediController) {
    // Configuración de CORS
    engine.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

    // Rutas relacionadas con predicciones (prodes)
    engine.POST("/prodes/carrera", prediController.CreateProdeCarrera)
    engine.POST("/prodes/session", prediController.CreateProdeSession)

    engine.GET("/prodes/carrera/:id", prediController.GetProdeCarreraByID)
    engine.GET("/prodes/session/:id", prediController.GetProdeSessionByID)
    engine.GET("/prodes/user/:userID", prediController.GetProdesByUserID)

    engine.PUT("/prodes/carrera/:id", prediController.UpdateProdeCarrera)
    engine.PUT("/prodes/session/:id", prediController.UpdateProdeSession)

    engine.DELETE("/prodes/:id", prediController.DeleteProdeByID)

    fmt.Println("Finishing mappings configurations")
}
