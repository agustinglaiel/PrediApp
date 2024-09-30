package utils

import (
	"fmt"
	"prodes/internal/model"
	"prodes/pkg/config"
	"time"

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

    // Obtener la conexión SQL subyacente
    sqlDB, err := db.DB()
    if err != nil {
        return nil, fmt.Errorf("error getting sql.DB from gorm: %v", err)
    }

    // Configurar el pool de conexiones
    sqlDB.SetMaxIdleConns(10)                    // Conexiones inactivas máximas
    sqlDB.SetMaxOpenConns(100)                   // Conexiones máximas abiertas al mismo tiempo
    sqlDB.SetConnMaxLifetime(10 * time.Minute)   // Tiempo máximo de vida de una conexión (10 Min)

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
func StartDbEngine() {
    DB.AutoMigrate(&model.ProdeCarrera{})
    DB.AutoMigrate(&model.ProdeSession{})
    fmt.Println("Finishing Migration Database Tables")
}
