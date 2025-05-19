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

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        DisableForeignKeyConstraintWhenMigrating: true, // Evita que GORM cree tablas o claves foráneas automáticamente
    })
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
    fmt.Println("Database connection closed successfully")
}

// StartDbEngine migrates the database tables
func StartDbEngine() error {
    // Verificar que las tablas dependientes existan
    for _, table := range []string{"users", "sessions", "drivers"} {
        var exists int64
        query := fmt.Sprintf("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = '%s'", table)
        if err := DB.Raw(query).Scan(&exists).Error; err != nil {
            return fmt.Errorf("error checking if table %s exists: %v", table, err)
        }
        if exists == 0 {
            return fmt.Errorf("table %s does not exist, cannot migrate prodes tables", table)
        }
    }

    // Migrar solo las tablas controladas por prodes
    if err := DB.AutoMigrate(
        &model.ProdeCarrera{},
        &model.ProdeSession{},
    ); err != nil {
        return fmt.Errorf("error migrating database tables: %v", err)
    }

    // Establecer la zona horaria de la sesión a UTC
    if err := DB.Exec("SET time_zone = 'UTC'").Error; err != nil {
        return fmt.Errorf("error setting timezone to UTC: %v", err)
    }

    fmt.Println("Finishing Migration Database Tables")
    return nil
}