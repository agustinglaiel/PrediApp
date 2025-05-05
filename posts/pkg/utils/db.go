package utils

import (
	"fmt"
	"posts/internal/model"
	"posts/pkg/config"

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
    DB.AutoMigrate(&model.Post{})
    fmt.Println("Finishing Migration Database Tables")
}