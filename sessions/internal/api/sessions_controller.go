package api

import (
	"net/http"
	"strconv"
	"time"

	dto "prediapp.local/sessions/internal/dto"
	service "prediapp.local/sessions/internal/service"
	e "prediapp.local/sessions/pkg/utils"

	"github.com/gin-gonic/gin"
)

type SessionController struct {
	sessionService service.SessionServiceInterface
}

func NewSessionController(sessionService service.SessionServiceInterface) *SessionController {
	return &SessionController{
		sessionService: sessionService,
	}
}

func (sc *SessionController) CreateSession(c *gin.Context) {
	var request dto.CreateSessionDTO

	// Bind the JSON payload to the DTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Datos inválidos"))
		return
	}

	// Validar que las fechas estén en UTC
	if !request.DateStart.IsZero() && request.DateStart.Location().String() != "UTC" {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("DateStart debe estar en UTC"))
		return
	}
	if !request.DateEnd.IsZero() && request.DateEnd.Location().String() != "UTC" {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("DateEnd debe estar en UTC"))
		return
	}

	// Llamar al servicio para crear la sesión
	response, apiErr := sc.sessionService.CreateSession(c.Request.Context(), request)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// Responder con el DTO de respuesta
	c.JSON(http.StatusCreated, response)
}

func (sc *SessionController) GetSessionById(c *gin.Context) {
	// Obtener el ID de la sesión desde los parámetros de la URL
	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("ID inválido"))
		return
	}

	// Llamar al servicio para obtener la sesión por ID
	response, apiErr := sc.sessionService.GetSessionById(c.Request.Context(), sessionID)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// Responder con el DTO de respuesta
	c.JSON(http.StatusOK, response)
}

func (sc *SessionController) UpdateSessionById(c *gin.Context) {
	// Obtener el ID de la sesión desde los parámetros de la URL
	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("ID inválido"))
		return
	}

	// Vincular la carga útil JSON al DTO de actualización
	var request dto.UpdateSessionDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Datos inválidos"))
		return
	}

	// Validar que las fechas estén en UTC si se proporcionan
	if request.DateStart != nil && request.DateStart.Location().String() != "UTC" {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("DateStart debe estar en UTC"))
		return
	}
	if request.DateEnd != nil && request.DateEnd.Location().String() != "UTC" {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("DateEnd debe estar en UTC"))
		return
	}

	// Llamar al servicio para actualizar la sesión
	response, apiErr := sc.sessionService.UpdateSessionById(c.Request.Context(), sessionID, request)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// Responder con el DTO de respuesta
	c.JSON(http.StatusOK, response)
}

func (sc *SessionController) DeleteSessionById(c *gin.Context) {
	// Obtener el ID de la sesión desde los parámetros de la URL
	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("ID inválido"))
		return
	}

	// Llamar al servicio para eliminar la sesión
	apiErr := sc.sessionService.DeleteSessionById(c.Request.Context(), sessionID)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// Responder con un estado 204 (No Content) si la eliminación fue exitosa
	c.Status(http.StatusNoContent)
}

func (sc *SessionController) ListSessionsByYear(c *gin.Context) {
	// Obtener el año desde los parámetros de la URL
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Año inválido"))
		return
	}

	// Llamar al servicio para obtener las sesiones por año
	response, apiErr := sc.sessionService.ListSessionsByYear(c.Request.Context(), year)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// Responder con el listado de sesiones
	c.JSON(http.StatusOK, response)
}

func (sc *SessionController) GetSessionNameAndTypeById(c *gin.Context) {
	// Obtener el ID de la sesión desde los parámetros de la URL
	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("ID inválido"))
		return
	}

	// Llamar al servicio para obtener el nombre y tipo de la sesión
	response, apiErr := sc.sessionService.GetSessionNameAndTypeById(c.Request.Context(), sessionID)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// Responder con el DTO de respuesta
	c.JSON(http.StatusOK, response)
}

func (sc *SessionController) ListSessionsByCircuitKey(c *gin.Context) {
	// Obtener el CircuitKey desde los parámetros de la URL
	circuitKey, err := strconv.Atoi(c.Param("circuitKey"))
	if err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("CircuitKey inválido"))
		return
	}

	// Llamar al servicio para obtener las sesiones por CircuitKey
	response, apiErr := sc.sessionService.ListSessionsByCircuitKey(c.Request.Context(), circuitKey)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// Responder con el listado de sesiones
	c.JSON(http.StatusOK, response)
}

func (sc *SessionController) ListSessionsByCountryCode(c *gin.Context) {
	// Obtener el CountryCode desde los parámetros de la URL
	countryCode := c.Param("countryCode")

	// Llamar al servicio para obtener las sesiones por CountryCode
	response, apiErr := sc.sessionService.ListSessionsByCountryCode(c.Request.Context(), countryCode)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// Responder con el listado de sesiones
	c.JSON(http.StatusOK, response)
}

func (sc *SessionController) ListUpcomingSessions(c *gin.Context) {
	// Llamar al servicio para obtener las próximas sesiones
	response, apiErr := sc.sessionService.ListUpcomingSessions(c.Request.Context())
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// Responder con el listado de próximas sesiones
	c.JSON(http.StatusOK, response)
}

func (sc *SessionController) ListPastSessions(c *gin.Context) {
	// Obtener el parámetro year de la ruta
	yearStr := c.Param("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 1900 || year > 2100 { // Validación básica del año
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Año inválido"))
		return
	}

	// Llamar al servicio para obtener las sesiones pasadas del año especificado
	response, apiErr := sc.sessionService.ListPastSessions(c.Request.Context(), year)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// Responder con el listado de sesiones pasadas
	c.JSON(http.StatusOK, response)
}

func (sc *SessionController) ListSessionsBetweenDates(c *gin.Context) {
	// Obtener las fechas desde los parámetros de la URL o el query
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	// Parsear las fechas a time.Time
	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Formato de fecha de inicio inválido"))
		return
	}
	if startDate.Location().String() != "UTC" {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Fecha de inicio debe estar en UTC"))
		return
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Formato de fecha de fin inválido"))
		return
	}
	if endDate.Location().String() != "UTC" {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Fecha de fin debe estar en UTC"))
		return
	}

	// Llamar al servicio para obtener las sesiones entre las fechas especificadas
	response, apiErr := sc.sessionService.ListSessionsBetweenDates(c.Request.Context(), startDate, endDate)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// Responder con el listado de sesiones
	c.JSON(http.StatusOK, response)
}

func (sc *SessionController) FindSessionsByNameAndType(c *gin.Context) {
	// Obtener el nombre y tipo de la sesión desde los parámetros de la URL o el query
	sessionName := c.Query("session_name")
	sessionType := c.Query("session_type")

	// Validar que los parámetros no estén vacíos
	if sessionName == "" || sessionType == "" {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Nombre y tipo de sesión son requeridos"))
		return
	}

	// Llamar al servicio para obtener las sesiones por nombre y tipo
	response, apiErr := sc.sessionService.FindSessionsByNameAndType(c.Request.Context(), sessionName, sessionType)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// Responder con el listado de sesiones
	c.JSON(http.StatusOK, response)
}

func (sc *SessionController) GetAllSessions(c *gin.Context) {
	// Llamar al servicio para obtener todas las sesiones
	response, apiErr := sc.sessionService.GetAllSessions(c.Request.Context())
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// Responder con el listado de todas las sesiones
	c.JSON(http.StatusOK, response)
}

func (sc *SessionController) UpdateResultSCAndVSC(c *gin.Context) {
	// Obtener el ID de la sesión desde los parámetros de la URL
	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("ID inválido"))
		return
	}

	// Llamar al servicio para actualizar los resultados de SC y VSC
	apiErr := sc.sessionService.UpdateResultSCAndVSC(c.Request.Context(), sessionID)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// Responder con un estado 200 (OK) si la actualización fue exitosa
	c.Status(http.StatusOK)
}

func (sc *SessionController) UpdateDNF(c *gin.Context) {
	// Obtener el ID de la sesión desde los parámetros de la URL
	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("ID inválido"))
		return
	}

	// Vincular la carga útil JSON al DTO de actualización del DNF
	var request dto.UpdateDNFDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Datos inválidos"))
		return
	}

	// Llamar al servicio para actualizar el DNF
	apiErr := sc.sessionService.UpdateDNFBySessionID(c.Request.Context(), sessionID, request.DNF)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// Responder con un estado 200 (OK) si la actualización fue exitosa
	c.Status(http.StatusOK)
}

func (sc *SessionController) UpdateSessionData(c *gin.Context) {
	// Obtener el ID de la sesión desde los parámetros de la URL
	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("ID de sesión inválido"))
		return
	}

	// Vincular la carga útil JSON al DTO
	var request dto.UpdateSessionDataDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Datos inválidos"))
		return
	}

	// Validar que los parámetros de entrada sean válidos (aunque las fechas se manejan en el servicio)
	if request.Location == "" || request.SessionName == "" || request.SessionType == "" || request.Year == 0 {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Todos los campos son requeridos"))
		return
	}

	// Llamar al servicio para actualizar el session_key, date_start, y date_end
	apiErr := sc.sessionService.UpdateSessionData(c.Request.Context(), sessionID, request.Location, request.SessionName, request.SessionType, request.Year)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// Responder con un estado 200 (OK) si la actualización fue exitosa
	c.Status(http.StatusOK)
}

func (sc *SessionController) GetSessionKeyBySessionID(c *gin.Context) {
	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("ID de sesión inválido"))
		return
	}

	sessionKey, apiErr := sc.sessionService.GetSessionKeyBySessionID(c.Request.Context(), sessionID)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"session_key": sessionKey})
}

func (sc *SessionController) UpdateSessionKeyAdmin(c *gin.Context) {
	// Obtener el ID de la sesión desde los parámetros de la URL
	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("ID inválido"))
		return
	}

	// Bindear el sessionKey del cuerpo de la solicitud
	var request struct {
		SessionKey int `json:"session_key"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Datos inválidos"))
		return
	}

	// Llamar al servicio para actualizar manualmente el session_key
	apiErr := sc.sessionService.UpdateSessionKeyAdmin(c.Request.Context(), sessionID, request.SessionKey)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// Responder con un estado 200 si la actualización fue exitosa
	c.Status(http.StatusOK)
}
