package api

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"users/internal/dto"
	"users/internal/service"
	e "users/pkg/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
    userService service.UserServiceInterface
}

func NewUserController(userService service.UserServiceInterface) *UserController {
    return &UserController{
        userService: userService,
    }
}

func (ctrl *UserController) SignUp(c *gin.Context) {
    var request dto.UserSignUpRequestDTO
    if err := c.ShouldBindJSON(&request); err != nil {
        log.Printf("Error binding JSON: %v", err)
        apiErr := e.NewBadRequestApiError("invalid request")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    //log.Printf("Request received: %+v", request)

    response, apiErr := ctrl.userService.SignUp(c.Request.Context(), request)
    if apiErr != nil {
        log.Printf("Error in user service SignUp: %v", apiErr)
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    c.JSON(http.StatusCreated, response)
}


func (ctrl *UserController) Login(c *gin.Context) {
    var request dto.UserLoginRequestDTO
    if err := c.ShouldBindJSON(&request); err != nil {
        apiErr := e.NewBadRequestApiError("invalid request")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    response, apiErr := ctrl.userService.Login(c.Request.Context(), request)
    if apiErr != nil {
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    c.JSON(http.StatusOK, response)
}

func (ctrl *UserController) GetUserByID(c *gin.Context) {
    id := c.Param("id")
    intID, err := strconv.Atoi(id) // Cambiado a Atoi para int
    if err != nil {
        apiErr := e.NewBadRequestApiError("invalid user ID")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    user, apiErr := ctrl.userService.GetUserById(c.Request.Context(), intID) // Cambiado a int
    if apiErr != nil {
        c.JSON(apiErr.Status(), apiErr)
        return
    }
    c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) GetUserByUsername(c *gin.Context) {
    username := c.Param("username")
    user, apiErr := ctrl.userService.GetUserByUsername(c.Request.Context(), username)
    if apiErr != nil {
        c.JSON(apiErr.Status(), apiErr)
        return
    }
    c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) GetUsers(c *gin.Context) {
    users, apiErr := ctrl.userService.GetUsers(c.Request.Context())
    if apiErr != nil {
        c.JSON(apiErr.Status(), apiErr)
        return
    }
    c.JSON(http.StatusOK, users)
}

func (ctrl *UserController) UpdateUserByID(c *gin.Context) {
    id := c.Param("id")
    intID, err := strconv.Atoi(id) // Cambiado a Atoi para int
    if err != nil {
        apiErr := e.NewBadRequestApiError("invalid user ID")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    var request dto.UserUpdateRequestDTO
    if err := c.ShouldBindJSON(&request); err != nil {
        apiErr := e.NewBadRequestApiError("invalid request")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    user, apiErr := ctrl.userService.UpdateUserById(c.Request.Context(), intID, request) 
    if apiErr != nil {
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    c.JSON(http.StatusOK, user)
}


func (ctrl *UserController) UpdateUserByUsername(c *gin.Context) {
    username := c.Param("username")
    var request dto.UserUpdateRequestDTO
    if err := c.ShouldBindJSON(&request); err != nil {
        apiErr := e.NewBadRequestApiError("invalid request")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    user, apiErr := ctrl.userService.UpdateUserByUsername(c.Request.Context(), username, request)
    if apiErr != nil {
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    c.JSON(http.StatusOK, user)
}

// @Summary Delete user by ID
// @Description Deletes a user by their ID
// @Tags Users
// @Param id path int true "User ID"
// @Success 204
// @Failure 404 {object} utils.ApiError
// @Router /users/{id} [delete]
func (ctrl *UserController) DeleteUserByID(c *gin.Context) {
    id := c.Param("id")
    intID, err := strconv.Atoi(id) // Cambiado a Atoi para int
    if err != nil {
        apiErr := e.NewBadRequestApiError("invalid user ID")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    apiErr := ctrl.userService.DeleteUserById(c.Request.Context(), intID) // Cambiado a int
    if apiErr != nil {
        c.JSON(apiErr.Status(), apiErr)
        return
    }
    c.JSON(http.StatusNoContent, nil)
}

// DeleteUserByUsername handles deleting a user by their username
func (ctrl *UserController) DeleteUserByUsername(c *gin.Context) {
    username := c.Param("username")
    apiErr := ctrl.userService.DeleteUserByUsername(c.Request.Context(), username)
    if apiErr != nil {
        c.JSON(apiErr.Status(), apiErr)
        return
    }
    c.JSON(http.StatusNoContent, nil)
}

func (ctrl *UserController) UpdateRoleByUserId(c *gin.Context) {
    id := c.Param("id")
    intID, err := strconv.Atoi(id)
    if err != nil {
        apiErr := e.NewBadRequestApiError("invalid user ID")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    var request dto.UserUpdateRoleRequestDTO
    if err := c.ShouldBindJSON(&request); err != nil {
        apiErr := e.NewBadRequestApiError("invalid request")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    user, apiErr := ctrl.userService.UpdateRoleByUserId(c.Request.Context(), intID, request)
    if apiErr != nil {
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) StoreRefreshToken(c *gin.Context) {
    var request map[string]interface{}
    if err := c.ShouldBindJSON(&request); err != nil {
        apiErr := e.NewBadRequestApiError("invalid request")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    userIDFloat, ok := request["user_id"].(float64)
    if !ok {
        apiErr := e.NewBadRequestApiError("invalid user_id")
        c.JSON(apiErr.Status(), apiErr)
        return
    }
    userID := int(userIDFloat)
    token, ok := request["token"].(string)
    if !ok {
        apiErr := e.NewBadRequestApiError("invalid token")
        c.JSON(apiErr.Status(), apiErr)
        return
    }
    expiresAtStr, ok := request["expires_at"].(string)
    if !ok {
        apiErr := e.NewBadRequestApiError("invalid expires_at")
        c.JSON(apiErr.Status(), apiErr)
        return
    }
    expiresAt, err := time.Parse(time.RFC3339, expiresAtStr)
    if err != nil {
        apiErr := e.NewBadRequestApiError("invalid expires_at format")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    apiErr := ctrl.userService.StoreRefreshToken(c.Request.Context(), userID, token, expiresAt)
    if apiErr != nil {
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    c.JSON(http.StatusCreated, nil)
}

func (ctrl *UserController) Refresh(c *gin.Context) {
    var request map[string]interface{}
    if err := c.ShouldBindJSON(&request); err != nil {
        apiErr := e.NewBadRequestApiError("invalid request")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    refreshToken, ok := request["refresh_token"].(string)
    if !ok {
        apiErr := e.NewBadRequestApiError("invalid refresh_token")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    user, apiErr := ctrl.userService.ValidateRefreshToken(c.Request.Context(), refreshToken)
    if apiErr != nil {
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "user_id": user.ID,
        "role":    user.Role,
    })
}

func (ctrl *UserController) SignOut(c *gin.Context) {
    var request map[string]interface{}
    if err := c.ShouldBindJSON(&request); err != nil {
        apiErr := e.NewBadRequestApiError("invalid request")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    refreshToken, ok := request["refresh_token"].(string)
    if !ok {
        apiErr := e.NewBadRequestApiError("invalid refresh_token")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    apiErr := ctrl.userService.RevokeRefreshToken(c.Request.Context(), refreshToken)
    if apiErr != nil {
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    c.JSON(http.StatusNoContent, nil)
}