package db

import (
	"fmt"

	"prediapp.local/db/config"
	"prediapp.local/db/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init inicializa la conexión a la base de datos
func Init() error {
	db, err := gorm.Open(mysql.Open(config.DBConnectionURL), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error connecting to DB: %w", err)
	}
	DB = db
	return nil
}

// DisconnectDB desconecta de la base de datos
func DisconnectDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		fmt.Printf("Error getting DB instance: %v\n", err)
		return
	}
	sqlDB.Close()
}

// AutoMigrate aplica AutoMigrate sobre todos los modelos
func AutoMigrate() error {
	return DB.AutoMigrate(
		&model.User{},
		&model.Driver{},
		&model.Session{},
		&model.Result{},
		&model.ProdeCarrera{},
		&model.ProdeSession{},
		&model.Group{},
		&model.GroupXUsers{},
		&model.Post{},
	)
}
