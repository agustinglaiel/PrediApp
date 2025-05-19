package utils

import (
	"fmt"
	"sessions/internal/model"
	"sessions/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the database connection
func InitDB() (*gorm.DB, error) {
    dsn := config.DBConnectionURL

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("error connecting to the database: %v", err)
    }

    DB = db

    return db, nil
}

// DisconnectDB disconnects from the database
func DisconnectDB() {
    sqlDB, err := DB.DB()
    if err != nil {
        fmt.Printf("Error getting DB instance: %v\n", err)
        return
    }
    sqlDB.Close()
}

// StartDbEngine migrates the database tables
func StartDbEngine() error {
    // Migrar las tablas
    if err := DB.AutoMigrate(&model.Session{}); err != nil {
        return fmt.Errorf("error during database migration: %v", err)
    }

    // Establecer la zona horaria de la sesión a UTC
    if err := DB.Exec("SET time_zone = 'UTC'").Error; err != nil {
        return fmt.Errorf("error setting timezone to UTC: %v", err)
    }

    fmt.Println("Finishing Migration Database Tables")
    return nil
}
