package model

import "time"

type Result struct {
	ID             int       `gorm:"primaryKey" json:"id"`
	SessionID      int       `json:"session_id"`
	DriverID       int       `json:"driver_id"`
	Position       *int      `json:"position"`
	FastestLapTime float64   `json:"fastest_lap_time"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
