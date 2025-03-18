package model

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID              int           `gorm:"primaryKey" json:"id"`
	FirstName       string         `gorm:"size:255" json:"first_name"`
	LastName        string         `gorm:"size:255" json:"last_name"`
	Username        string         `gorm:"size:255;uniqueIndex" json:"username"`
	Email           string         `gorm:"size:255;uniqueIndex" json:"email"`
	Password        string         `gorm:"size:255" json:"-"` // omitir en la respuesta JSON
	Role            string         `gorm:"size:255" json:"role"`
	Score           int            `gorm:"default:0" json:"score"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	IsEmailVerified bool           `gorm:"default:false" json:"is_email_verified"`
	LastLoginAt     *time.Time     `json:"last_login_at,omitempty"`
	PhoneNumber     string         `gorm:"size:20" json:"phone_number,omitempty"`
	Provider        string         `gorm:"size:255" json:"provider,omitempty"`
	ProviderID      string         `gorm:"size:255" json:"provider_id,omitempty"`
	AvatarURL       string         `gorm:"size:255" json:"avatar_url,omitempty"`
}

type RefreshToken struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	UserID    int       `gorm:"index" json:"user_id"`
	Token     string    `gorm:"size:255;uniqueIndex" json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `gorm:"default:false" json:"revoked"`
}