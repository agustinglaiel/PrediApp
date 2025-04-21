package model

import "time"

// Result representa los resultados de una sesión para un piloto
type Result struct {
    ID             int    `gorm:"primaryKey" json:"id"`
    SessionID      int       `json:"session_id"` // Foreign key to sessions
    Session        Session   `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE;OnUpdate:CASCADE;"`
    DriverID       int       `json:"driver_id" gorm:"type:int"`  // Foreign key to drivers
    Driver         Driver    `gorm:"foreignKey:DriverID;constraint:OnDelete:CASCADE;OnUpdate:CASCADE;"`
    Position       *int       `json:"position"`   // Posición del piloto en la carrera
    Status         string    `json:"status"`     // Estado del piloto en la carrera
    FastestLapTime float64   `json:"fastest_lap_time"` // Duración de la vuelta rápida en segundos (con decimales)
    CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// Definir el modelo Driver solo con los campos que necesitas para el preload
type Driver struct {
    ID          int    `json:"id"`
    FirstName   string `json:"first_name"`
    LastName    string `json:"last_name"`
    FullName    string `json:"full_name"`
    NameAcronym string `json:"name_acronym"`
    TeamName    string `json:"team_name"`
}

// Definir el modelo Session solo con los campos que necesitas para el preload
type Session struct {
    ID               int       `json:"id"`
    CircuitShortName string    `json:"circuit_short_name"`
    CountryName      string    `json:"country_name"`
	Location         string    `json:"location"`
    SessionName      string    `json:"session_name"`
    SessionType      string    `json:"session_type"`
    DateStart        time.Time `json:"date_start"`
}