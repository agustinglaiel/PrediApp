package model

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID           int            `gorm:"primaryKey" json:"id"`
	UserID       int            `gorm:"index" json:"user_id"`
	User         *User          `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"-"`
	ParentPostID *int           `gorm:"index;foreignKey:ParentPostID;references:ID" json:"parent_post_id"`
	Body         string         `gorm:"type:varchar(500);not null" json:"body"`
	CreatedAt    time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Children     []*Post        `gorm:"-" json:"children"`
}

// El campo Children es un campo virtual (no se almacena en la base de datos,
// por eso tiene la etiqueta gorm:"-") que se usa para construir la estructura
// de hilos (es decir, un árbol de posts y comentarios) al devolver datos
// al cliente.

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

