package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

// DTO de respuesta para un pronóstico de carrera
type ResponseProdeCarreraDTO struct {
    ID         primitive.ObjectID `json:"id"`
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

// DTO de respuesta para un pronóstico de sesión
type ResponseProdeSessionDTO struct {
    ID          primitive.ObjectID `json:"id"`
    UserID      primitive.ObjectID `json:"user_id"`
    EventID     primitive.ObjectID `json:"event_id"`
    P1          primitive.ObjectID `json:"p1"` // driver_id
	P2          primitive.ObjectID `json:"p2"` // driver_id
	P3          primitive.ObjectID `json:"p3"` // driver_id
}

// son utilizados para estructurar y enviar la información de vuelta al cliente después de realizar 
// una operación en el servidor. Estos DTOs aseguran que solo se devuelvan los datos necesarios 
// en el formato adecuado, mejorando la seguridad y eficiencia de la comunicación entre el backend y el frontend.