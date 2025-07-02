package api

import (
	"log"
	"net/http"
	"strconv"

	"prediapp.local/groups/internal/dto"
	"prediapp.local/groups/internal/service"
	e "prediapp.local/groups/pkg/utils"

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
		apiErr := e.NewBadRequestApiError("Invalid request")
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
		apiErr := e.NewBadRequestApiError("Invalid group ID")
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

func (ctrl *GroupController) GetGroupsByUserId(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		apiErr := e.NewBadRequestApiError("Invalid user ID")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	responses, apiErr := ctrl.groupService.GetGroupsByUserId(c.Request.Context(), userID)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, responses)
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
		apiErr := e.NewBadRequestApiError("Invalid group ID")
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

func (ctrl *GroupController) JoinGroup(c *gin.Context) {
	var request dto.RequestJoinGroupDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		apiErr := e.NewBadRequestApiError("Invalid request")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	apiErr := ctrl.groupService.JoinGroup(c.Request.Context(), request)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Solicitud enviada correctamente"})
}

func (ctrl *GroupController) ManageGroupInvitation(c *gin.Context) {
	var request dto.ManageGroupInvitationDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		apiErr := e.NewBadRequestApiError("Invalid request")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	apiErr := ctrl.groupService.ManageGroupInvitation(c.Request.Context(), request)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Solicitud manejada correctamente"})
}

func (ctrl *GroupController) GetJoinRequests(c *gin.Context) {
	groupId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		apiErr := e.NewBadRequestApiError("Invalid group ID")
		c.JSON(apiErr.Status(), apiErr)
		return
	}
	requests, apiErr := ctrl.groupService.GetJoinRequests(c.Request.Context(), groupId)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}
	c.JSON(http.StatusOK, requests)
}
