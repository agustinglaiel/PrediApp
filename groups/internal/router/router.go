package router

import (
	"fmt"
	"groups/internal/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)


func MapUrls(engine *gin.Engine, groupController *api.GroupController) {
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	engine.POST("groups", groupController.CreateGroup)
	engine.GET("groups/:id", groupController.GetGroupByID)
	engine.GET("groups", groupController.GetGroups)
	engine.DELETE("groups/:id", groupController.DeleteGroupByID)
	engine.DELETE("groups/groupName/:group_name", groupController.DeleteGroupByGroupName)

	fmt.Println("Finishing mappings configurations")
}
