package model

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID           int            `gorm:"primaryKey" json:"id"`
	UserID       int            `json:"user_id"`
	ParentPostID *int           `json:"parent_post_id"`
	Body         string         `json:"body"`
	CreatedAt    time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
