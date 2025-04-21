package handlers

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Estructura para el JWT Payload
type Claims struct {
	UserID int   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Estructura para la respuesta de Login y Signup
type UserResponseDTO struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	Token        string `json:"token"`        // JWT token
	RefreshToken string `json:"refresh_token"` // Refresh token
	CreatedAt    string `json:"created_at"`
}

func GenerateTokens(userID int, role string, secretKey string) (string, string, error) {
	// Token de acceso (15 minutos)
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", fmt.Errorf("error generating access token: %w", err)
	}

	// Refresh token (7 días)
	refreshTokenBytes := make([]byte, 32)
	if _, err := rand.Read(refreshTokenBytes); err != nil {
		return "", "", fmt.Errorf("error generating refresh token: %w", err)
	}
	refreshToken := base64.URLEncoding.EncodeToString(refreshTokenBytes)

	return accessToken, refreshToken, nil
}

func LoginHandler(c *gin.Context) {
	// Define la URL del microservicio de usuarios
	usersServiceURL := os.Getenv("USERS_SERVICE_URL")
	secretKey := os.Getenv("JWT_SECRET")

	// 1. Recibir la Solicitud de Login
	var requestBody map[string]interface{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Printf("Error parsing request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// 2. Comunicarse con el Microservicio de Usuarios
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Error marshaling request body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	req, err := http.NewRequest("POST", usersServiceURL+"/users/login", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Error creating request to users service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request to users service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer resp.Body.Close()

	// 3. Validar la Respuesta del Microservicio
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Error response from users service: %s", string(body))
		c.JSON(resp.StatusCode, gin.H{"error": "invalid credentials"})
		return
	}

	// Leer la respuesta del microservicio
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	var userResponse UserResponseDTO
	if err := json.Unmarshal(responseBody, &userResponse); err != nil {
		log.Printf("Error unmarshaling response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// 4. Generar ambos tokens
	accessToken, refreshToken, err := GenerateTokens(userResponse.ID, userResponse.Role, secretKey)
	if err != nil {
		log.Printf("Error generating tokens: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// 5. Almacenar el refresh token en el microservicio de users
	refreshReqBody := map[string]interface{}{
		"user_id":    userResponse.ID,
		"token":      refreshToken,
		"expires_at": time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339),
	}
	refreshJson, err := json.Marshal(refreshReqBody)
	if err != nil {
		log.Printf("Error marshaling refresh token body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	refreshReq, err := http.NewRequest("POST", usersServiceURL+"/users/refresh-token", bytes.NewBuffer(refreshJson))
	if err != nil {
		log.Printf("Error creating refresh token request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	refreshReq.Header.Set("Content-Type", "application/json")
	refreshResp, err := client.Do(refreshReq)
	if err != nil || refreshResp.StatusCode != http.StatusCreated {
		log.Printf("Error storing refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer refreshResp.Body.Close()

	// 6. Asignar tokens a la respuesta
	userResponse.Token = accessToken
	userResponse.RefreshToken = refreshToken

	// 7. Devolver la Respuesta con los tokens
	c.JSON(http.StatusOK, userResponse)
}

func SignupHandler(c *gin.Context) {
	// Define la URL del microservicio de usuarios
	usersServiceURL := os.Getenv("USERS_SERVICE_URL")
	secretKey := os.Getenv("JWT_SECRET")

	// 1. Recibir la Solicitud de Signup
	var requestBody map[string]interface{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Printf("Error parsing request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// 2. Comunicarse con el Microservicio de Usuarios
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Error marshaling request body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	req, err := http.NewRequest("POST", usersServiceURL+"/users/signup", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Error creating request to users service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request to users service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer resp.Body.Close()

	// 3. Validar la Respuesta del Microservicio
	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Error response from users service: %s", string(body))
		c.JSON(resp.StatusCode, gin.H{"error": "signup failed"})
		return
	}

	// Leer la respuesta del microservicio
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	var userResponse UserResponseDTO
	if err := json.Unmarshal(responseBody, &userResponse); err != nil {
		log.Printf("Error unmarshaling response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// 4. Generar ambos tokens
	accessToken, refreshToken, err := GenerateTokens(userResponse.ID, userResponse.Role, secretKey)
	if err != nil {
		log.Printf("Error generating tokens: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// 5. Almacenar el refresh token en el microservicio de users
	refreshReqBody := map[string]interface{}{
		"user_id":    userResponse.ID,
		"token":      refreshToken,
		"expires_at": time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339),
	}
	refreshJson, err := json.Marshal(refreshReqBody)
	if err != nil {
		log.Printf("Error marshaling refresh token body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	refreshReq, err := http.NewRequest("POST", usersServiceURL+"/users/refresh-token", bytes.NewBuffer(refreshJson))
	if err != nil {
		log.Printf("Error creating refresh token request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	refreshReq.Header.Set("Content-Type", "application/json")
	refreshResp, err := client.Do(refreshReq)
	if err != nil || refreshResp.StatusCode != http.StatusCreated {
		log.Printf("Error storing refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer refreshResp.Body.Close()

	// 6. Asignar tokens a la respuesta
	userResponse.Token = accessToken
	userResponse.RefreshToken = refreshToken

	// 7. Devolver la Respuesta con los tokens
	c.JSON(http.StatusCreated, userResponse)
}

func RefreshTokenHandler(c *gin.Context) {
	secretKey := os.Getenv("JWT_SECRET")
	usersServiceURL := os.Getenv("USERS_SERVICE_URL")

	// 1. Recibir el refresh token
	var request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Error parsing refresh token request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// 2. Validar el refresh token con el microservicio de users
	refreshReqBody := map[string]interface{}{
		"refresh_token": request.RefreshToken,
	}
	refreshJson, err := json.Marshal(refreshReqBody)
	if err != nil {
		log.Printf("Error marshaling refresh token body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	req, err := http.NewRequest("POST", usersServiceURL+"/users/refresh", bytes.NewBuffer(refreshJson))
	if err != nil {
		log.Printf("Error creating refresh request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending refresh request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Error response from users service: %s", string(body))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	// 3. Leer la respuesta del microservicio
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading refresh response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	var refreshResponse struct {
		UserID int    `json:"user_id"`
		Role   string `json:"role"`
	}
	if err := json.Unmarshal(responseBody, &refreshResponse); err != nil {
		log.Printf("Error unmarshaling refresh response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// 4. Generar un nuevo token de acceso
	accessToken, _, err := GenerateTokens(refreshResponse.UserID, refreshResponse.Role, secretKey)
	if err != nil {
		log.Printf("Error generating new access token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// 5. Devolver el nuevo token de acceso
	c.JSON(http.StatusOK, gin.H{
		"token": accessToken,
	})
}

func SignOutHandler(c *gin.Context) {
	usersServiceURL := os.Getenv("USERS_SERVICE_URL")

	// 1. Recibir el refresh token
	var request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Error parsing signout request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// 2. Enviar solicitud de revocación al microservicio de users
	signoutReqBody := map[string]interface{}{
		"refresh_token": request.RefreshToken,
	}
	signoutJson, err := json.Marshal(signoutReqBody)
	if err != nil {
		log.Printf("Error marshaling signout body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	req, err := http.NewRequest("POST", usersServiceURL+"/users/signout", bytes.NewBuffer(signoutJson))
	if err != nil {
		log.Printf("Error creating signout request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending signout request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Error response from users service: %s", string(body))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// 3. Devolver respuesta exitosa
	c.JSON(http.StatusNoContent, nil)
}