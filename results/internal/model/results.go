package model

import (
	"time"

	"gorm.io/gorm"
)

// Result representa los resultados de una sesión para un piloto
type Result struct {
	ID              int            `gorm:"primaryKey" json:"id"`
	SessionID       int            `gorm:"index" json:"session_id"`
	Session         *Session       `gorm:"foreignKey:SessionID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"-"` // Relación definida aquí
	DriverID        int            `gorm:"index" json:"driver_id"`
	Driver          *Driver        `gorm:"foreignKey:DriverID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"-"` // Relación definida aquí
	Position        *int            `json:"position"`
	FastestLapTime  float64        `json:"fastest_lap_time"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	Status          string         `gorm:"type:longtext" json:"status"`
}

// Definir el modelo Driver solo con los campos que necesitas para el preload
type Driver struct {
	ID            int    `gorm:"primaryKey;type:int" json:"id"`
	BroadcastName string `json:"broadcast_name" gorm:"type:varchar(100)"`
	CountryCode   string `json:"country_code" gorm:"type:varchar(10)"`
	DriverNumber  int    `json:"driver_number"`
	FirstName     string `json:"first_name" gorm:"type:varchar(50);index:idx_driver_name,priority:1"`
	LastName      string `json:"last_name" gorm:"type:varchar(50);index:idx_driver_name,priority:2"`
	FullName      string `json:"full_name" gorm:"type:varchar(100)"`
	NameAcronym   string `json:"name_acronym" gorm:"type:varchar(10)"`
	HeadshotURL   string `json:"headshot_url" gorm:"type:varchar(200)"` // Añadimos el campo de la foto
	TeamName      string `json:"team_name" gorm:"type:varchar(100)"`
	Activo        bool   `json:"activo"`
}

// Definir el modelo Session solo con los campos que necesitas para el preload
type Session struct {
	ID               int            `gorm:"primaryKey" json:"id"`
	WeekendID        int            `json:"weekend_id"`           
	CircuitKey       int            `json:"circuit_key"`          
	CircuitShortName string         `json:"circuit_short_name"`   
	CountryCode      string         `json:"country_code"`         
	CountryKey       int            `json:"country_key"`          
	CountryName      string         `json:"country_name"`         
	Location         string         `json:"location"`             
	SessionKey       *int           `json:"session_key"`          
	SessionName      string         `json:"session_name"`         
	SessionType      string         `json:"session_type"`         
	DateStart        time.Time      `json:"date_start" gorm:"type:timestamp"`           
	DateEnd          time.Time      `json:"date_end" gorm:"type:timestamp"`             
	Year             int            `json:"year"`                 
	VSC              *bool          `json:"vsc,omitempty"`        
	SF               *bool          `json:"sf,omitempty"`         
	DNF              *int           `json:"dnf,omitempty"`        
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}