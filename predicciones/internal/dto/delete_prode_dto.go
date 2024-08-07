package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

// DTO para eliminar un pronóstico
type DeleteProdeDTO struct {
    ProdeID primitive.ObjectID `json:"prode_id"`
    UserID  primitive.ObjectID `json:"user_id"`
}