package dto

// DTO de respuesta para un pronóstico de carrera
type ResponseProdeCarreraDTO struct {
    ID         uint `json:"id"`
    UserID     uint `json:"user_id"`
    EventID    uint `json:"event_id"`
    P1         uint `json:"p1"` // driver_id
    P2         uint `json:"p2"` // driver_id
    P3         uint `json:"p3"` // driver_id
    P4         uint `json:"p4"` // driver_id
    P5         uint `json:"p5"` // driver_id
    FastestLap uint `json:"fastest_lap"` // driver_id
    VSC        bool `json:"vsc"`
    SC         bool `json:"sc"`
    DNF        int  `json:"dnf"`
}

// DTO de respuesta para un pronóstico de sesión
type ResponseProdeSessionDTO struct {
    ID         uint `json:"id"`
    UserID     uint `json:"user_id"`
    EventID    uint `json:"event_id"`
    P1         uint `json:"p1"` // driver_id
    P2         uint `json:"p2"` // driver_id
    P3         uint `json:"p3"` // driver_id
}


// son utilizados para estructurar y enviar la información de vuelta al cliente después de realizar 
// una operación en el servidor. Estos DTOs aseguran que solo se devuelvan los datos necesarios 
// en el formato adecuado, mejorando la seguridad y eficiencia de la comunicación entre el backend y el frontend.