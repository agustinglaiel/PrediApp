package router

import (
	"prediapp.local/results/internal/api"

	"github.com/gin-gonic/gin"
)

func MapUrls(engine *gin.Engine, resultController *api.ResultController) {
	// Rutas relacionadas con resultados
	engine.GET("/results/api/:sessionId", resultController.FetchResultsFromExternalAPI) // Obtener resultados de la API externa para una sesión
	engine.GET("/results/session/api/:sessionId", resultController.FetchNonRaceSessionResults)
	engine.POST("/results", resultController.CreateResult)                                         // Crear un nuevo resultado
	engine.POST("/results/admin", resultController.CreateSessionResultsAdmin)                      // Crear los resultados de una session CUALQUIERA siendo ADMIN
	engine.PUT("/results/:id", resultController.UpdateResult)                                      // Actualizar un resultado
	engine.DELETE("/results/:id", resultController.DeleteResult)                                   // Eliminar un resultado por su ID
	engine.GET("/results/session/:sessionID", resultController.GetResultsOrderedByPosition)        // Obtener resultados de una sesión ordenados por posición
	engine.GET("/results/session/:sessionID/fastest-lap", resultController.GetFastestLapInSession) // Obtener la vuelta más rápida en una sesión
	engine.GET("/results", resultController.GetAllResults)                                         // Obtener todos los resultados
	engine.GET("/results/session/:sessionID/top/:n", resultController.GetTopNDriversInSession)     // Obtener los mejores N pilotos de una sesión
	engine.DELETE("/results/session/:sessionID", resultController.DeleteAllResultsForSession)      // Eliminar todos los resultados de una sesión

	// Ruta para verificar si el servidor está en funcionamiento
	engine.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}
