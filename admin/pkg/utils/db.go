package utils

import (
	modelD "admin/internal/model/drivers"
	modelP "admin/internal/model/prodes"
	modelS "admin/internal/model/sessions"
	modelU "admin/internal/model/users"
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
		&modelU.User{},          // Migrar el modelo User
		&modelP.ProdeCarrera{},  // Migrar el modelo ProdeCarrera
		&modelP.ProdeSession{},  // Migrar el modelo ProdeSession
		&modelS.Session{},       // Migrar el modelo Session
		&modelD.Driver{},        // Migrar el modelo Driver
		&modelD.DriverEvent{},   // Migrar el modelo DriverEvent (asegúrate de tenerlo definido)
	)
	if err != nil {
		fmt.Printf("Error migrating database tables: %v\n", err)
		return
	}

	// Verificar y agregar índices o restricciones únicos
	checkAndCreateUniqueConstraint(DB, "sessions", "uni_sessions_session_key", "session_key")

	fmt.Println("Finished migrating database tables")
}

// checkAndCreateUniqueConstraint verifica si existe un índice único y lo crea si no existe
func checkAndCreateUniqueConstraint(db *gorm.DB, tableName, indexName, columnName string) {
	var exists int64
	query := fmt.Sprintf("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = '%s' AND index_name = '%s'", tableName, indexName)
	db.Raw(query).Scan(&exists)

	if exists == 0 {
		err := db.Exec(fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s UNIQUE (%s)", tableName, indexName, columnName)).Error
		if err != nil {
			fmt.Printf("Error creating unique constraint: %v\n", err)
		} else {
			fmt.Printf("Unique constraint `%s` created successfully.\n", indexName)
		}
	} else {
		fmt.Printf("Unique constraint `%s` already exists, skipping creation.\n", indexName)
	}
}