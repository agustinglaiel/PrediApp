package db

import (
	"fmt"
	"log"
	"time"

	"prediapp.local/db/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init inicializa la conexión a la base de datos y configura el pool de conexiones.
func Init() error {
	// Inicializar configuración y construir URL
	if err := config.Init(); err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}
	// Añadimos multiStatements para que la migración SQL con varios statements funcione
	dsn := config.DBConnectionURL + "&multiStatements=true"

	// Conectar a la base de datos
	dbConn, err := gorm.Open(gormmysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error connecting to DB: %w", err)
	}

	// Configurar el pool de conexiones
	sqlDB, err := dbConn.DB()
	if err != nil {
		return fmt.Errorf("error getting DB instance: %w", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	DB = dbConn
	return nil
}

// DisconnectDB desconecta de la base de datos.
func DisconnectDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("error getting DB instance: %w", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("error closing DB connection: %w", err)
	}
	log.Println("Database connection closed")
	return nil
}

// Migrate ejecuta las migraciones definidas en db/migration.
func Migrate() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("error getting DB instance: %w", err)
	}

	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("error creating migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migration",
		"mysql",
		driver,
	)
	if err != nil {
		return fmt.Errorf("error creating migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error applying migrations: %w", err)
	}

	log.Println("Database migrations applied successfully")
	return nil
}
