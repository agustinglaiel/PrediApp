package config

import (
	"fmt"
	"os"
)

var (
    DBUser     = getEnv("DB_USER", "root")
    DBPassword = getEnv("DB_PASSWORD", "1A2g3u4s.")
    DBHost     = getEnv("DB_HOST", "localhost")
    DBPort     = getEnv("DB_PORT", "3306")
    DBName     = getEnv("DB_NAME", "prediApp")

    DBConnectionURL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        DBUser, DBPassword, DBHost, DBPort, DBName)
)

// getEnv obtiene el valor de una variable de entorno o devuelve un valor por defecto si no está establecida
func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}
