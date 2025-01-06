package utils

import (
	"fmt"
	"results/internal/model"
	"results/pkg/config"

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
func StartDbEngine() {
    // Migraciones automáticas de GORM
    if err := DB.AutoMigrate(&model.Result{}); err != nil {
        fmt.Printf("Error migrating database tables: %v\n", err)
        return
    }

    // Migraciones personalizadas en SQL en bruto
    addForeignKeys()
    fmt.Println("Finishing Migration Database Tables")
}

// addForeignKeys agrega relaciones personalizadas entre las tablas
func addForeignKeys() {
    // Relación entre `results` y `sessions`
    result := DB.Exec(`
        ALTER TABLE results
        ADD CONSTRAINT fk_results_sessions
        FOREIGN KEY (session_id)
        REFERENCES sessions(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE;
    `)
    if result.Error != nil {
        fmt.Printf("Error adding foreign key for sessions: %v\n", result.Error)
    }

    // Relación entre `results` y `drivers`
    result = DB.Exec(`
        ALTER TABLE results
        ADD CONSTRAINT fk_results_drivers
        FOREIGN KEY (driver_id)
        REFERENCES drivers(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE;
    `)
    if result.Error != nil {
        fmt.Printf("Error adding foreign key for drivers: %v\n", result.Error)
    }
}
