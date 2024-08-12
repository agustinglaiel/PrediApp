package api

import (
	"log"
	"net/http"
	"predicciones/internal/dto"
	"predicciones/internal/service"
	e "predicciones/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PrediController struct {
	prediService service.PrediServiceInterface
}

func NewPrediController(prediService service.PrediServiceInterface) *PrediController {
	return &PrediController{
		prediService: prediService,
	}
}

// CreateProdeCarrera handles the creation of a new race prediction
func (ctrl *PrediController) CreateProdeCarrera(c *gin.Context) {
	var request dto.CreateProdeCarreraDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Error binding JSON: %v", err)
		apiErr := e.NewBadRequestApiError("invalid request")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	response, apiErr := ctrl.prediService.CreateProdeCarrera(c.Request.Context(), request)
	if apiErr != nil {
		log.Printf("Error in predi service CreateProdeCarrera: %v", apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// CreateProdeSession handles the creation of a new session prediction
func (ctrl *PrediController) CreateProdeSession(c *gin.Context) {
	var request dto.CreateProdeSessionDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		apiErr := e.NewBadRequestApiError("invalid request")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	response, apiErr := ctrl.prediService.CreateProdeSession(c.Request.Context(), request)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetProdeCarreraByID handles retrieving a race prediction by ID
func (ctrl *PrediController) GetProdeCarreraByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		apiErr := e.NewBadRequestApiError("invalid ID")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	prode, apiErr := ctrl.prediService.GetProdeCarreraByID(c.Request.Context(), uint(id))
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, prode)
}

// GetProdeSessionByID handles retrieving a session prediction by ID
func (ctrl *PrediController) GetProdeSessionByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		apiErr := e.NewBadRequestApiError("invalid ID")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	prode, apiErr := ctrl.prediService.GetProdeSessionByID(c.Request.Context(), uint(id))
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, prode)
}

// GetProdesByUserID handles retrieving all predictions by a user
func (ctrl *PrediController) GetProdesByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		apiErr := e.NewBadRequestApiError("invalid user ID")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	prodesCarrera, prodesSession, apiErr := ctrl.prediService.GetProdesByUserID(c.Request.Context(), uint(userID))
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"carrera": prodesCarrera,
		"session": prodesSession,
	})
}                                 

// UpdateProdeCarrera handles updating an existing race prediction
func (ctrl *PrediController) UpdateProdeCarrera(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        apiErr := e.NewBadRequestApiError("invalid ID")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    var request dto.UpdateProdeCarreraDTO
    if err := c.ShouldBindJSON(&request); err != nil {
        apiErr := e.NewBadRequestApiError("invalid request")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    request.ProdeID = id
    prode, apiErr := ctrl.prediService.UpdateProdeCarrera(c.Request.Context(), request)
    if apiErr != nil {
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    c.JSON(http.StatusOK, prode)
}

// UpdateProdeSession handles updating an existing session prediction
func (ctrl *PrediController) UpdateProdeSession(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        apiErr := e.NewBadRequestApiError("invalid ID")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    var request dto.UpdateProdeSessionDTO
    if err := c.ShouldBindJSON(&request); err != nil {
        apiErr := e.NewBadRequestApiError("invalid request")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    request.ProdeID = id
    prode, apiErr := ctrl.prediService.UpdateProdeSession(c.Request.Context(), request)
    if apiErr != nil {
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    c.JSON(http.StatusOK, prode)
}

// DeleteProdeByID handles deleting a prediction by ID
func (ctrl *PrediController) DeleteProdeByID(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        apiErr := e.NewBadRequestApiError("invalid ID")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    userID, err := strconv.Atoi(c.Query("userID"))
    if err != nil {
        apiErr := e.NewBadRequestApiError("invalid user ID")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

	eventID, err := strconv.Atoi(c.Query("eventID"))
    if err != nil {
        apiErr := e.NewBadRequestApiError("invalid event ID")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    // DeleteProdeByID se encargará de diferenciar si es un prode de carrera o de otra sesión.
    apiErr := ctrl.prediService.DeleteProdeByID(c.Request.Context(), uint(id), uint(eventID) ,uint(userID))
    if apiErr != nil {
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    c.JSON(http.StatusNoContent, nil)
}