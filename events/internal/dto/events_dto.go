package dto

import "time"

// CreateEventDTO - DTO para la creación de un nuevo evento
type CreateEventDTO struct {
    SessionID            int       `json:"session_id"`
    Date                 time.Time `json:"date"`
    RaceResultID         *int      `json:"race_result_id,omitempty"`
    SprintRaceResultID   *int      `json:"sprint_race_result_id,omitempty"`
    QualyResultID        *int      `json:"qualy_result_id,omitempty"`
    SprintQualyResultID  *int      `json:"sprint_qualy_result_id,omitempty"`
    FP1ID                *int      `json:"fp1_id,omitempty"`
    FP2ID                *int      `json:"fp2_id,omitempty"`
    FP3ID                *int      `json:"fp3_id,omitempty"`
}

// UpdateEventDTO - DTO para la actualización de un evento existente
type UpdateEventDTO struct {
    Date                 *time.Time `json:"date,omitempty"`
    RaceResultID         *int       `json:"race_result_id,omitempty"`
    SprintRaceResultID   *int       `json:"sprint_race_result_id,omitempty"`
    QualyResultID        *int       `json:"qualy_result_id,omitempty"`
    SprintQualyResultID  *int       `json:"sprint_qualy_result_id,omitempty"`
    FP1ID                *int       `json:"fp1_id,omitempty"`
    FP2ID                *int       `json:"fp2_id,omitempty"`
    FP3ID                *int       `json:"fp3_id,omitempty"`
}

// EventResponseDTO - DTO para devolver detalles completos de un evento
type EventResponseDTO struct {
    ID                   int         `json:"id"`
    Session              SessionDTO  `json:"session"`
    Date                 time.Time   `json:"date"`
    RaceResultID         *int        `json:"race_result_id,omitempty"`
    SprintRaceResultID   *int        `json:"sprint_race_result_id,omitempty"`
    QualyResultID        *int        `json:"qualy_result_id,omitempty"`
    SprintQualyResultID  *int        `json:"sprint_qualy_result_id,omitempty"`
    FP1ID                *int        `json:"fp1_id,omitempty"`
    FP2ID                *int        `json:"fp2_id,omitempty"`
    FP3ID                *int        `json:"fp3_id,omitempty"`
}

// ListEventDTO - DTO para listar eventos
type ListEventDTO struct {
    ID       int       `json:"id"`
    Session  string    `json:"session_name"`
    Date     time.Time `json:"date"`
}

type SessionDTO struct {
    ID            int    `json:"id"`
    SessionName   string `json:"session_name"`
    SessionType   string `json:"session_type"`
    CircuitShortName string `json:"circuit_short_name"`
    CountryCode   string `json:"country_code"`
    CountryName   string `json:"country_name"`
    Location      string `json:"location"`
    SessionKey    string `json:"session_key"`
}