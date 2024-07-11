package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    FirstName       string             `bson:"first_name" json:"first_name"`
    LastName        string             `bson:"last_name" json:"last_name"`
    Username        string             `bson:"username" json:"username"`
    Email           string             `bson:"email" json:"email"`
    Password        string             `bson:"password" json:"-"`
    Role            string             `bson:"role" json:"role"`
    Score           int                `bson:"score" json:"score"`
    CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
    DeletedAt       *time.Time         `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
    IsActive        bool               `bson:"is_active" json:"is_active"`
    IsEmailVerified bool               `bson:"is_email_verified" json:"is_email_verified"`
    LastLoginAt     *time.Time         `bson:"last_login_at,omitempty" json:"last_login_at,omitempty"`
    PhoneNumber     string             `bson:"phone_number,omitempty" json:"phone_number,omitempty"`
}

