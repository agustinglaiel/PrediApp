package router

import (
	"fmt"
	"groups/internal/api"

	"github.com/gin-gonic/gin"
)


func MapUrls(engine *gin.Engine, groupController *api.GroupController) {
	groupsGroup := engine.Group("/groups")
	{
		groupsGroup.POST("", groupController.CreateGroup)        // Crear grupo
		groupsGroup.GET("/:id", groupController.GetGroupByID)    // Obtener grupo por ID
		groupsGroup.GET("", groupController.GetGroups)           // Listar grupos
		groupsGroup.DELETE("/:id", groupController.DeleteGroupByID) // Eliminar grupo
		groupsGroup.POST("/join", groupController.JoinGroup)  // Unirse a un grupo con código
		groupsGroup.POST("/manage-invitation", groupController.ManageGroupInvitation) // Aceptar/rechazar usuario
	}
	
	fmt.Println("Finishing mappings configurations")
}
