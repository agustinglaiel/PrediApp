package dto

import "time"

// CreateResultDTO - DTO para la creación de un nuevo resultado de evento
type CreateResultDTO struct {
    EventID         int       `json:"event_id"`          // ID del evento asociado
    Date            time.Time `json:"date"`              // Fecha del evento
    DriverPositions []int     `json:"driver_positions"`  // IDs de los pilotos en cada posición, indexado de 0 a 19
    FastestLap      int       `json:"fastest_lap"`       // ID del piloto con la vuelta más rápida
    VSC             bool      `json:"vsc"`               // Indica si hubo coche de seguridad virtual
    SC              bool      `json:"sc"`                // Indica si hubo coche de seguridad
    DNF             int       `json:"dnf"`               // Número de pilotos que no terminaron
}

// ResponseResultDTO - DTO para devolver detalles de un resultado de evento
type ResponseResultDTO struct {
    ID             int            `json:"id"`              // ID del resultado
    EventID        int            `json:"event_id"`        // ID del evento asociado
    Date           time.Time      `json:"date"`            // Fecha del evento
    DriverResults  []int          `json:"driver_results"`  // ID de los pilotos, ordenados por posición
    FastestLap     int            `json:"fastest_lap"`     // ID del piloto con la vuelta más rápida
    VSC            bool           `json:"vsc"`             // Indica si hubo coche de seguridad virtual
    SC             bool           `json:"sc"`              // Indica si hubo coche de seguridad
    DNF            int            `json:"dnf"`             // Número de pilotos que no terminaron
}

// UpdateResultDTO - DTO para la actualización de un resultado de evento
type UpdateResultDTO struct {
    ID             int            `json:"id"`              // ID del resultado
    EventID        int            `json:"event_id"`        // ID del evento para validación
    Date           time.Time      `json:"date,omitempty"`  // Fecha del evento, opcional
    DriverResults  []int          `json:"driver_results,omitempty"`  // ID de los pilotos actualizados, ordenados por posición
    FastestLap     *int            `json:"fastest_lap,omitempty"`     // ID del piloto con la vuelta más rápida actualizado
    VSC            *bool          `json:"vsc,omitempty"`             // Actualizar si hubo coche de seguridad virtual
    SC             *bool          `json:"sc,omitempty"`              // Actualizar si hubo coche de seguridad
    DNF            *int           `json:"dnf,omitempty"`             // Actualizar el número de pilotos que no terminaron
}

// DeleteResultDTO - DTO para la solicitud de eliminación de un resultado
type DeleteResultDTO struct {
    ID       int `json:"id"`       // ID del resultado
    EventID  int `json:"event_id"` // ID del evento asociado
}