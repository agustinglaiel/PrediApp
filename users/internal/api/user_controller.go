package api

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"prediapp.local/users/internal/dto"
	"prediapp.local/users/internal/service"
	e "prediapp.local/users/pkg/utils"

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

func (ctrl *UserController) GetUserScoreByUserId(c *gin.Context) {
	id := c.Param("id")
	intID, err := strconv.Atoi(id) // Cambiado a Atoi para int
	if err != nil {
		apiErr := e.NewBadRequestApiError("invalid user ID")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	score, apiErr := ctrl.userService.GetUserScoreByUserId(c.Request.Context(), intID) // Cambiado a int
	if apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, score)
}

func (ctrl *UserController) UploadProfilePicture(c *gin.Context) {
	// 1) Parsear y validar el ID
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		apiErr := e.NewBadRequestApiError("ID de usuario no válido")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// 2) Obtener el archivo
	file, err := c.FormFile("profile_picture")
	if err != nil {
		apiErr := e.NewBadRequestApiError("error al subir el archivo")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// 3) Validar tamaño (200 KB)
	const maxSize = 200 * 1024
	if file.Size > maxSize {
		apiErr := e.NewBadRequestApiError("El archivo no puede excederse de 200 KB")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// 4) Validar tipo MIME
	mimeType := file.Header.Get("Content-Type")
	if mimeType != "image/jpeg" && mimeType != "image/jpg" {
		apiErr := e.NewBadRequestApiError("Solo se permiten JPG/JPEG")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// 5) Leer datos
	f, err := file.Open()
	if err != nil {
		apiErr := e.NewInternalServerApiError("No se pudo abrir el archivo", err)
		c.JSON(apiErr.Status(), apiErr)
		return
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		apiErr := e.NewInternalServerApiError("Error al leer el archivo", err)
		c.JSON(apiErr.Status(), apiErr)
		return
	}
	// doble check por seguridad
	if len(data) > maxSize {
		apiErr := e.NewBadRequestApiError("El archivo no puede excederse de 200 KB")
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// 6) Llamar al servicio con los bytes y el MIME
	if apiErr := ctrl.userService.UploadProfilePicture(
		c.Request.Context(),
		userID,
		data,
		mimeType,
	); apiErr != nil {
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// 7) Respuesta
	c.JSON(http.StatusOK, gin.H{"message": "Profile picture uploaded successfully"})
}
