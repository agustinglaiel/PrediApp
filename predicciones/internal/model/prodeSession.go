package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProdeSession struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       primitive.ObjectID `bson:"user_id"`
	EventID      primitive.ObjectID `bson:"event_id"`
	P1           primitive.ObjectID `bson:"p1"` // driver_id
	P2           primitive.ObjectID `bson:"p2"` // driver_id
	P3           primitive.ObjectID `bson:"p3"` // driver_id
}
