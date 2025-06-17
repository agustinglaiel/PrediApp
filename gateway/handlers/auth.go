package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID    int    `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Score     int    `json:"score"`
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
	Score		 int    `json:"score"`
	Token        string `json:"token"`
	CreatedAt    string `json:"created_at"`
}

func GenerateTokens(id int, firstName, lastName, username, email, role string, score int, secretKey string) (string, error) {
	claims := Claims{
		UserID:    id,
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		Email:     email,
		Role:      role,
		Score:     score,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString([]byte(secretKey))
}

func LoginHandler(c *gin.Context) {
	usersServiceURL := os.Getenv("USERS_SERVICE_URL")
	secretKey := os.Getenv("JWT_SECRET")

	// 1. Capturar credenciales
	var creds map[string]interface{}
	if err := c.ShouldBindJSON(&creds); err != nil {
		log.Printf("Login: payload inválido: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// 2. Enviar al users service
	body, _ := json.Marshal(creds)
	req, err := http.NewRequest("POST", usersServiceURL+"/users/login", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Login: error creando request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Login: error llamando al servicio de users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer resp.Body.Close()

	// 3. Validar estado
	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Login: users service respondió %d: %s", resp.StatusCode, string(b))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// 4. Leer y parsear la respuesta de users service (sin tokens)
	var userResp UserResponseDTO
	raw, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(raw, &userResp); err != nil {
		log.Printf("Login: no pude parsear respuesta: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// 5. Generar JWT
	token, err := GenerateTokens(
		userResp.ID,
		userResp.FirstName,
		userResp.LastName,
		userResp.Username,
		userResp.Email,
		userResp.Role,
		userResp.Score,
		secretKey,
	)
	if err != nil {
		log.Printf("Login: fallo generando JWT: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// 6. Devolver solo el JWT y la info pública
	c.JSON(http.StatusOK, gin.H{
		"id":          userResp.ID,
		"first_name":  userResp.FirstName,
		"last_name":   userResp.LastName,
		"username":    userResp.Username,
		"email":       userResp.Email,
		"role":        userResp.Role,
		"score":       userResp.Score,
		"token":       token,
		"created_at":  userResp.CreatedAt,
	})
}

func SignupHandler(c *gin.Context) {
	usersServiceURL := os.Getenv("USERS_SERVICE_URL")
	secretKey := os.Getenv("JWT_SECRET")

	// 1. Capturar datos de registro
	var signUp map[string]interface{}
	if err := c.ShouldBindJSON(&signUp); err != nil {
		log.Printf("Signup: payload inválido: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// 2. Enviar al users service
	body, _ := json.Marshal(signUp)
	req, err := http.NewRequest("POST", usersServiceURL+"/users/signup", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Signup: error creando request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Signup: error llamando al servicio de users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer resp.Body.Close()

	// 3. Validar estado
	if resp.StatusCode != http.StatusCreated {
		b, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Signup: users service respondió %d: %s", resp.StatusCode, string(b))
		c.JSON(http.StatusBadRequest, gin.H{"error": "signup failed"})
		return
	}

	// 4. Leer y parsear la respuesta de users service
	var userResp UserResponseDTO
	raw, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(raw, &userResp); err != nil {
		log.Printf("Signup: no pude parsear respuesta: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// 5. Generar JWT
	token, err := GenerateTokens(
		userResp.ID,
		userResp.FirstName,
		userResp.LastName,
		userResp.Username,
		userResp.Email,
		userResp.Role,
		userResp.Score,
		secretKey,
	)
	if err != nil {
		log.Printf("Signup: fallo generando JWT: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// 6. Devolver solo el JWT y la info pública
	c.JSON(http.StatusCreated, gin.H{
		"id":          userResp.ID,
		"first_name":  userResp.FirstName,
		"last_name":   userResp.LastName,
		"username":    userResp.Username,
		"email":       userResp.Email,
		"role":        userResp.Role,
		"score":       userResp.Score,
		"token":       token,
		"created_at":  userResp.CreatedAt,
	})
}

func MeHandler(c *gin.Context) {
	raw, ok := c.Get("claims")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	claims := raw.(*Claims)
	c.JSON(http.StatusOK, gin.H{
		"id":         claims.UserID,
		"first_name": claims.FirstName,
		"last_name":  claims.LastName,
		"username":   claims.Username,
		"email":      claims.Email,
		"role":       claims.Role,
		"score":      claims.Score,
	})
}