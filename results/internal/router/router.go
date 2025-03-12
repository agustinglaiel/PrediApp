package router

import (
	"results/internal/api"

	"github.com/gin-gonic/gin"
)

func MapUrls(engine *gin.Engine, resultController *api.ResultController) {
	// Use CORS middleware
	// engine.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:3000"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// }))

	// Rutas relacionadas con resultados
	engine.GET("/results/api/:sessionId", resultController.FetchResultsFromExternalAPI)           // Obtener resultados de la API externa para una sesión
	engine.POST("/results", resultController.CreateResult)                                      // Crear un nuevo resultado
	engine.GET("/results/:id", resultController.GetResultByID)                                   // Obtener un resultado por su ID
	engine.PUT("/results/:id", resultController.UpdateResult)                                    // Actualizar un resultado
	engine.DELETE("/results/:id", resultController.DeleteResult)                                 // Eliminar un resultado por su ID
	engine.GET("/results/session/:sessionID", resultController.GetResultsOrderedByPosition)       // Obtener resultados de una sesión ordenados por posición
	engine.GET("/results/session/:sessionID/fastest-lap", resultController.GetFastestLapInSession) // Obtener la vuelta más rápida en una sesión
	engine.GET("/results/driver/:driverID", resultController.GetResultsForDriverAcrossSessions)   // Obtener todos los resultados de un piloto en todas las sesiones
	engine.GET("/results", resultController.GetAllResults)                                       // Obtener todos los resultados
	engine.GET("/results/session/:sessionID/top/:n", resultController.GetTopNDriversInSession)    // Obtener los mejores N pilotos de una sesión
	engine.GET("/results/session/:sessionID/driver/:driverName", resultController.GetResultsForSessionByDriverName) // Obtener resultados de una sesión para un piloto por su nombre
	engine.GET("/results/driver/:driverID/best-position", resultController.GetBestPositionForDriver)  // Obtener la mejor posición de un piloto
	engine.DELETE("/results/session/:sessionID", resultController.DeleteAllResultsForSession)    // Eliminar todos los resultados de una sesión
	engine.GET("/results/driver/:driverID/fastest-laps", resultController.GetTotalFastestLapsForDriver) // Obtener la cantidad total de vueltas más rápidas de un piloto
	engine.GET("/results/driver/:driverID/last", resultController.GetLastResultForDriver)        // Obtener el último resultado de un piloto

	// Ruta para verificar si el servidor está en funcionamiento
	engine.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}
