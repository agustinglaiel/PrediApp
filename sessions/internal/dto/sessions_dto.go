package dto

import "time"

type CreateSessionDTO struct {
	WeekendID        int       `json:"weekend_id" binding:"required"`
	CircuitKey       int       `json:"circuit_key" binding:"required"`
	CircuitShortName string    `json:"circuit_short_name" binding:"required"`
	CountryCode      string    `json:"country_code" binding:"required"`
	CountryName      string    `json:"country_name" binding:"required"`
	DateStart        time.Time `json:"date_start" binding:"required"`
	DateEnd          time.Time `json:"date_end" binding:"required"`
	Location         string    `json:"location" binding:"required"`
	SessionKey       *int      `json:"session_key"`
	SessionName      string    `json:"session_name" binding:"required"`
	SessionType      string    `json:"session_type" binding:"required"`
	Year             int       `json:"year" binding:"required"`

	// Solo para "Race"
	DFastLap *int  `json:"d_fast_lap,omitempty"` // ID del piloto con la vuelta más rápida
	VSC      *bool `json:"vsc,omitempty"`        // Indica si hubo Virtual Safety Car
	SF       *bool `json:"sf,omitempty"`         // Indica si hubo Safety Car
	DNF      *int  `json:"dnf,omitempty"`        // Número de pilotos que no terminaron
}

type UpdateSessionDTO struct {
	WeekendID        *int       `json:"weekend_id,omitempty"`
	CircuitKey       *int       `json:"circuit_key,omitempty"`
	CircuitShortName *string    `json:"circuit_short_name,omitempty"`
	CountryCode      *string    `json:"country_code,omitempty"`
	CountryKey       *int       `json:"country_key,omitempty"`
	CountryName      *string    `json:"country_name,omitempty"`
	DateStart        *time.Time `json:"date_start,omitempty"`
	DateEnd          *time.Time `json:"date_end,omitempty"`
	Location         *string    `json:"location,omitempty"`
	SessionKey       *int       `json:"session_key,omitempty"`
	SessionName      *string    `json:"session_name,omitempty"`
	SessionType      *string    `json:"session_type,omitempty"`
	Year             *int       `json:"year,omitempty"`

	// Solo para "Race"
	DFastLap *int  `json:"d_fast_lap,omitempty"`
	VSC      *bool `json:"vsc,omitempty"`
	SF       *bool `json:"sf,omitempty"`
	DNF      *int  `json:"dnf,omitempty"`
}

type ResponseSessionDTO struct {
	ID               int       `json:"id"`
	WeekendID        int       `json:"weekend_id"`
	CircuitKey       int       `json:"circuit_key"`
	CircuitShortName string    `json:"circuit_short_name"`
	CountryCode      string    `json:"country_code"`
	CountryName      string    `json:"country_name"`
	DateStart        time.Time `json:"date_start"`
	DateEnd          time.Time `json:"date_end"`
	Location         string    `json:"location"`
	SessionKey       *int      `json:"session_key"`
	SessionName      string    `json:"session_name"`
	SessionType      string    `json:"session_type"`
	Year             int       `json:"year"`
	VSC              *bool     `json:"vsc"`
	SF               *bool     `json:"sf"`
	DNF              *int      `json:"dnf"`
}

type DeleteSessionDTO struct {
	SessionID int `json:"session_id"`
}

type SessionNameAndTypeDTO struct {
	SessionName string `json:"session_name"`
	SessionType string `json:"session_type"`
}

type RaceResultsDTO struct {
	DNF      *int  `json:"dnf,omitempty"`
	VSC      *bool `json:"vsc,omitempty"`
	SF       *bool `json:"sf,omitempty"`
	DFastLap *int  `json:"d_fast_lap,omitempty"`
}

// RaceControlEvent representa un evento de control de carrera como SC, VSC o banderas
type RaceControlEvent struct {
	Category     string `json:"category"`
	Date         string `json:"date"`
	DriverNumber *int   `json:"driver_number,omitempty"`
	Flag         string `json:"flag,omitempty"`
	LapNumber    *int   `json:"lap_number,omitempty"`
	MeetingKey   int    `json:"meeting_key"`
	Message      string `json:"message"`
	SessionKey   int    `json:"session_key"`
}

type LapData struct {
	DriverNumber int      `json:"driver_number"`
	LapDuration  *float64 `json:"lap_duration"` // Tiempo de vuelta en segundos
	LapNumber    int      `json:"lap_number"`   // Número de la vuelta
}

type UpdateDNFDTO struct {
	DNF int `json:"dnf" binding:"required"`
}

// SessionKeyResponseDTO representa la respuesta que incluye el session_key, date_start, y date_end
type SessionKeyResponseDTO struct {
	SessionKey *int       `json:"session_key"`
	DateStart  *time.Time `json:"date_start"`
	DateEnd    *time.Time `json:"date_end"`
	CountryKey *int       `json:"country_key"`
	CircuitKey *int       `json:"circuit_key"`
}

type UpdateSessionKeyDTO struct {
	Location    string `json:"location" binding:"required"`
	SessionName string `json:"session_name" binding:"required"`
	SessionType string `json:"session_type" binding:"required"`
	Year        int    `json:"year" binding:"required"`
}

type UpdateSessionDataDTO struct {
	Location    string `json:"location" binding:"required"`
	SessionName string `json:"session_name" binding:"required"`
	SessionType string `json:"session_type" binding:"required"`
	Year        int    `json:"year" binding:"required"`
}

// FastestLapDTO representa el piloto con la vuelta más rápida de una sesión
type FastestLapDTO struct {
	ID             int                `json:"id"`
	Position       int                `json:"position"`
	FastestLapTime float64            `json:"fastest_lap_time"`
	Session        ResponseSessionDTO `json:"session"` // Información de la sesión
	Driver         ResponseDriverDTO  `json:"driver"`  // Información del piloto
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

//ACA MISMO PUDE USAR ResponseSessionDTO DEFINIDO ARRIBA

// ResponseDriverDTO representa la información del piloto
type ResponseDriverDTO struct {
	ID          int    `json:"driver_id"`    // ID del piloto
	FirstName   string `json:"first_name"`   // Nombre del piloto
	LastName    string `json:"last_name"`    // Apellido del piloto
	FullName    string `json:"full_name"`    // Nombre completo del piloto
	NameAcronym string `json:"name_acronym"` // Acrónimo del nombre del piloto
	TeamName    string `json:"team_name"`    // Nombre del equipo del piloto
}
