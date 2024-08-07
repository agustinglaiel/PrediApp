package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

// DTO para obtener pronósticos por usuario o evento
type GetProdesDTO struct {
    UserID  primitive.ObjectID `json:"user_id,omitempty"`
    EventID primitive.ObjectID `json:"event_id,omitempty"`
}