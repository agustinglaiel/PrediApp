package api

import (
	dto "admin/internal/dto/drivers"
	service "admin/internal/service/drivers"
	e "admin/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DriverController struct {
	driversService service.DriverService
}

func NewDriverController(driversService service.DriverService) *DriverController {
	return &DriverController{
		driversService: driversService,
	}
}

func (c *DriverController) CreateDriver(ctx *gin.Context) {
	var request dto.CreateDriverDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid request payload"))
		return
	}

	response, apiErr := c.driversService.CreateDriver(ctx.Request.Context(), request)
	if apiErr != nil {
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *DriverController) GetDriverByID(ctx *gin.Context) {
	driverID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid driver ID"))
		return
	}

	response, apiErr := c.driversService.GetDriverByID(ctx.Request.Context(), uint(driverID))
	if apiErr != nil {
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *DriverController) UpdateDriver(ctx *gin.Context) {
	driverID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid driver ID"))
		return
	}

	var request dto.UpdateDriverDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid request payload"))
		return
	}

	response, apiErr := c.driversService.UpdateDriver(ctx.Request.Context(), uint(driverID), request)
	if apiErr != nil {
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *DriverController) DeleteDriver(ctx *gin.Context) {
	driverID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid driver ID"))
		return
	}

	apiErr := c.driversService.DeleteDriver(ctx.Request.Context(), uint(driverID))
	if apiErr != nil {
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Driver deleted successfully"})
}

func (c *DriverController) ListDrivers(ctx *gin.Context) {
	response, apiErr := c.driversService.ListDrivers(ctx.Request.Context())
	if apiErr != nil {
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *DriverController) ListDriversByTeam(ctx *gin.Context) {
	teamName := ctx.Param("team")

	response, apiErr := c.driversService.ListDriversByTeam(ctx.Request.Context(), teamName)
	if apiErr != nil {
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *DriverController) ListDriversByCountry(ctx *gin.Context) {
	countryCode := ctx.Param("country")

	response, apiErr := c.driversService.ListDriversByCountry(ctx.Request.Context(), countryCode)
	if apiErr != nil {
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *DriverController) ListDriversByFullName(ctx *gin.Context) {
	fullName := ctx.Param("fullname")

	response, apiErr := c.driversService.ListDriversByFullName(ctx.Request.Context(), fullName)
	if apiErr != nil {
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *DriverController) ListDriversByAcronym(ctx *gin.Context) {
	acronym := ctx.Param("acronym")

	response, apiErr := c.driversService.ListDriversByAcronym(ctx.Request.Context(), acronym)
	if apiErr != nil {
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
