package api

import (
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
	sessionID, err := ParseUintParam(c.Param("sessionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("ID de sesión inválido"))
		return
	}

	// Llamar al servicio para hacer fetch de los resultados desde la API externa
	results, apiErr := rc.resultService.FetchResultsFromExternalAPI(c.Request.Context(), sessionID)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetResultByID obtiene un resultado por su ID
func (rc *ResultController) GetResultByID(c *gin.Context) {
	resultID, err := ParseUintParam(c.Param("resultID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("ID de resultado inválido"))
		return
	}

	result, apiErr := rc.resultService.GetResultByID(c.Request.Context(), resultID)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, result)
}

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
	resultID, err := ParseUintParam(c.Param("resultID"))
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
	resultID, err := ParseUintParam(c.Param("resultID"))
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
	sessionID, err := ParseUintParam(c.Param("sessionID"))
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
	sessionID, err := ParseUintParam(c.Param("sessionID"))
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
func (rc *ResultController) GetResultsForDriverAcrossSessions(c *gin.Context) {
	driverID, err := ParseUintParam(c.Param("driverID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("ID de piloto inválido"))
		return
	}

	results, apiErr := rc.resultService.GetResultsForDriverAcrossSessions(c.Request.Context(), driverID)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, results)
}

// ParseUintParam obtiene un parámetro de la URL y lo convierte a uint
func ParseUintParam(param string) (uint, error) {
	id, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		return 0, e.NewBadRequestApiError("ID inválido")
	}
	return uint(id), nil
}
