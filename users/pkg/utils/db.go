package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

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
    clientOpts := options.Client().ApplyURI("mongodb://localhost:27017")
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

    // Cambiar el nombre de la base de datos aquí
    MongoDb = client.Database("prediApp")

    fmt.Println("Available databases:")
    fmt.Println(dbNames)

    return nil
}