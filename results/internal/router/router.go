package router

import (
	"results/internal/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func MapUrls(engine *gin.Engine, resultController *api.ResultController) {
	// Use CORS middleware
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Rutas relacionadas con resultados
	engine.POST("/results", resultController.CreateResult)
	engine.GET("/results/:id", resultController.GetResultByID)
	engine.PUT("/results/:id", resultController.UpdateResult)
	engine.DELETE("/results/:id", resultController.DeleteResult)
	engine.GET("/results/api/:sessionId", resultController.FetchResultsFromExternalAPI)
	//engine.GET("/results/session/:sessionID", resultController.GetResultsBySession)
	engine.GET("/results/session/:sessionID/fastest-lap", resultController.GetFastestLapInSession)
	//engine.GET("/results/driver/:driverID", resultController.GetResultsForDriver)
	engine.GET("/results", resultController.GetAllResults)
	//engine.GET("/results/session/:sessionID/top/:n", resultController.GetTopNDriversInSession)
	//engine.GET("/results/session/:sessionID/driver/:driverName", resultController.GetResultsForSessionByDriverName)
	//engine.GET("/results/driver/:driverID/best-position", resultController.GetBestPositionForDriver)
	//engine.DELETE("/results/session/:sessionID", resultController.DeleteAllResultsForSession)
	//engine.GET("/results/driver/:driverID/fastest-laps", resultController.GetTotalFastestLapsForDriver)
	//engine.GET("/results/driver/:driverID/last", resultController.GetLastResultForDriver)

	// Debugging purpose
	engine.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}
