package model

import (
	"time"

	"gorm.io/gorm"
)

type ProdeSession struct {
	ID        int            `gorm:"primaryKey" json:"id"`
	UserID    int            `gorm:"index;not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
	SessionID int            `gorm:"index;not null" json:"session_id"`
	Session   Session        `gorm:"foreignKey:SessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"session"`
	P1        int            `json:"p1"`
	DriverP1  Driver         `gorm:"foreignKey:P1;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"driver_p1"`
	P2        int            `json:"p2"`
	DriverP2  Driver         `gorm:"foreignKey:P2;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"driver_p2"`
	P3        int            `json:"p3"`
	DriverP3  Driver         `gorm:"foreignKey:P3;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"driver_p3"`
	Score     int            `gorm:"default:0" json:"score"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
