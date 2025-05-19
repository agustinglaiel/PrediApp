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
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // Evita que GORM cree tablas o claves foráneas automáticamente
	})
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
	// Verificar que las tablas dependientes existan
	for _, table := range []string{"users"} {
		var exists int64
		query := fmt.Sprintf("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = '%s'", table)
		if err := DB.Raw(query).Scan(&exists).Error; err != nil {
			return fmt.Errorf("error checking if table %s exists: %v", table, err)
		}
		if exists == 0 {
			return fmt.Errorf("table %s does not exist, cannot migrate posts tables", table)
		}
	}

	// Migrar solo las tablas controladas por posts
	if err := DB.AutoMigrate(&model.Post{}); err != nil {
		return fmt.Errorf("error migrating database tables: %v", err)
	}

	// Establecer la zona horaria de la sesión a UTC
	if err := DB.Exec("SET time_zone = 'UTC'").Error; err != nil {
		return fmt.Errorf("error setting timezone to UTC: %v", err)
	}

	// Verificar y agregar el índice FULLTEXT si no existe
	var indexExists int64
	query := `
		SELECT COUNT(*)
		FROM information_schema.statistics
		WHERE table_schema = DATABASE()
		AND table_name = 'posts'
		AND index_name = 'idx_posts_body';
	`
	if err := DB.Raw(query).Scan(&indexExists).Error; err != nil {
		fmt.Printf("Debug: Error checking index existence: %v\n", err) // Depuración
		return fmt.Errorf("error checking if index idx_posts_body exists: %v", err)
	}

	if indexExists == 0 {
		if err := DB.Exec(`
			ALTER TABLE posts
			ADD FULLTEXT INDEX idx_posts_body (body);
		`).Error; err != nil {
			return fmt.Errorf("error creating FULLTEXT index idx_posts_body: %v", err)
		}
		fmt.Println("FULLTEXT index idx_posts_body created successfully")
	} else {
		fmt.Println("FULLTEXT index idx_posts_body already exists, skipping creation")
	}

	fmt.Println("Finishing Migration Database Tables")
	return nil
}