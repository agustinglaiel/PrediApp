package router

import (
	"sessions/internal/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func MapUrls(engine *gin.Engine, sessionController *api.SessionController) {
	// Use CORS middleware
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Rutas relacionadas con sesiones
	engine.POST("/sessions", sessionController.CreateSession)
	engine.GET("/sessions/:id", sessionController.GetSessionById)
	engine.PUT("/sessions/:id", sessionController.UpdateSessionById)
	engine.DELETE("/sessions/:id", sessionController.DeleteSessionById)
	engine.GET("/sessions/year/:year", sessionController.ListSessionsByYear)
	engine.GET("/sessions/circuit/:circuitKey", sessionController.ListSessionsByCircuitKey)
	engine.GET("/sessions/country/:countryCode", sessionController.ListSessionsByCountryCode)
	engine.GET("/sessions/upcoming", sessionController.ListUpcomingSessions)
	// engine.GET("/sessions/date-range", sessionController.ListSessionsBetweenDates) VER DESPUES
	engine.GET("/sessions/name-type", sessionController.FindSessionsByNameAndType)
	engine.GET("/sessions/:id/name-type", sessionController.GetSessionNameAndTypeById)
	engine.GET("/sessions", sessionController.GetAllSessions)
	engine.PUT("/sessions/:id/scvsc", sessionController.UpdateResultSCAndVSC)
	//engine.PUT("/sessions/:id/dnf", sessionController.CalculateDNF)
	engine.PUT("/sessions/:id/dnf", sessionController.UpdateDNF)


	// Debugging purpose
	engine.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}
