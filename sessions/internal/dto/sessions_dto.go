package sessions

import "time"

type CreateSessionDTO struct {
    CircuitKey        int       `json:"circuit_key" binding:"required"`
    CircuitShortName  string    `json:"circuit_short_name" binding:"required"`
    CountryCode       string    `json:"country_code" binding:"required"`
    CountryName       string    `json:"country_name" binding:"required"`
    DateStart         time.Time `json:"date_start" binding:"required"`
    DateEnd           time.Time `json:"date_end" binding:"required"`
    Location          string    `json:"location" binding:"required"`
    SessionKey        int       `json:"session_key" binding:"required"`
    SessionName       string    `json:"session_name" binding:"required"`
    SessionType       string    `json:"session_type" binding:"required"`
    Year              int       `json:"year" binding:"required"`
    
    // Solo para "Race"
    DFastLap          *int      `json:"d_fast_lap,omitempty"`  // ID del piloto con la vuelta más rápida
    VSC               *bool     `json:"vsc,omitempty"`         // Indica si hubo Virtual Safety Car
    SF                *bool     `json:"sf,omitempty"`          // Indica si hubo Safety Car
    DNF               *int      `json:"dnf,omitempty"`         // Número de pilotos que no terminaron
}

type UpdateSessionDTO struct {
    CircuitKey        *int       `json:"circuit_key,omitempty"`
    CircuitShortName  *string    `json:"circuit_short_name,omitempty"`
    CountryCode       *string    `json:"country_code,omitempty"`
    CountryKey        *int       `json:"country_key,omitempty"`
    CountryName       *string    `json:"country_name,omitempty"`
    DateStart         *time.Time `json:"date_start,omitempty"`
    DateEnd           *time.Time `json:"date_end,omitempty"`
    Location          *string    `json:"location,omitempty"`
    SessionKey        *int       `json:"session_key,omitempty"`
    SessionName       *string    `json:"session_name,omitempty"`
    SessionType       *string    `json:"session_type,omitempty"`
    Year              *int       `json:"year,omitempty"`

    // Solo para "Race"
    DFastLap          *int       `json:"d_fast_lap,omitempty"`
    VSC               *bool      `json:"vsc,omitempty"`
    SF                *bool      `json:"sf,omitempty"`
    DNF               *int       `json:"dnf,omitempty"`
}

type ResponseSessionDTO struct {
    ID               uint      `json:"id"`  // Identificador único en la base de datos
    CircuitKey       int       `json:"circuit_key"`
    CircuitShortName string    `json:"circuit_short_name"`
    CountryCode      string    `json:"country_code"`
    CountryName      string    `json:"country_name"`
    DateStart        time.Time `json:"date_start"`
    DateEnd          time.Time `json:"date_end"`
    Location         string    `json:"location"`
    SessionKey       int       `json:"session_key"`  // Identificador lógico o de negocio
    SessionName      string    `json:"session_name"`
    SessionType      string    `json:"session_type"`
    Year             int       `json:"year"`

    // Solo para "Race"
    DFastLap         *int      `json:"d_fast_lap,omitempty"`
    VSC              *bool     `json:"vsc,omitempty"`
    SF               *bool     `json:"sf,omitempty"`
    DNF              *int      `json:"dnf,omitempty"`
}

type DeleteSessionDTO struct {
    SessionID uint `json:"session_id"`
}

type SessionNameAndTypeDTO struct {
	SessionName string `json:"session_name"`
	SessionType string `json:"session_type"`
}

type RaceResultsDTO struct {
    DNF          *int  `json:"dnf,omitempty"`
    VSC          *bool `json:"vsc,omitempty"`
    SF           *bool `json:"sf,omitempty"`
    DFastLap     *int  `json:"d_fast_lap,omitempty"`
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
    DriverNumber   int      `json:"driver_number"`
    LapDuration    *float64 `json:"lap_duration"`  // Tiempo de vuelta en segundos
    LapNumber      int      `json:"lap_number"`    // Número de la vuelta
}

type UpdateDNFDTO struct {
    DNF int `json:"dnf" binding:"required"`
}