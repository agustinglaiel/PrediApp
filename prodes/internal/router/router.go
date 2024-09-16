package router

import (
	prodes "prodes/internal/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func MapUrls(engine *gin.Engine, prodeController *prodes.ProdeController) {
	// Use CORS middleware
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Rutas relacionadas con prodes
	engine.POST("/prodes/carrera", prodeController.CreateProdeCarrera)
	engine.POST("/prodes/session", prodeController.CreateProdeSession)
	engine.PUT("/prodes/carrera/:id", prodeController.UpdateProdeCarrera)
	engine.PUT("/prodes/session/:id", prodeController.UpdateProdeSession)
	engine.DELETE("/prodes/:id", prodeController.DeleteProdeById)
	engine.GET("/prodes/user/:user_id", prodeController.GetProdesByUserId)

	// Debugging purpose
	engine.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}
