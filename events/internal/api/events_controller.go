package api

import (
	"net/http"
	"strconv"

	dto "events/internal/dto"
	events "events/internal/service"
	e "events/pkg/utils"

	"github.com/gin-gonic/gin"
)

type EventController struct {
	eventService events.EventService
}

func NewEventController(eventService events.EventService) *EventController {
	return &EventController{
		eventService: eventService,
	}
}

// CreateEvent - Crea un nuevo evento
func (ctrl *EventController) CreateEvent(c *gin.Context) {
	var request dto.CreateEventDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		apiErr := e.NewBadRequestApiError("Invalid request")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	response, apiErr := ctrl.eventService.CreateEvent(c.Request.Context(), request)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetEventByID - Obtiene un evento por su ID
func (c *EventController) GetEventByID(ctrl *gin.Context) {
	idStr := ctrl.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctrl.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	session, err := c.eventService.GetSessionById(ctrl.Request.Context(), id)
	if err != nil {
		ctrl.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve session"})
		return
	}

	ctrl.JSON(http.StatusOK, session)
}

// UpdateEvent - Actualiza un evento existente
func (ctrl *EventController) UpdateEvent(c *gin.Context) {
	id := c.Param("id")
	eventID, err := strconv.Atoi(id)
	if err != nil {
		apiErr := e.NewBadRequestApiError("Invalid event ID")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	var request dto.UpdateEventDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		apiErr := e.NewBadRequestApiError("Invalid request")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	response, apiErr := ctrl.eventService.UpdateEvent(c.Request.Context(), eventID, request)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteEvent - Elimina un evento por su ID
func (ctrl *EventController) DeleteEvent(c *gin.Context) {
	id := c.Param("id")
	eventID, err := strconv.Atoi(id)
	if err != nil {
		apiErr := e.NewBadRequestApiError("Invalid event ID")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	apiErr := ctrl.eventService.DeleteEvent(c.Request.Context(), eventID)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListEvents - Lista todos los eventos
func (ctrl *EventController) ListEvents(c *gin.Context) {
	response, apiErr := ctrl.eventService.ListEvents(c.Request.Context())
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, response)
}
