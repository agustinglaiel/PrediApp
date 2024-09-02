package api

import (
	dto "admin/internal/dto/prodes"
	prodes "admin/internal/service/prodes"
	e "admin/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProdeController struct {
	prodeService prodes.ProdeServiceInterface
}

func NewProdeController(prodeService prodes.ProdeServiceInterface) *ProdeController {
	return &ProdeController{
		prodeService: prodeService,
	}
}

func (c *ProdeController) CreateProdeCarrera(ctx *gin.Context) {
	var request dto.CreateProdeCarreraDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid JSON data"))
		return
	}

	response, err := c.prodeService.CreateProdeCarrera(ctx.Request.Context(), request)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *ProdeController) CreateProdeSession(ctx *gin.Context) {
	var request dto.CreateProdeSessionDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid JSON data"))
		return
	}

	response, err := c.prodeService.CreateProdeSession(ctx.Request.Context(), request)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *ProdeController) UpdateProdeCarrera(ctx *gin.Context) {
	var request dto.UpdateProdeCarreraDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid JSON data"))
		return
	}

	response, err := c.prodeService.UpdateProdeCarrera(ctx.Request.Context(), request)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *ProdeController) UpdateProdeSession(ctx *gin.Context) {
	var request dto.UpdateProdeSessionDTO
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid JSON data"))
		return
	}

	response, err := c.prodeService.UpdateProdeSession(ctx.Request.Context(), request)
	if err != nil {
		ctx.JSON(err.Status(), err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *ProdeController) DeleteProdeById(ctx *gin.Context) {
	prodeID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid prode ID"))
		return
	}

	if err := c.prodeService.DeleteProdeById(ctx.Request.Context(), prodeID); err != nil {
		ctx.JSON(err.Status(), err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *ProdeController) GetProdesByUserId(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, e.NewBadRequestApiError("Invalid user ID"))
		return
	}

	carreraProdes, sessionProdes, apiErr := c.prodeService.GetProdesByUserId(ctx.Request.Context(), userID)
	if apiErr != nil {
		ctx.JSON(apiErr.Status(), apiErr)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"carrera_prodes": carreraProdes,
		"session_prodes": sessionProdes,
	})
}

