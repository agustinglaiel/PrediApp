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
    // Verificar y agregar la clave foránea para `results -> sessions`
    err := addForeignKeyIfNotExists("results", "fk_results_sessions", "session_id", "sessions", "id")
    if err != nil {
        fmt.Printf("Error adding foreign key for sessions: %v\n", err)
    }

    // Verificar y agregar la clave foránea para `results -> drivers`
    err = addForeignKeyIfNotExists("results", "fk_results_drivers", "driver_id", "drivers", "id")
    if err != nil {
        fmt.Printf("Error adding foreign key for drivers: %v\n", err)
    }
}

// addForeignKeyIfNotExists verifica si existe una clave foránea y la agrega si no existe
func addForeignKeyIfNotExists(tableName, constraintName, foreignKey, referenceTable, referenceKey string) error {
    // Verificar si la clave foránea ya existe
    query := `
        SELECT COUNT(*)
        FROM information_schema.REFERENTIAL_CONSTRAINTS
        WHERE CONSTRAINT_NAME = ? AND TABLE_NAME = ?;
    `
    var count int64
    result := DB.Raw(query, constraintName, tableName).Scan(&count)
    if result.Error != nil {
        return fmt.Errorf("error checking foreign key %s: %w", constraintName, result.Error)
    }

    if count == 0 {
        // Agregar la clave foránea si no existe
        alterQuery := fmt.Sprintf(`
            ALTER TABLE %s
            ADD CONSTRAINT %s
            FOREIGN KEY (%s)
            REFERENCES %s(%s)
            ON DELETE CASCADE
            ON UPDATE CASCADE;
        `, tableName, constraintName, foreignKey, referenceTable, referenceKey)

        result = DB.Exec(alterQuery)
        if result.Error != nil {
            return fmt.Errorf("error adding foreign key %s: %w", constraintName, result.Error)
        }
    }

    return nil
}
