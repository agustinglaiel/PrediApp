package model

import (
	"time"

	"gorm.io/gorm"
)

type ProdeCarrera struct {
	ID         int      `gorm:"primaryKey" json:"id"`
	UserID     int      `json:"user_id"`
	EventID    int      `json:"event_id"`
	P1         int      `json:"p1"` // driver_id
	P2         int      `json:"p2"` // driver_id
	P3         int      `json:"p3"` // driver_id
	P4         int      `json:"p4"` // driver_id
	P5         int      `json:"p5"` // driver_id
	FastestLap int      `json:"fastest_lap"` // driver_id
	VSC        bool      `json:"vsc"` // Virtual Safety Car
	SC         bool      `json:"sc"`  // Safety Car
	DNF        int       `json:"dnf"` // Did Not Finish
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type ProdeSession struct {
	ID        int      `gorm:"primaryKey" json:"id"`
	UserID    int      `json:"user_id"`
	EventID   int      `json:"event_id"`
	P1        int      `json:"p1"` // driver_id
	P2        int      `json:"p2"` // driver_id
	P3        int      `json:"p3"` // driver_id
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
