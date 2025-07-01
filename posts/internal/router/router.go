package router

import (
	"fmt"

	"prediapp.local/posts/internal/api"

	"github.com/gin-gonic/gin"
)

func MapUrls(engine *gin.Engine, postController *api.PostController) {
	postsGroup := engine.Group("/posts")
	{
		postsGroup.POST("", postController.CreatePost)                    // Crear post o comentario
		postsGroup.GET("", postController.GetPosts)                       // Obtener todos los posts principales (con hilos)
		postsGroup.GET("/:id", postController.GetPostByID)                // Obtener un post por ID (con hilos)
		postsGroup.GET("/user/:user_id", postController.GetPostsByUserID) // Obtener posts de un usuario
		postsGroup.GET("/search", postController.SearchPosts)             // Buscar posts por texto
		postsGroup.DELETE("/:id", postController.DeletePostByID)          // Eliminar un post o comentario (soft delete)
	}

	fmt.Println("Finishing mappings configurations")
}
