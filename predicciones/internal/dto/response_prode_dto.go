package dto

// DTO de respuesta para un pronóstico de carrera
type ResponseProdeCarreraDTO struct {
    ID         int `json:"id"`
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

// DTO de respuesta para un pronóstico de sesión
type ResponseProdeSessionDTO struct {
    ID         int `json:"id"`
    UserID     int `json:"user_id"`
    EventID    int `json:"event_id"`
    P1         int `json:"p1"` // driver_id
    P2         int `json:"p2"` // driver_id
    P3         int `json:"p3"` // driver_id
}


// son utilizados para estructurar y enviar la información de vuelta al cliente después de realizar 
// una operación en el servidor. Estos DTOs aseguran que solo se devuelvan los datos necesarios 
// en el formato adecuado, mejorando la seguridad y eficiencia de la comunicación entre el backend y el frontend.