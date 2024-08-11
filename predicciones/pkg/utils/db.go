package utils

import (
	"fmt"
	"predicciones/internal/model"
	"predicciones/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB 

func InitDB() (*gorm.DB, error){
    dsn := config.DBConnectionURL

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("error connection to the database: %v", err)
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

func StartDbEngine(){
    DB.AutoMigrate(&model.ProdeCarrera{})
    DB.AutoMigrate(&model.ProdeSession{})
    fmt.Println("Finishing Migration Database Tables")
}

/* Esto es para MONGO
var MongoDb *mongo.Database
var client *mongo.Client

// DisconnectDB desconecta la base de datos MongoDB
func DisconnectDB() {
    if client != nil {
        err := client.Disconnect(context.TODO())
        if err != nil {
            log.Fatalf("Error disconnecting from MongoDB: %v", err)
        }
        fmt.Println("Disconnected from MongoDB")
    }
}

// InitDB inicializa la conexión a la base de datos MongoDB
func InitDB() error {
    mongoURI := os.Getenv("MONGO_URI")
    if mongoURI == "" {
        mongoURI = "mongodb://localhost:27017"
    }

    clientOpts := options.Client().ApplyURI(mongoURI)
    cli, err := mongo.Connect(context.TODO(), clientOpts)
    if err != nil {
        return fmt.Errorf("error connecting to MongoDB: %v", err)
    }
    client = cli

    // Verificar la conexión
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    err = client.Ping(ctx, nil)
    if err != nil {
        return fmt.Errorf("error pinging MongoDB: %v", err)
    }

    dbNames, err := client.ListDatabaseNames(context.TODO(), bson.M{})
    if err != nil {
        return fmt.Errorf("error listing database names: %v", err)
    }

    MongoDb = client.Database("prediApp")

    fmt.Println("Available databases:")
    fmt.Println(dbNames)

    return nil
}*/