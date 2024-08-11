package model

import (
	"time"

	"gorm.io/gorm"
)

type ProdeSession struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	EventID   uint      `json:"event_id"`
	P1        uint      `json:"p1"` // driver_id
	P2        uint      `json:"p2"` // driver_id
	P3        uint      `json:"p3"` // driver_id
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
