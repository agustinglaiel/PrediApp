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
}

type UpdateSessionDTO struct {
    CircuitKey        int       `json:"circuit_key,omitempty"`
    CircuitShortName  string    `json:"circuit_short_name,omitempty"`
    CountryCode       string    `json:"country_code,omitempty"`
    CountryKey        int       `json:"country_key, omitempty"`
    CountryName       string    `json:"country_name,omitempty"`
    DateStart         time.Time `json:"date_start,omitempty"`
    DateEnd           time.Time `json:"date_end,omitempty"`
    Location          string    `json:"location,omitempty"`
    SessionKey        int       `json:"session_key,omitempty"`
    SessionName       string    `json:"session_name,omitempty"`
    SessionType       string    `json:"session_type,omitempty"`
    Year              int       `json:"year,omitempty"`
}

type ResponseSessionDTO struct {
    ID                uint      `json:"id"`
    CircuitKey        int       `json:"circuit_key"`
    CircuitShortName  string    `json:"circuit_short_name"`
    CountryCode       string    `json:"country_code"`
    CountryName       string    `json:"country_name"`
    DateStart         time.Time `json:"date_start"`
    DateEnd           time.Time `json:"date_end"`
    Location          string    `json:"location"`
    SessionKey        int       `json:"session_key"`
    SessionName       string    `json:"session_name"`
    SessionType       string    `json:"session_type"`
    Year              int       `json:"year"`
}

type DeleteSessionDTO struct {
    SessionID uint `json:"session_id"`
}
