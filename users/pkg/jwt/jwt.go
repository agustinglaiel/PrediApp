package utils

import (
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
	UserID uint `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Genera un token JWT con el user_id
func GenerateJWT(userID uint, role string) (string, e.ApiError) {
	claims := JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)), // Expiración de 12 horas
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", e.NewInternalServerApiError("error generating JWT token", err)
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (*JWTClaims, e.ApiError) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	// Verificar si el token es válido y contiene los claims
	if err != nil {
		// Si hubo un error al analizar el token, envuelve el error en un ApiError
		return nil, e.NewUnauthorizedApiError("Invalid or expired token")
	}

	// Si el token es válido pero no tiene los claims esperados o no es válido, manejamos el error
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, e.NewUnauthorizedApiError("Invalid or expired token")
	}

	return claims, nil
}

// Middleware para verificar el JWT
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extrae el token del encabezado Authorization (Bearer <token>)
		tokenString := strings.Split(authHeader, "Bearer ")[1]
		claims, err := ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Verificar si el rol es admin
		if claims.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return 
		}

		// Guarda el user_id en el contexto de la solicitud para su uso posterior
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}