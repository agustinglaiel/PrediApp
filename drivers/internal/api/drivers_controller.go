package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	dto "prediapp.local/drivers/internal/dto"
	service "prediapp.local/drivers/internal/service"
	"prediapp.local/drivers/pkg/utils"
	e "prediapp.local/drivers/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
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
		// Convertir el error en un mensaje detallado
		validationErrors := make([]interface{}, 0)
		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range ve {
				// Construir el mensaje de error manualmente
				errorMessage := fmt.Sprintf("Field %s: validation failed on tag '%s'", fieldErr.Field(), fieldErr.Tag())
				if fieldErr.Param() != "" {
					errorMessage += fmt.Sprintf(" with parameter '%s'", fieldErr.Param())
				}
				validationErrors = append(validationErrors, errorMessage)
			}
		} else {
			validationErrors = append(validationErrors, err.Error())
		}
		apiErr := utils.NewValidationApiError("Invalid request payload", "bad_request", validationErrors)
		ctx.JSON(apiErr.Status(), apiErr)
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
	driverID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid driver ID"))
		return
	}

	response, apiErr := c.driversService.GetDriverByID(ctx.Request.Context(), driverID)
	if apiErr != nil {
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *DriverController) UpdateDriver(ctx *gin.Context) {
	driverID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid driver ID"))
		return
	}

	var request dto.UpdateDriverDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid request payload"))
		return
	}

	response, apiErr := c.driversService.UpdateDriver(ctx.Request.Context(), driverID, request)
	if apiErr != nil {
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *DriverController) DeleteDriver(ctx *gin.Context) {
	driverID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid driver ID"))
		return
	}

	apiErr := c.driversService.DeleteDriver(ctx.Request.Context(), driverID)
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
	teamName := ctx.Query("team")

	response, apiErr := c.driversService.ListDriversByTeam(ctx.Request.Context(), teamName)
	if apiErr != nil {
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *DriverController) ListDriversByCountry(ctx *gin.Context) {
	countryCode := ctx.Param("countryCode")
	log.Printf("Country code: %s", countryCode)
	response, apiErr := c.driversService.ListDriversByCountry(ctx.Request.Context(), countryCode)
	if apiErr != nil {
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *DriverController) ListDriversByFullName(ctx *gin.Context) {
	fullName := ctx.Param("fullName")
	log.Printf("Full name: %s", fullName)
	response, apiErr := c.driversService.ListDriversByFullName(ctx.Request.Context(), fullName)
	if apiErr != nil {
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *DriverController) ListDriversByAcronym(ctx *gin.Context) {
	acronym := ctx.Param("acronym")
	log.Printf("Acronym: %s", acronym)
	response, apiErr := c.driversService.ListDriversByAcronym(ctx.Request.Context(), acronym)
	if apiErr != nil {
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *DriverController) FetchAllDriversFromExternalAPI(ctx *gin.Context) {
	response, apiErr := c.driversService.FetchAllDriversFromExternalAPI(ctx.Request.Context())
	if apiErr != nil {
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *DriverController) GetDriverByNumber(ctx *gin.Context) {
	driverNumber, err := strconv.Atoi(ctx.Param("driver_number"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid driver number"))
		return
	}

	response, apiErr := c.driversService.GetDriverByNumber(ctx.Request.Context(), driverNumber)
	if apiErr != nil {
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *DriverController) GetDriverByFirstAndLastName(ctx *gin.Context) {
	firstName := ctx.Param("first_name")
	lastName := ctx.Param("last_name")

	response, apiErr := c.driversService.GetDriverByFirstAndLastName(ctx.Request.Context(), firstName, lastName)
	if apiErr != nil {
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
