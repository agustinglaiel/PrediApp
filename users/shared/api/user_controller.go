package api

import (
	"log"
	"net/http"
	"strconv"
	e "users/pkg/utils"
	"users/shared/dto"
	"users/shared/service"

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

// @Summary Register a new user
// @Description Creates a new user account in the system
// @Tags Users
// @Accept json
// @Produce json
// @Param user body dto.UserSignUpRequestDTO true "User sign-up details"
// @Success 201 {object} dto.UserResponseDTO
// @Failure 400 {object} utils.ApiError
// @Router /users/signup [post]
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

// @Summary Log in a user
// @Description Authenticates a user and returns a JWT token
// @Tags Users
// @Accept json
// @Produce json
// @Param credentials body dto.UserLoginRequestDTO true "User login details"
// @Success 200 {object} dto.UserLoginRequestDTO
// @Failure 401 {object} utils.ApiError
// @Router /users/login [post]
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

// @Summary Get user by ID
// @Description Retrieves a user's details by their ID
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dto.UserResponseDTO
// @Failure 404 {object} utils.ApiError
// @Router /users/{id} [get]
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

// @Summary Get user by username
// @Description Retrieves a user's details by their username
// @Tags Users
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} dto.UserResponseDTO
// @Failure 404 {object} utils.ApiError
// @Router /users/username/{username} [get]
func (ctrl *UserController) GetUserByUsername(c *gin.Context) {
    username := c.Param("username")
    user, apiErr := ctrl.userService.GetUserByUsername(c.Request.Context(), username)
    if apiErr != nil {
        c.JSON(apiErr.Status(), apiErr)
        return
    }
    c.JSON(http.StatusOK, user)
}

// @Summary Get all users
// @Description Retrieves all users in the system
// @Tags Users
// @Produce json
// @Success 200 {array} dto.UserResponseDTO
// @Failure 500 {object} utils.ApiError
// @Router /users [get]
func (ctrl *UserController) GetUsers(c *gin.Context) {
    users, apiErr := ctrl.userService.GetUsers(c.Request.Context())
    if apiErr != nil {
        c.JSON(apiErr.Status(), apiErr)
        return
    }
    c.JSON(http.StatusOK, users)
}

// @Summary Update user by ID
// @Description Updates a user's details by their ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body dto.UserUpdateRequestDTO true "User update details"
// @Success 200 {object} dto.UserResponseDTO
// @Failure 400 {object} utils.ApiError
// @Failure 404 {object} utils.ApiError
// @Router /users/{id} [put]
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