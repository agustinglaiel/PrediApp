package utils

import (
	"admin/internal/model/drivers"
	modelP "admin/internal/model/prodes"
	modelS "admin/internal/model/sessions"
	"admin/internal/model/users"
	"admin/pkg/config"
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
}

// StartDbEngine migrates the database tables
func StartDbEngine() {
    err := DB.AutoMigrate(
        &users.User{},        // Migrar el modelo User
        &modelP.ProdeCarrera{},  // Migrar el modelo ProdeCarrera
        &modelP.ProdeSession{},  // Migrar el modelo ProdeSession
        &modelS.Session{},  // Migrar el modelo Session
        &drivers.Driver{},    // Migrar el modelo Driver
    )
    if err != nil {
        fmt.Printf("Error migrating database tables: %v\n", err)
        return
    }

    fmt.Println("Finished migrating database tables")
}