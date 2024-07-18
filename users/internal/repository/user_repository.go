package repository

import (
	"context"
	"log"
	"users/internal/model"
	e "users/pkg/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// userRepository es una estructura vacía que implementa la interfaz UserRepository
type userRepository struct{}

// UserRepository define los métodos que deben ser implementados por el repositorio de usuarios
type UserRepository interface {
    CreateUser(ctx context.Context, user *model.User) e.ApiError
    GetUserByEmail(ctx context.Context, email string) (*model.User, e.ApiError)
    GetUserByUsername(ctx context.Context, username string) (*model.User, e.ApiError)
    GetUserByID(ctx context.Context, id primitive.ObjectID) (*model.User, e.ApiError)
    GetUsers(ctx context.Context) ([]*model.User, e.ApiError)
    UpdateUserByID(ctx context.Context, id primitive.ObjectID, user *model.User) e.ApiError
    UpdateUserByUsername(ctx context.Context, username string, user *model.User) e.ApiError
    DeleteUserByID(ctx context.Context, id primitive.ObjectID) e.ApiError
    DeleteUserByUsername(ctx context.Context, username string) e.ApiError
}

// NewUserRepository crea una nueva instancia de userRepository
func NewUserRepository() UserRepository {
    return &userRepository{}
}

// CreateUser inserta un nuevo usuario en la base de datos
func (r *userRepository) CreateUser(ctx context.Context, user *model.User) e.ApiError {
    db := e.MongoDb
    //log.Printf("Attempting to insert user into database: %+v", user)
    result, err := db.Collection("users").InsertOne(ctx, user)
    if err != nil {
        log.Printf("Error creating user: %v", err)
        return e.NewInternalServerApiError("error creating user", err)
    }
    //log.Printf("User inserted with ID: %v", result.InsertedID)
    
    // Agregar log para confirmar que el usuario se ha insertado
    var insertedUser model.User
    err = db.Collection("users").FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&insertedUser)
    if err != nil {
        log.Printf("Error verifying inserted user: %v", err)
    } /*else {
        log.Printf("Inserted user: %+v", insertedUser)
    }*/
    
    return nil
}

// GetUserByEmail obtiene un usuario por su email
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, e.ApiError) {
    var user model.User
    db := e.MongoDb
    //log.Printf("Searching for user by email: %s", email)
    err := db.Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            log.Printf("User not found: %s", email)
            return nil, e.NewNotFoundApiError("user not found")
        }
        log.Printf("Error finding user by email: %v", err)
        return nil, e.NewInternalServerApiError("error finding user by email", err)
    }
    log.Printf("User found: %+v", user)
    return &user, nil
}

// GetUserByUsername obtiene un usuario por su nombre de usuario
func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, e.ApiError) {
    var user model.User
    db := e.MongoDb
    err := db.Collection("users").FindOne(ctx, bson.M{"username": username}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, e.NewNotFoundApiError("user not found")
        }
        return nil, e.NewInternalServerApiError("error finding user by username", err)
    }
    return &user, nil
}

// GetUserByID obtiene un usuario por su ID
func (r *userRepository) GetUserByID(ctx context.Context, id primitive.ObjectID) (*model.User, e.ApiError) {
    var user model.User
    db := e.MongoDb
    err := db.Collection("users").FindOne(ctx, bson.M{"_id": id}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, e.NewNotFoundApiError("user not found")
        }
        return nil, e.NewInternalServerApiError("error finding user by ID", err)
    }
    return &user, nil
}

// GetUsers obtiene todos los usuarios
func (r *userRepository) GetUsers(ctx context.Context) ([]*model.User, e.ApiError) {
    var users []*model.User
    db := e.MongoDb
    cursor, err := db.Collection("users").Find(ctx, bson.M{})
    if err != nil {
        return nil, e.NewInternalServerApiError("error finding users", err)
    }
    defer cursor.Close(ctx)
    for cursor.Next(ctx) {
        var user model.User
        if err := cursor.Decode(&user); err != nil {
            return nil, e.NewInternalServerApiError("error decoding user", err)
        }
        users = append(users, &user)
    }
    if err := cursor.Err(); err != nil {
        return nil, e.NewInternalServerApiError("cursor error", err)
    }
    return users, nil
}

// UpdateUserByID actualiza un usuario por su ID en la base de datos
func (r *userRepository) UpdateUserByID(ctx context.Context, id primitive.ObjectID, user *model.User) e.ApiError {
    db := e.MongoDb
    filter := bson.M{"_id": id}
    update := bson.M{
        "$set": user,
    }
    _, err := db.Collection("users").UpdateOne(ctx, filter, update)
    if err != nil {
        return e.NewInternalServerApiError("error updating user by ID", err)
    }
    return nil
}

// UpdateUserByUsername actualiza un usuario por su nombre de usuario en la base de datos
func (r *userRepository) UpdateUserByUsername(ctx context.Context, username string, user *model.User) e.ApiError {
    db := e.MongoDb
    filter := bson.M{"username": username}
    update := bson.M{
        "$set": user,
    }
    _, err := db.Collection("users").UpdateOne(ctx, filter, update)
    if err != nil {
        return e.NewInternalServerApiError("error updating user by username", err)
    }
    return nil
}

// DeleteUserByID elimina un usuario por su ID de la base de datos
func (r *userRepository) DeleteUserByID(ctx context.Context, id primitive.ObjectID) e.ApiError {
    db := e.MongoDb
    filter := bson.M{"_id": id}
    _, err := db.Collection("users").DeleteOne(ctx, filter)
    if err != nil {
        return e.NewInternalServerApiError("error deleting user by ID", err)
    }
    return nil
}

// DeleteUserByUsername elimina un usuario por su nombre de usuario de la base de datos
func (r *userRepository) DeleteUserByUsername(ctx context.Context, username string) e.ApiError {
    db := e.MongoDb
    filter := bson.M{"username": username}
    _, err := db.Collection("users").DeleteOne(ctx, filter)
    if err != nil {
        return e.NewInternalServerApiError("error deleting user by username", err)
    }
    return nil
}