package api

import (
	dto "admin/internal/dto/users"
	service "admin/internal/service/users"
	e "admin/pkg/utils"
	"log"
	"net/http"
	"strconv"

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
func (ctrl *UserController) OAuthSignIn(c *gin.Context) {
    var request dto.GoogleOAuthRequestDTO
    if err := c.ShouldBindJSON(&request); err != nil {
        apiErr := e.NewBadRequestApiError("invalid request")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    response, apiErr := ctrl.userService.OAuthSignIn(c.Request.Context(), request)
    if apiErr != nil {
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    c.JSON(http.StatusOK, response)
}

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

// GetUsers handles fetching all users with pagination
func (ctrl *UserController) GetUsers(c *gin.Context) {
    // Parse limit and offset from query parameters
    limit, err := strconv.Atoi(c.Query("limit"))
    if err != nil || limit <= 0 {
        limit = 10 // Set a default value if not provided or invalid
    }

    offset, err := strconv.Atoi(c.Query("offset"))
    if err != nil || offset < 0 {
        offset = 0 // Set to 0 if not provided or invalid
    }

    // Call the service layer with the parsed limit and offset
    users, apiErr := ctrl.userService.GetUsers(c.Request.Context(), limit, offset)
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

// UpdateUserRoleByID handles updating a user's role by their ID
func (ctrl *UserController) UpdateUserRoleByID(c *gin.Context) {
    id := c.Param("id")
    uintID, err := strconv.ParseUint(id, 10, 32)
    if err != nil {
        apiErr := e.NewBadRequestApiError("invalid user ID")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    var request struct {
        Role string `json:"role"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        apiErr := e.NewBadRequestApiError("invalid request")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    if request.Role == "" {
        apiErr := e.NewBadRequestApiError("role is required")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    apiErr := ctrl.userService.UpdateUserRole(c.Request.Context(), uint(uintID), request.Role)
    if apiErr != nil {
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User role updated successfully"})
}


// DeactivateUserByID handles deactivating a user by their ID
func (ctrl *UserController) DeactivateUserByID(c *gin.Context) {
    id := c.Param("id")
    uintID, err := strconv.ParseUint(id, 10, 32)
    if err != nil {
        apiErr := e.NewBadRequestApiError("invalid user ID")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    apiErr := ctrl.userService.DeactivateUserById(c.Request.Context(), uint(uintID))
    if apiErr != nil {
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User deactivated successfully"})
}

// ReactivateUserByID handles reactivating a user by their ID
func (ctrl *UserController) ReactivateUserByID(c *gin.Context) {
    id := c.Param("id")
    uintID, err := strconv.ParseUint(id, 10, 32)
    if err != nil {
        apiErr := e.NewBadRequestApiError("invalid user ID")
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    apiErr := ctrl.userService.ReactivateUserById(c.Request.Context(), uint(uintID))
    if apiErr != nil {
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User reactivated successfully"})
}