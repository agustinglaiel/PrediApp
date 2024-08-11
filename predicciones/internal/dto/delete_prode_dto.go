package dto

// DTO para eliminar un pronóstico
type DeleteProdeDTO struct {
    ProdeID int `json:"prode_id"`
    UserID  int `json:"user_id"`
}
