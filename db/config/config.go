package config

import (
	"fmt"
	"os"
)

// DBConnectionURL construye la URL de conexión a la base de datos usando variables de entorno.
var DBConnectionURL string

// Init inicializa la configuración validando las variables de entorno necesarias.
func Init() error {
	requiredEnvVars := []string{
		"DB_USER",
		"DB_PASS",
		"DB_HOST",
		"DB_PORT",
		"DB_NAME",
	}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			return fmt.Errorf("%s environment variable is not set", envVar)
		}
	}

	DBConnectionURL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	return nil
}
