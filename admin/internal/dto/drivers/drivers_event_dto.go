package dto

type DriverEventDTO struct {
    EventID  uint `json:"event_id" binding:"required"`
    DriverID uint `json:"driver_id" binding:"required"`
}

type ResponseDriverEventDTO struct {
    ID       uint `json:"id"`
    EventID  uint `json:"event_id"`
    DriverID uint `json:"driver_id"`
}