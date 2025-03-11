package router

import (
	prodes "prodes/internal/api"

	"github.com/gin-gonic/gin"
)

func MapUrls(engine *gin.Engine, prodeController *prodes.ProdeController) {
	// Rutas relacionadas con prodes de carrera
	engine.POST("/prodes/carrera", prodeController.CreateProdeCarrera)
	engine.PUT("/prodes/carrera/:prode_id", prodeController.UpdateProdeCarrera)
	// engine.GET("/prodes/carrera/user/:user_id/session/:session_id", prodeController.GetRaceProdeByUserAndSession)
	engine.GET("/prodes/carrera/session/:session_id", prodeController.GetRaceProdesBySession)
	engine.PUT("/prodes/carrera/user/:user_id/session/:session_id", prodeController.UpdateRaceProdeForUserBySessionId)
	engine.POST("/prodes/carrera/:session_id/score", prodeController.UpdateScoresForRace)

	// Rutas relacionadas con prodes de sesión
	engine.POST("/prodes/session", prodeController.CreateProdeSession)
	engine.PUT("/prodes/session/:session_id", prodeController.UpdateProdeSession)
	// engine.GET("/prodes/session/user/:user_id/session/:session_id", prodeController.GetSessionProdeByUserAndSession)
	engine.GET("/prodes/session/:session_id", prodeController.GetSessionProdesBySession)
	engine.POST("/prodes/session/:session_id/score", prodeController.UpdateScoresForSession)

	// Rutas para eliminar prodes
	engine.DELETE("/prodes/:id", prodeController.DeleteProdeById)

	// Rutas relacionadas con usuarios
	engine.GET("/prodes/user/:user_id", prodeController.GetProdesByUserId)
	engine.GET("/prodes/user/:user_id/session/:session_id", prodeController.GetProdeByUserAndSession)

	// Rutas relacionadas con pilotos
	engine.GET("/drivers/:driver_id", prodeController.GetDriverDetails)
	engine.GET("/drivers", prodeController.GetAllDrivers)
	engine.GET("/drivers/top/session/:session_id/n/:n", prodeController.GetTopDriversBySessionId)

	// Debugging purpose
	engine.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}