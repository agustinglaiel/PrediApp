package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProdeCarrera struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       primitive.ObjectID `bson:"user_id"`
	EventID      primitive.ObjectID `bson:"event_id"`
	P1           primitive.ObjectID `bson:"p1"` // driver_id
	P2           primitive.ObjectID `bson:"p2"` // driver_id
	P3           primitive.ObjectID `bson:"p3"` // driver_id
	P4           primitive.ObjectID `bson:"p4"` // driver_id
	P5           primitive.ObjectID `bson:"p5"` // driver_id
	FastestLap   primitive.ObjectID `bson:"fastest_lap"` // driver_id
	VSC          bool               `bson:"vsc"` // Virtual Safety Car
	SC           bool               `bson:"sc"`  // Safety Car
	DNF          int                `bson:"dnf"` // Did Not Finish
}
