package model

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	ID               int            `gorm:"primaryKey" json:"id"`
	WeekendID        int            `json:"weekend_id"`           // Representa el evento o fin de semana de carreras
	CircuitKey       int            `json:"circuit_key"`          // Identificador único del circuito
	CircuitShortName string         `json:"circuit_short_name"`   // Nombre corto del circuito
	CountryCode      string         `json:"country_code"`         // Código del país (por ejemplo, "GBR")
	CountryKey       int            `json:"country_key"`          // Clave única del país
	CountryName      string         `json:"country_name"`         // Nombre del país
	Location         string         `json:"location"`             // Ubicación del circuito
	SessionKey       *int           `json:"session_key"`          // Clave única de la sesión
	SessionName      string         `json:"session_name"`         // Nombre de la sesión (por ejemplo, Carrera, Qualy, Práctica)
	SessionType      string         `json:"session_type"`         // Tipo de sesión (por ejemplo, Carrera, Práctica)
	DateStart        time.Time      `json:"date_start"`           // Fecha y hora de inicio de la sesión
	DateEnd          time.Time      `json:"date_end"`             // Fecha y hora de fin de la sesión
	Year             int            `json:"year"`                 // Año en el que ocurre la sesión
	// DFastLap         *int           `json:"d_fast_lap,omitempty"` // ID del piloto que hizo la vuelta más rápida (solo relevante para "Race")
	VSC              *bool          `json:"vsc,omitempty"`        // Indica si hubo VSC (solo relevante para "Race")
	SF               *bool          `json:"sf,omitempty"`         // Indica si hubo SC (solo relevante para "Race")
	DNF              *int           `json:"dnf,omitempty"`        // Número de pilotos que no terminaron (solo relevante para "Race")
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}