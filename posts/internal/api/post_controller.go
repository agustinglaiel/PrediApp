package api

import (
	"net/http"
	"posts/internal/dto"
	"posts/internal/service"
	e "posts/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	postService service.PostService
}

func NewPostController(postService service.PostService) *PostController {
	return &PostController{postService: postService}
}

func (ctrl *PostController) CreatePost(c *gin.Context) {
	var request dto.PostCreateRequestDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		apiErr := e.NewBadRequestApiError("invalid request: " + err.Error())
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	if request.UserID <= 0 {
		apiErr := e.NewBadRequestApiError("user_id is required and must be greater than 0")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	response, apiErr := ctrl.postService.CreatePost(c.Request.Context(), request)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (ctrl *PostController) GetPostByID(c *gin.Context) {
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		apiErr := e.NewBadRequestApiError("invalid post ID: " + err.Error())
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	post, apiErr := ctrl.postService.GetPostByID(c.Request.Context(), intID)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}
	c.JSON(http.StatusOK, post)
}

func (ctrl *PostController) GetPosts(c *gin.Context) {
	// Obtener parámetros de paginación
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "10")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		apiErr := e.NewBadRequestApiError("invalid offset value")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		apiErr := e.NewBadRequestApiError("invalid limit value")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	posts, apiErr := ctrl.postService.GetPosts(c.Request.Context(), offset, limit)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (ctrl *PostController) GetPostsByUserID(c *gin.Context) {
	userID := c.Param("user_id")
	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		apiErr := e.NewBadRequestApiError("invalid user ID: " + err.Error())
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	posts, apiErr := ctrl.postService.GetPostsByUserID(c.Request.Context(), intUserID)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (ctrl *PostController) DeletePostByID(c *gin.Context) {
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		apiErr := e.NewBadRequestApiError("invalid post ID: " + err.Error())
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	var request struct {
		UserID int `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		apiErr := e.NewBadRequestApiError("invalid request: " + err.Error())
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	if request.UserID <= 0 {
		apiErr := e.NewBadRequestApiError("user_id is required and must be greater than 0")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	apiErr := ctrl.postService.DeletePostByID(c.Request.Context(), intID, request.UserID)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (ctrl *PostController) SearchPosts(c *gin.Context) {
	// Obtener parámetro de búsqueda
	query := c.Query("query")
	if query == "" {
		apiErr := e.NewBadRequestApiError("query parameter is required")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// Obtener parámetros de paginación
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "10")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		apiErr := e.NewBadRequestApiError("invalid offset value")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		apiErr := e.NewBadRequestApiError("invalid limit value")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	posts, apiErr := ctrl.postService.SearchPosts(c.Request.Context(), query, offset, limit)
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}
	c.JSON(http.StatusOK, posts)
}