package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProdeQualySprint struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       primitive.ObjectID `bson:"user_id"`
	EventID      primitive.ObjectID `bson:"event_id"`
	P1           string             `bson:"p1"`
	P2           string             `bson:"p2"`
	P3           string             `bson:"p3"`
}