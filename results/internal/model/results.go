package model

import "time"

// Result representa los resultados de una sesión para un piloto
type Result struct {
    ID             int    `gorm:"primaryKey" json:"id"`
    SessionID      int       `json:"session_id"` 
    Session        Session   `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE;OnUpdate:CASCADE;"`
    DriverID       int       `json:"driver_id" gorm:"type:int"`
    Driver         Driver    `gorm:"foreignKey:DriverID;constraint:OnDelete:CASCADE;OnUpdate:CASCADE;"`
    Position       *int       `json:"position"`   
    Status         string    `json:"status"` 
    FastestLapTime float64   `json:"fastest_lap_time"`
    CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// Definir el modelo Driver solo con los campos que necesitas para el preload
type Driver struct {
    ID            int    `gorm:"primaryKey" json:"id"`
    BroadcastName string `json:"broadcast_name"`
    CountryCode   string `json:"country_code"`
    DriverNumber  int    `json:"driver_number"`
    FirstName     string `json:"first_name"`
    LastName      string `json:"last_name"`
    FullName      string `json:"full_name"`
    NameAcronym   string `json:"name_acronym"`
    TeamName      string `json:"team_name"`
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