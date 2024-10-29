package api

import (
	"log"
	"net/http"
	"strconv"
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

// SignUp handles user registration
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

    //log.Printf("User created successfully: %+v", response)
    c.JSON(http.StatusCreated, response)
}

// Login handles user login
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

// OAuthSignIn handles Google OAuth sign-in
// func (ctrl *UserController) OAuthSignIn(c *gin.Context) {
//     var request dto.GoogleOAuthRequestDTO
//     if err := c.ShouldBindJSON(&request); err != nil {
//         apiErr := e.NewBadRequestApiError("invalid request")
//         c.JSON(apiErr.Status(), apiErr)
//         return
//     }
//     response, apiErr := ctrl.userService.OAuthSignIn(c.Request.Context(), request)
//     if apiErr != nil {
//         c.JSON(apiErr.Status(), apiErr)
//         return
//     }
//     c.JSON(http.StatusOK, response)
// }

// GetUserByID handles fetching a user by their ID
func (ctrl *UserController) GetUserByID(c *gin.Context) {
    id := c.Param("id")
    uintID, err := strconv.ParseUint(id, 10, 32)
    if err != nil {
        apiErr := e.NewBadRequestApiError("invalid user ID")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    user, apiErr := ctrl.userService.GetUserById(c.Request.Context(), uint(uintID))
    if apiErr != nil {
        c.JSON(apiErr.Status(), apiErr)
        return
    }
    c.JSON(http.StatusOK, user)
}

// GetUserByUsername handles fetching a user by their username
func (ctrl *UserController) GetUserByUsername(c *gin.Context) {
    username := c.Param("username")
    user, apiErr := ctrl.userService.GetUserByUsername(c.Request.Context(), username)
    if apiErr != nil {
        c.JSON(apiErr.Status(), apiErr)
        return
    }
    c.JSON(http.StatusOK, user)
}

// GetUsers handles fetching all users
func (ctrl *UserController) GetUsers(c *gin.Context) {
    users, apiErr := ctrl.userService.GetUsers(c.Request.Context())
    if apiErr != nil {
        c.JSON(apiErr.Status(), apiErr)
        return
    }
    c.JSON(http.StatusOK, users)
}

// UpdateUserByID handles updating a user by their ID
func (ctrl *UserController) UpdateUserByID(c *gin.Context) {
    id := c.Param("id")
    uintID, err := strconv.ParseUint(id, 10, 32)
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

    user, apiErr := ctrl.userService.UpdateUserById(c.Request.Context(), uint(uintID), request)
    if apiErr != nil {
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    c.JSON(http.StatusOK, user)
}

// UpdateUserByUsername handles updating a user by their username
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

// DeleteUserByID handles deleting a user by their ID
func (ctrl *UserController) DeleteUserByID(c *gin.Context) {
    id := c.Param("id")
    uintID, err := strconv.ParseUint(id, 10, 32)
    if err != nil {
        apiErr := e.NewBadRequestApiError("invalid user ID")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    apiErr := ctrl.userService.DeleteUserById(c.Request.Context(), uint(uintID))
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