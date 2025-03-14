package utils

import (
	"fmt"
	"groups/internal/model"
	"groups/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)



var DB * gorm.DB 

func InitDB() (*gorm.DB, error) {
    dsn := config.DBConnectionURL

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("error connecting to the database: %v", err)
    }

    DB = db

    return db, nil
}

func DisconnectDB() {
    sqlDB, err := DB.DB()
    if err != nil {
        fmt.Printf("Error getting DB instance: %v\n", err)
        return
    }
    sqlDB.Close()
}

func StartDbEngine() {
    DB.AutoMigrate(&model.Group{})
    DB.AutoMigrate(&model.GroupXUsers{})
    fmt.Println("Finishing Migration Database Tables")
}

