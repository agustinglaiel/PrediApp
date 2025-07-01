package model

import (
	"time"

	"gorm.io/gorm"
)

type ProdeCarrera struct {
	ID        int            `gorm:"primaryKey" json:"id"`
	UserID    int            `json:"user_id"`
	SessionID int            `json:"session_id"`
	P1        int            `json:"p1"`
	P2        int            `json:"p2"`
	P3        int            `json:"p3"`
	P4        int            `json:"p4"`
	P5        int            `json:"p5"`
	VSC       bool           `json:"vsc"`
	SC        bool           `json:"sc"`
	DNF       int            `json:"dnf"`
	Score     int            `gorm:"default:0" json:"score"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
