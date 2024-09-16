package api

import (
	dto "drivers/internal/dto/drivers"
	drivers "drivers/internal/service/drivers"
	"drivers/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DriverEventController struct {
	driverEventService drivers.DriverEventServiceInterface
}

func NewDriverEventController(driverEventService drivers.DriverEventServiceInterface) *DriverEventController {
	return &DriverEventController{
		driverEventService: driverEventService,
	}
}

func (c *DriverEventController) AddDriverToEvent(ctx *gin.Context) {
	var request dto.DriverEventDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewBadRequestApiError("Invalid request payload"))
		return
	}

	response, apiErr := c.driverEventService.AddDriverToEvent(ctx.Request.Context(), request)
	if apiErr != nil {
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *DriverEventController) RemoveDriverFromEvent(ctx *gin.Context) {
	driverEventID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewBadRequestApiError("Invalid driver event ID"))
		return
	}

	apiErr := c.driverEventService.RemoveDriverFromEvent(ctx.Request.Context(), uint(driverEventID))
	if apiErr != nil {
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Driver removed from event successfully"})
}

func (c *DriverEventController) ListDriversByEvent(ctx *gin.Context) {
	eventID, err := strconv.ParseUint(ctx.Param("event_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewBadRequestApiError("Invalid event ID"))
		return
	}

	response, apiErr := c.driverEventService.ListDriversByEvent(ctx.Request.Context(), uint(eventID))
	if apiErr != nil {
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *DriverEventController) ListEventsByDriver(ctx *gin.Context) {
	driverID, err := strconv.ParseUint(ctx.Param("driver_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewBadRequestApiError("Invalid driver ID"))
		return
	}

	response, apiErr := c.driverEventService.ListEventsByDriver(ctx.Request.Context(), uint(driverID))
	if apiErr != nil {
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
