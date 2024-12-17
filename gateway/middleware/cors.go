package middleware

// Contiene el middleware de CORS, que configura los headers
// necesarios para permitir peticiones desde el frontend y otros orígenes.

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc{
	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != ""  && contains(splitString(allowedOrigins), origin){
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else if allowedOrigins == "*" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
		}
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func splitString(s string) []string {
    var result []string
	if s == "" {
		return result
	}
    for _, part := range strings.Split(s, ",") {
        result = append(result, strings.TrimSpace(part))
    }
    return result
}
	