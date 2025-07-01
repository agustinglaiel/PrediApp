package model

import "time"

type Result struct {
	ID             int       `gorm:"primaryKey" json:"id"`
	SessionID      int       `gorm:"index" json:"session_id"`
	Session        *Session  `gorm:"foreignKey:SessionID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"-"` // Relación definida aquí
	DriverID       int       `gorm:"index" json:"driver_id"`
	Driver         *Driver   `gorm:"foreignKey:DriverID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"-"` // Relación definida aquí
	Position       *int      `json:"position"`
	FastestLapTime float64   `json:"fastest_lap_time"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Status         string    `gorm:"type:longtext" json:"status"`
}
