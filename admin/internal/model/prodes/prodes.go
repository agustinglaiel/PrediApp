package prodes

import (
	"time"

	"gorm.io/gorm"
)

type ProdeCarrera struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `json:"user_id"`
	EventID    uint      `json:"event_id"`
	P1         uint      `json:"p1"` // driver_id
	P2         uint      `json:"p2"` // driver_id
	P3         uint      `json:"p3"` // driver_id
	P4         uint      `json:"p4"` // driver_id
	P5         uint      `json:"p5"` // driver_id
	FastestLap uint      `json:"fastest_lap"` // driver_id
	VSC        bool      `json:"vsc"` // Virtual Safety Car
	SC         bool      `json:"sc"`  // Safety Car
	DNF        int       `json:"dnf"` // Did Not Finish
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

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
