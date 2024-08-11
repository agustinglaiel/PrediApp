package dto

// DTO para obtener pronósticos por usuario o evento
type GetProdesDTO struct {
    UserID  int `json:"user_id,omitempty"`
    EventID int `json:"event_id,omitempty"`
}
