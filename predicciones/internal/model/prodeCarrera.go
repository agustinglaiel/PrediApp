package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProdeCarrera struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       primitive.ObjectID `bson:"user_id"`
	EventID      primitive.ObjectID `bson:"event_id"`
	P1           string             `bson:"p1"`
	P2           string             `bson:"p2"`
	P3           string             `bson:"p3"`
	P4           string             `bson:"p4"`
	P5           string             `bson:"p5"`
	FastestLap   string             `bson:"fastest_lap"`
	VSC          bool               `bson:"vsc"` // Virtual Safety Car
	SC           bool               `bson:"sc"`  // Safety Car
	DNF          int                `bson:"dnf"` // Did Not Finish
}