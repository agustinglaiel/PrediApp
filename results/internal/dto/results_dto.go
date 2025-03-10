package dto

import "time"

// DTO para crear un nuevo resultado
type CreateResultDTO struct {
	SessionID      int  `json:"session_id" binding:"required"`
	DriverID       int    `json:"driver_id" binding:"required"`
	Position       int     `json:"position" binding:"required"`
	FastestLapTime float64 `json:"fastest_lap_time"`
}

// DTO para actualizar un resultado existente
type UpdateResultDTO struct {
	Position       int     `json:"position,omitempty"`
	FastestLapTime float64 `json:"fastest_lap_time,omitempty"`
}

// DTO para devolver un resultado con detalles del piloto y la sesión
type ResponseResultDTO struct {
	ID             int                  `json:"id"`
	Position       int                  `json:"position"`
	FastestLapTime float64              `json:"fastest_lap_time"`
	Session        ResponseSessionDTO   `json:"session"`
	Driver         ResponseDriverDTO    `json:"driver"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
}

// DTO para los detalles de la sesión asociados al resultado
type ResponseSessionDTO struct {
	ID               int      `json:"id"`
	CircuitShortName string    `json:"circuit_short_name"`
	CountryName      string    `json:"country_name"`
	Location         string    `json:"location"`
	SessionName      string    `json:"session_name"`
	SessionType      string    `json:"session_type"`
	DateStart        time.Time `json:"date_start"`
}

// DTO para los detalles del piloto asociados al resultado
type ResponseDriverDTO struct {
    ID             int    `json:"id"`
    BroadcastName  string `json:"broadcast_name"`
    CountryCode    string `json:"country_code"`
    DriverNumber   int    `json:"driver_number"`
    FirstName      string `json:"first_name"`
    LastName       string `json:"last_name"`
    FullName       string `json:"full_name"`
    NameAcronym    string `json:"name_acronym"`
    TeamName       string `json:"team_name"`
}

type Position struct {
	DriverNumber int       `json:"driver_number"`
	Position     int       `json:"position"`
	// Date         time.Time `json:"date"`
}

type Lap struct {
	LapNumber   int     `json:"lap_number"`
	LapDuration float64 `json:"lap_duration"`
}

// DTO para devolver solo la posición y el ID del piloto
type TopDriverDTO struct {
    Position int `json:"position"`
    DriverID int `json:"driver_id"`
}
