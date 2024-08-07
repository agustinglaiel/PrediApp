package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

// DTO para actualizar un pronóstico de carrera
type UpdateProdeCarreraDTO struct {
    ProdeID    primitive.ObjectID `json:"prode_id"`
    UserID     primitive.ObjectID `json:"user_id"`
    EventID    primitive.ObjectID `json:"event_id"`
    P1         primitive.ObjectID `json:"p1"` // driver_id
	P2         primitive.ObjectID `json:"p2"` // driver_id
	P3         primitive.ObjectID `json:"p3"` // driver_id
    P4         primitive.ObjectID `json:"p4"` // driver_id
    P5         primitive.ObjectID `json:"p5"` // driver_id
    FastestLap string             `json:"fastest_lap"`
    VSC        bool               `json:"vsc"`
    SC         bool               `json:"sc"`
    DNF        int                `json:"dnf"`
}

// DTO para actualizar un pronóstico de sesión que no sea carrera normal
type UpdateProdeSessionDTO struct {
    ProdeID     primitive.ObjectID `json:"prode_id"`
    UserID      primitive.ObjectID `json:"user_id"`
    EventID     primitive.ObjectID `json:"event_id"`
    P1          primitive.ObjectID `json:"p1"` // driver_id
	P2          primitive.ObjectID `json:"p2"` // driver_id
	P3          primitive.ObjectID `json:"p3"` // driver_id
}