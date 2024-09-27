package model

import (
	"time"

	"gorm.io/gorm"
)

type ProdeCarrera struct {
	ID         int      	  `gorm:"primaryKey" json:"id"`
	UserID     int      	  `json:"user_id"`
	SessionID  int      	  `json:"session_id"` // foreign key to sessions
	Session    Session  	  `gorm:"foreignKey:SessionID"` // Relación con la tabla Session
	P1         int      	  `json:"p1"` // driver_id
	P2         int      	  `json:"p2"` // driver_id
	P3         int      	  `json:"p3"` // driver_id
	P4         int      	  `json:"p4"` // driver_id
	P5         int       	  `json:"p5"` // driver_id
	FastestLap int      	  `json:"fastest_lap"` // driver_id
	VSC        bool      	  `json:"vsc"` // Virtual Safety Car
	SC         bool      	  `json:"sc"`  // Safety Car
	DNF        int       	  `json:"dnf"` // Did Not Finish
	CreatedAt  time.Time 	  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time 	  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type ProdeSession struct {
	ID        int      `gorm:"primaryKey" json:"id"`
	UserID    int      `json:"user_id"`
	SessionID int      `json:"session_id"` // foreign key to sessions
	Session   Session       `gorm:"foreignKey:SessionID"` // Relación con la tabla Session
	P1        int      `json:"p1"` // driver_id
	P2        int      `json:"p2"` // driver_id
	P3        int      `json:"p3"` // driver_id
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// Session es un modelo simplificado para representar la sesión que se relaciona con los prodes
type Session struct {
	ID               int       `gorm:"primaryKey" json:"id"`
	CircuitShortName string    `json:"circuit_short_name"`
	CountryCode      string    `json:"country_code"`
	CountryName      string    `json:"country_name"`
	DateStart        time.Time `json:"date_start"`
	DateEnd          time.Time `json:"date_end"`
	Location         string    `json:"location"`
}