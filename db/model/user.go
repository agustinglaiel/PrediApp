package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID              int            `gorm:"primaryKey" json:"id"`
	FirstName       string         `json:"first_name"`
	LastName        string         `json:"last_name"`
	Username        string         `gorm:"uniqueIndex" json:"username"`
	Email           string         `gorm:"uniqueIndex" json:"email"`
	Password        string         `json:"-"`
	Role            string         `json:"role"`
	Score           int            `gorm:"default:0" json:"score"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	IsEmailVerified bool           `gorm:"default:false" json:"is_email_verified"`
	LastLoginAt     *time.Time     `json:"last_login_at,omitempty"`
	PhoneNumber     string         `json:"phone_number,omitempty"`
	Provider        string         `json:"provider,omitempty"`
	ProviderID      string         `json:"provider_id,omitempty"`
	AvatarURL       string         `json:"avatar_url,omitempty"`
}
