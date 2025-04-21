package api

import (
	"fmt"
	"net/http"
	"results/internal/dto"
	"results/internal/service"
	e "results/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ResultController struct {
	resultService service.ResultService
}

// NewResultController crea un nuevo controlador de resultados
func NewResultController(resultService service.ResultService) *ResultController {
	return &ResultController{
		resultService: resultService,
	}
}

// FetchResultsFromExternalAPI obtiene los resultados desde una API externa y los inserta o actualiza
func (rc *ResultController) FetchResultsFromExternalAPI(c *gin.Context) {
	fmt.Println("Controller: Iniciando FetchResultsFromExternalAPI")

    // Obtener el parámetro sessionId desde la URL
    sessionIDStr := c.Param("sessionId")

    // Convertir sessionId de string a int
    sessionID, err := strconv.Atoi(sessionIDStr)
    if err != nil {
        fmt.Println("Controller: Error al parsear sessionId", err)
        c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("ID de sesión inválido"))
        return
    }
	
	// Llamar al servicio para hacer fetch de los resultados desde la API externa
	results, apiErr := rc.resultService.FetchResultsFromExternalAPI(c.Request.Context(), sessionID)
	if apiErr != nil {
		fmt.Println("Controller: Error en servicio FetchResultsFromExternalAPI", apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// fmt.Println("Controller: Resultados obtenidos:", results)
	c.JSON(http.StatusOK, results)
}

// GetResultByID obtiene un resultado por su ID
// func (rc *ResultController) GetResultByID(c *gin.Context) {
// 	resultID, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("ID de resultado inválido"))
// 		return
// 	}

// 	result, apiErr := rc.resultService.GetResultByID(c.Request.Context(), resultID)
// 	if apiErr != nil {
// 		c.JSON(apiErr.Status(), apiErr)
// 		return
// 	}

// 	c.JSON(http.StatusOK, result)
// }

//ESTO SOLO SIRVE PARA CREAR UN RESULTADO A LA VEZ
// CreateResult crea un nuevo resultado
func (rc *ResultController) CreateResult(c *gin.Context) {
	var request dto.CreateResultDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Datos inválidos"))
		return
	}

	result, apiErr := rc.resultService.CreateResult(c.Request.Context(), request)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusCreated, result)
}

// UpdateResult actualiza un resultado existente
func (rc *ResultController) UpdateResult(c *gin.Context) {
	resultID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("ID de resultado inválido"))
		return
	}

	var request dto.UpdateResultDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Datos inválidos"))
		return
	}

	updatedResult, apiErr := rc.resultService.UpdateResult(c.Request.Context(), resultID, request)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, updatedResult)
}

// DeleteResult elimina un resultado por su ID
func (rc *ResultController) DeleteResult(c *gin.Context) {
	resultID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("ID de resultado inválido"))
		return
	}

	apiErr := rc.resultService.DeleteResult(c.Request.Context(), resultID)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetResultsOrderedByPosition obtiene los resultados de una sesión ordenados por posición
func (rc *ResultController) GetResultsOrderedByPosition(c *gin.Context) {
	sessionID, err := strconv.Atoi(c.Param("sessionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("ID de sesión inválido"))
		return
	}

	results, apiErr := rc.resultService.GetResultsOrderedByPosition(c.Request.Context(), sessionID)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetFastestLapInSession obtiene el piloto con la vuelta más rápida en una sesión
func (rc *ResultController) GetFastestLapInSession(c *gin.Context) {
	sessionID, err := strconv.Atoi(c.Param("sessionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("ID de sesión inválido"))
		return
	}

	fastestLap, apiErr := rc.resultService.GetFastestLapInSession(c.Request.Context(), sessionID)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, fastestLap)
}

// GetAllResults obtiene todos los resultados
func (rc *ResultController) GetAllResults(c *gin.Context) {
	results, apiErr := rc.resultService.GetAllResults(c.Request.Context())
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetResultsForDriverAcrossSessions obtiene todos los resultados de un piloto a través de las sesiones
// func (rc *ResultController) GetResultsForDriverAcrossSessions(c *gin.Context) {
// 	driverID, err := strconv.Atoi(c.Param("driverID"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("ID de piloto inválido"))
// 		return
// 	}

// 	results, apiErr := rc.resultService.GetResultsForDriverAcrossSessions(c.Request.Context(), driverID)
// 	if apiErr != nil {
// 		c.JSON(apiErr.Status(), apiErr)
// 		return
// 	}

// 	c.JSON(http.StatusOK, results)
// }

// GetBestPositionForDriver obtiene la mejor posición de un piloto en cualquier sesión
// func (rc *ResultController) GetBestPositionForDriver(c *gin.Context) {
// 	driverID, err := strconv.Atoi(c.Param("driverID"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("ID de piloto inválido"))
// 		return
// 	}

// 	bestPosition, apiErr := rc.resultService.GetBestPositionForDriver(c.Request.Context(), driverID)
// 	if apiErr != nil {
// 		c.JSON(apiErr.Status(), apiErr)
// 		return
// 	}

// 	c.JSON(http.StatusOK, bestPosition)
// }

// GetTopNDriversInSession obtiene los mejores N pilotos en una sesión
func (rc *ResultController) GetTopNDriversInSession(c *gin.Context) {
	sessionID, err := strconv.Atoi(c.Param("sessionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("ID de sesión inválido"))
		return
	}

	n, err := strconv.Atoi(c.Param("n"))
	if err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Número de pilotos inválido"))
		return
	}

	topDrivers, apiErr := rc.resultService.GetTopNDriversInSession(c.Request.Context(), sessionID, n)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, topDrivers)
}

// DeleteAllResultsForSession elimina todos los resultados de una sesión específica
func (rc *ResultController) DeleteAllResultsForSession(c *gin.Context) {
	sessionID, err := strconv.Atoi(c.Param("sessionID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("ID de sesión inválido"))
		return
	}

	apiErr := rc.resultService.DeleteAllResultsForSession(c.Request.Context(), sessionID)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetResultsForSessionByDriverName obtiene los resultados de una sesión filtrados por el nombre del piloto
// func (rc *ResultController) GetResultsForSessionByDriverName(c *gin.Context) {
// 	sessionID, err := strconv.Atoi(c.Param("sessionID"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("ID de sesión inválido"))
// 		return
// 	}

// 	driverName := c.Param("driverName")
// 	if driverName == "" {
// 		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Nombre del piloto inválido"))
// 		return
// 	}

// 	results, apiErr := rc.resultService.GetResultsForSessionByDriverName(c.Request.Context(), sessionID, driverName)
// 	if apiErr != nil {
// 		c.JSON(apiErr.Status(), apiErr)
// 		return
// 	}

// 	c.JSON(http.StatusOK, results)
// }

// GetTotalFastestLapsForDriver obtiene el total de vueltas rápidas de un piloto
// func (rc *ResultController) GetTotalFastestLapsForDriver(c *gin.Context) {
// 	driverID, err := strconv.Atoi(c.Param("driverID"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("ID de piloto inválido"))
// 		return
// 	}

// 	totalFastestLaps, apiErr := rc.resultService.GetTotalFastestLapsForDriver(c.Request.Context(), driverID)
// 	if apiErr != nil {
// 		c.JSON(apiErr.Status(), apiErr)
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"total_fastest_laps": totalFastestLaps})
// }

// GetLastResultForDriver obtiene el último resultado registrado de un piloto
// func (rc *ResultController) GetLastResultForDriver(c *gin.Context) {
// 	driverID, err := strconv.Atoi(c.Param("driverID"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("ID de piloto inválido"))
// 		return
// 	}

// 	lastResult, apiErr := rc.resultService.GetLastResultForDriver(c.Request.Context(), driverID)
// 	if apiErr != nil {
// 		c.JSON(apiErr.Status(), apiErr)
// 		return
// 	}

// 	c.JSON(http.StatusOK, lastResult)
// }

// CreateBulkResults crea múltiples resultados para una sesión en una sola operación
func (rc *ResultController) CreateSessionResultsAdmin(c *gin.Context) {
    var bulkRequest dto.CreateBulkResultsDTO
    if err := c.ShouldBindJSON(&bulkRequest); err != nil {
        c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Datos inválidos para creación masiva"))
        return
    }

    createdResults, apiErr := rc.resultService.CreateSessionResultsAdmin(c.Request.Context(), bulkRequest)
    if apiErr != nil {
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    // Retornamos la lista de resultados creados
    c.JSON(http.StatusCreated, createdResults)
}
