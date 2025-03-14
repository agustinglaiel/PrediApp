package handlers

import (
	"bytes"
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

// Directorio para los manejadores de rutas específicas del gateway (ej: /login),
// que se encargan de la comunicación con los microservicios y la generación de tokens JWT.

// Dentro de handlers, creamos un archivo auth.go para manejar toda
// la lógica de autenticación, como el login, el refresh token, etc.

// Estructura para el JWT Payload
type Claims struct {
	UserID int   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Estructura para la respuesta de Login y Signup
type UserResponseDTO struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Token     string `json:"token"` // JWT token
	CreatedAt string `json:"created_at"`
}

func GenerateJWT(userID int, role string, secretKey string) (string, error) {
    claims := Claims{
        UserID: userID,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(2400 * time.Hour)),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(secretKey))
    if err != nil {
        log.Printf("Error generating JWT token for user ID %d: %v", userID, err)
        return "", fmt.Errorf("error generating JWT token: %w", err)
    }
    log.Printf("Generated JWT token for user ID %d with role %s", userID, role)
    return tokenString, nil
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

	// 4. Generar el Token JWT
	token, err := GenerateJWT(userResponse.ID, userResponse.Role, secretKey)
	if err != nil {
		log.Printf("Error generating JWT token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	userResponse.Token = token

	// 5. Devolver la Respuesta con el token JWT
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

	// 4. Generar el Token JWT
	token, err := GenerateJWT(userResponse.ID, userResponse.Role, secretKey)
	if err != nil {
		log.Printf("Error generating JWT token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	userResponse.Token = token

	// 5. Devolver la Respuesta con el token JWT
	c.JSON(http.StatusCreated, userResponse)
}