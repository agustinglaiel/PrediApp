package utils

import (
	"fmt"
	"time"

	"prediapp.local/prodes/internal/model"
	"prediapp.local/prodes/pkg/config"

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
	sqlDB.SetMaxIdleConns(10)                  // Conexiones inactivas máximas
	sqlDB.SetMaxOpenConns(100)                 // Conexiones máximas abiertas al mismo tiempo
	sqlDB.SetConnMaxLifetime(10 * time.Minute) // Tiempo máximo de vida de una conexión (10 Min)

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
	// Verificar existencia de tablas dependientes
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

	// Migrar las tablas de prodes
	if err := DB.AutoMigrate(&model.ProdeCarrera{}, &model.ProdeSession{}); err != nil {
		return fmt.Errorf("error migrating database tables: %v", err)
	}

	// Establecer zona horaria a UTC
	if err := DB.Exec("SET time_zone = 'UTC'").Error; err != nil {
		return fmt.Errorf("error setting timezone to UTC: %v", err)
	}

	// Agregar Foreign Keys
	if err := addForeignKeys(); err != nil {
		return fmt.Errorf("error adding foreign keys: %v", err)
	}

	fmt.Println("Finishing Migration Database Tables")
	return nil
}

func addForeignKeys() error {
	fkDefinitions := []struct {
		tableName      string
		constraintName string
		foreignKey     string
		referenceTable string
		referenceKey   string
	}{
		{"prode_sessions", "fk_prode_sessions_user", "user_id", "users", "id"},
		{"prode_sessions", "fk_prode_sessions_session", "session_id", "sessions", "id"},
		{"prode_carreras", "fk_prode_carreras_user", "user_id", "users", "id"},
		{"prode_carreras", "fk_prode_carreras_session", "session_id", "sessions", "id"},
	}

	for _, fk := range fkDefinitions {
		var count int64
		query := `
            SELECT COUNT(*)
            FROM information_schema.REFERENTIAL_CONSTRAINTS
            WHERE CONSTRAINT_NAME = ? AND TABLE_NAME = ?;
        `
		if err := DB.Raw(query, fk.constraintName, fk.tableName).Scan(&count).Error; err != nil {
			return fmt.Errorf("error checking foreign key %s: %w", fk.constraintName, err)
		}

		if count == 0 {
			alterQuery := fmt.Sprintf(`
                ALTER TABLE %s
                ADD CONSTRAINT %s
                FOREIGN KEY (%s)
                REFERENCES %s(%s)
                ON DELETE CASCADE
                ON UPDATE CASCADE;
            `, fk.tableName, fk.constraintName, fk.foreignKey, fk.referenceTable, fk.referenceKey)

			if err := DB.Exec(alterQuery).Error; err != nil {
				return fmt.Errorf("error adding foreign key %s: %w", fk.constraintName, err)
			}

			fmt.Printf("Foreign key %s added successfully.\n", fk.constraintName)
		} else {
			fmt.Printf("Foreign key %s already exists, skipping.\n", fk.constraintName)
		}
	}

	return nil
}
