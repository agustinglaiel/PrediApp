package utils

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	e "users/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Define una estructura para los claims del token
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Genera un token JWT con el user_id
func GenerateJWT(userID uint, role string) (string, e.ApiError) {
	claims := JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(240 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Printf("Error generating JWT token for user ID %d: %v", userID, err)
		return "", e.NewInternalServerApiError("error generating JWT token", err)
	}

	log.Printf("Generated JWT token for user ID %d with role %s", userID, role)
	return tokenString, nil
}

func GenerateRefreshToken() (string, e.ApiError) {
    refreshToken := make([]byte, 32)
    _, err := rand.Read(refreshToken)
    if err != nil {
        return "", e.NewInternalServerApiError("Error generating refresh token", err)
    }
    return base64.URLEncoding.EncodeToString(refreshToken), nil
}


func ValidateJWT(tokenString string) (*JWTClaims, e.ApiError) {
	log.Printf("Validating JWT token: %s", tokenString)
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	// Verificar si el token es válido y contiene los claims
	if err != nil {
		log.Printf("Invalid or expired token: %v", err)
		return nil, e.NewUnauthorizedApiError("Invalid or expired token")
	}

	// Si el token es válido pero no tiene los claims esperados o no es válido, manejamos el error
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		log.Printf("Token is valid but claims are missing or invalid")
		return nil, e.NewUnauthorizedApiError("Invalid or expired token")
	}

	log.Printf("Token validated successfully for user ID %d with role %s", claims.UserID, claims.Role)
	return claims, nil
}

// Middleware para verificar el JWT
func JWTAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            log.Println("Authorization header is missing")
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        // Extrae el token del encabezado Authorization (Bearer <token>)
        tokenParts := strings.Split(authHeader, "Bearer ")
        if len(tokenParts) != 2 {
            log.Println("Invalid Authorization header format")
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
            c.Abort()
            return
        }

        tokenString := tokenParts[1]
        log.Printf("Received JWT token: %s", tokenString) // Print the received token

        claims, err := ValidateJWT(tokenString)
        if err != nil {
            log.Printf("Token validation failed: %v", err)
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
            c.Abort()
            return
        }

        // Verificar si el rol es admin
        if claims.Role != "admin" {
            log.Printf("Access denied: User ID %d does not have admin role", claims.UserID)
            c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
            c.Abort()
            return 
        }

        // Guarda el user_id en el contexto de la solicitud para su uso posterior
        log.Printf("Access granted for user ID %d with admin role", claims.UserID)
        c.Set("user_id", claims.UserID)
        c.Next()
    }
}