package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID              int            `gorm:"primaryKey" json:"id"`
	FirstName       string         `gorm:"size:255" json:"first_name"`
	LastName        string         `gorm:"size:255" json:"last_name"`
	Username        string         `gorm:"size:255;uniqueIndex" json:"username"`
	Email           string         `gorm:"size:255;uniqueIndex" json:"email"`
	Password        string         `gorm:"size:255" json:"-"`
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
	ImagenPerfil    []byte         `gorm:"type:mediumblob" json:"imagen_perfil,omitempty"`
	ImagenMimeType  string         `gorm:"size:50" json:"imagen_mime_type,omitempty"`
}
