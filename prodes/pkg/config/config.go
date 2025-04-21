package config

import (
	"fmt"
	"os"
)

var (
    // Base de datos
    DBUser     = getEnv("DB_USER", "root")
    DBPassword = getEnv("DB_PASSWORD", "1A2g3u4s.")
    DBHost     = getEnv("DB_HOST", "localhost")
    DBPort     = getEnv("DB_PORT", "3306")
    DBName     = getEnv("DB_NAME", "prediapp")

    DBConnectionURL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        DBUser, DBPassword, DBHost, DBPort, DBName)

    // // RabbitMQ
    // RabbitMQUser     = getEnv("RABBITMQ_USER", "prediapp")
    // RabbitMQPassword = getEnv("RABBITMQ_PASSWORD", "1A2g3u4s.")
    // RabbitMQHost     = getEnv("RABBITMQ_HOST", "localhost")
    // RabbitMQPort     = getEnv("RABBITMQ_PORT", "5672")
    // RabbitMQQueue    = getEnv("RABBITMQ_QUEUE", "session_listener_queue")  // Aquí agregamos la cola
    // RabbitMQURL      = fmt.Sprintf("amqp://%s:%s@%s:%s/", RabbitMQUser, RabbitMQPassword, RabbitMQHost, RabbitMQPort)
)

// getEnv obtiene el valor de una variable de entorno o devuelve un valor por defecto si no está establecida
func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}
