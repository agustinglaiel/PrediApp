package api

import (
	"log"
	"net/http"
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

    log.Printf("Request received: %+v", request)

    response, apiErr := ctrl.userService.SignUp(c.Request.Context(), request)
    if apiErr != nil {
        log.Printf("Error in user service SignUp: %v", apiErr)
        c.JSON(apiErr.Status(), apiErr)
        return
    }

    log.Printf("User created successfully: %+v", response)
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