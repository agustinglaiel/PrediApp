package dto

// DTO para actualizar un pronóstico de carrera
type UpdateProdeCarreraDTO struct {
    ProdeID    int `json:"prode_id"`
    UserID     int `json:"user_id"`
    EventID    int `json:"event_id"`
    P1         int `json:"p1"` // driver_id
    P2         int `json:"p2"` // driver_id
    P3         int `json:"p3"` // driver_id
    P4         int `json:"p4"` // driver_id
    P5         int `json:"p5"` // driver_id
    FastestLap int `json:"fastest_lap"` // driver_id
    VSC        bool `json:"vsc"`
    SC         bool `json:"sc"`
    DNF        int  `json:"dnf"`
}

// DTO para actualizar un pronóstico de sesión que no sea carrera normal
type UpdateProdeSessionDTO struct {
    ProdeID    int `json:"prode_id"`
    UserID     int `json:"user_id"`
    EventID    int `json:"event_id"`
    P1         int `json:"p1"` // driver_id
    P2         int `json:"p2"` // driver_id
    P3         int `json:"p3"` // driver_id
}
