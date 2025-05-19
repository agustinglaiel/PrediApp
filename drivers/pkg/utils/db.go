package utils

import (
	modelD "drivers/internal/model"
	"drivers/pkg/config"
	"fmt"

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
    fmt.Println("Database connection closed successfully")
}

// StartDbEngine migrates the database tables
func StartDbEngine() error {
    // Migrar las tablas controladas por este microservicio
    err := DB.AutoMigrate(
        &modelD.Driver{},        // Migrar el modelo Driver
        //&modelD.DriverEvent{}, // Descomentar si DriverEvent existe
    )
    if err != nil {
        return fmt.Errorf("error migrating database tables: %v", err)
    }

    // Establecer la zona horaria de la sesión a UTC
    if err := DB.Exec("SET time_zone = 'UTC'").Error; err != nil {
        return fmt.Errorf("error setting timezone to UTC: %v", err)
    }

    fmt.Println("Finished migrating database tables")
    return nil
}