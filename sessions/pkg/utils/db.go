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
    fmt.Println("Database connection closed successfully")
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

    // Verificar y agregar índices o restricciones únicos
    if err := checkAndCreateUniqueConstraint(DB, "sessions", "uni_sessions_session_key", "session_key"); err != nil {
        return fmt.Errorf("error creating unique constraint: %v", err)
    }

    fmt.Println("Finishing Migration Database Tables")
    return nil
}

func checkAndCreateUniqueConstraint(db *gorm.DB, tableName, indexName, columnName string) error {
    var exists int64
    query := fmt.Sprintf("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = '%s' AND index_name = '%s'", tableName, indexName)
    if err := db.Raw(query).Scan(&exists).Error; err != nil {
        return fmt.Errorf("error checking unique constraint %s: %v", indexName, err)
    }

    if exists == 0 {
        err := db.Exec(fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s UNIQUE (%s)", tableName, indexName, columnName)).Error
        if err != nil {
            return fmt.Errorf("error creating unique constraint %s: %v", indexName, err)
        }
        fmt.Printf("Unique constraint `%s` created successfully.\n", indexName)
    } else {
        fmt.Printf("Unique constraint `%s` already exists, skipping creation.\n", indexName)
    }
    return nil
}