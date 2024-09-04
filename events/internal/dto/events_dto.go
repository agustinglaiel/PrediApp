package dto

import "time"

// BasicResultDTO - DTO para manejar referencias básicas a los resultados
type BasicResultDTO struct {
    ID       int    `json:"id"`
    EventID  int    `json:"event_id"`
    Result   string `json:"result"`  // Este campo puede contener el resultado en un formato simple, como una lista de posiciones.
}

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
    ID                   int            `json:"id"`
    Session              SessionDTO     `json:"session"`
    Date                 time.Time      `json:"date"`
    RaceResult           *BasicResultDTO `json:"race_result,omitempty"`
    SprintRaceResult     *BasicResultDTO `json:"sprint_race_result,omitempty"`
    QualyResult          *BasicResultDTO `json:"qualy_result,omitempty"`
    SprintQualyResult    *BasicResultDTO `json:"sprint_qualy_result,omitempty"`
    FP1Result            *BasicResultDTO `json:"fp1_result,omitempty"`
    FP2Result            *BasicResultDTO `json:"fp2_result,omitempty"`
    FP3Result            *BasicResultDTO `json:"fp3_result,omitempty"`
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