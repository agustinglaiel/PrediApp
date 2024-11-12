package api

import (
	"groups/internal/dto"
	"groups/internal/service"
	e "groups/pkg/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GroupController struct {
	groupService service.GroupServiceInterface
}

func NewGroupController(groupService service.GroupServiceInterface) *GroupController {
	return &GroupController{
		groupService: groupService,
	}
}

func (ctrl *GroupController) CreateGroup(c *gin.Context) {
	var request dto.CreateGroupRequestDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Error binding JSON: %v", err)
		apiErr := e.NewBadRequestApiError("invalid request")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	response, apiErr := ctrl.groupService.CreateGroup(c.Request.Context(), request)
	if apiErr != nil {
		log.Printf("Error in group service CreateGroup: %v", apiErr)
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (ctrl *GroupController) GetGroupByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		apiErr := e.NewBadRequestApiError("invalid id")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	response, apiErr := ctrl.groupService.GetGroupByID(c.Request.Context(), id)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (ctrl *GroupController) GetGroups(c *gin.Context) {
	response, apiErr := ctrl.groupService.GetGroups(c.Request.Context())
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (ctrl *GroupController) DeleteGroupByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		apiErr := e.NewBadRequestApiError("invalid id")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	apiErr := ctrl.groupService.DeleteGroupByID(c.Request.Context(), id)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (ctrl *GroupController) DeleteGroupByGroupName(c *gin.Context) {
	groupName := c.Param("group_name")

	apiErr := ctrl.groupService.DeleteGroupByGroupName(c.Request.Context(), groupName)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}