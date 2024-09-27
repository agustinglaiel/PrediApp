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

	// Rutas relacionadas con prodes de carrera
	engine.POST("/prodes/carrera", prodeController.CreateProdeCarrera)
	engine.PUT("/prodes/carrera/:id", prodeController.UpdateProdeCarrera)
	engine.GET("/prodes/carrera/user/:user_id/session/:session_id", prodeController.GetRaceProdeByUserAndSession)
	engine.GET("/prodes/carrera/session/:session_id", prodeController.GetRaceProdesBySession)
	engine.PUT("/prodes/carrera/user/:user_id/session/:session_id", prodeController.UpdateRaceProdeForUserBySessionId)

	// Rutas relacionadas con prodes de sesión
	engine.POST("/prodes/session", prodeController.CreateProdeSession)
	engine.PUT("/prodes/session/:id", prodeController.UpdateProdeSession)
	engine.GET("/prodes/session/user/:user_id/session/:session_id", prodeController.GetSessionProdeByUserAndSession)
	engine.GET("/prodes/session/:session_id", prodeController.GetSessionProdesBySession)

	// Rutas para eliminar prodes
	engine.DELETE("/prodes/:id", prodeController.DeleteProdeById)

	// Rutas relacionadas con usuarios
	engine.GET("/prodes/user/:user_id", prodeController.GetProdesByUserId)

	// Rutas relacionadas con pilotos
	engine.GET("/drivers/:driver_id", prodeController.GetDriverDetails)
	engine.GET("/drivers", prodeController.GetAllDrivers)
	engine.GET("/drivers/top/session/:session_id/n/:n", prodeController.GetTopDriversBySessionId)

	// Debugging purpose
	engine.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}