package router

import (
	"fmt"

	"prediapp.local/groups/internal/api"

	"github.com/gin-gonic/gin"
)

func MapUrls(engine *gin.Engine, groupController *api.GroupController) {
	groupsGroup := engine.Group("/groups")
	{
		groupsGroup.POST("", groupController.CreateGroup)
		groupsGroup.GET("/:id", groupController.GetGroupByID)
		groupsGroup.GET("/user/:userId", groupController.GetGroupsByUserId)
		groupsGroup.GET("", groupController.GetGroups)
		groupsGroup.DELETE("/:id", groupController.DeleteGroupByID)
		groupsGroup.POST("/join", groupController.JoinGroup)
		groupsGroup.POST("/manage-invitation", groupController.ManageGroupInvitation)
	}

	fmt.Println("Finishing mappings configurations")
}
