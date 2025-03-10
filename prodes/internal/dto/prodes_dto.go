package prodes

import "time"

// DTO para crear un pronóstico de carrera
type CreateProdeCarreraDTO struct {
    UserID     int `json:"user_id"`
    SessionID  int `json:"session_id"` // Vinculado a la sesión
    P1         int `json:"p1"` // driver_id
    P2         int `json:"p2"` // driver_id
    P3         int `json:"p3"` // driver_id
    P4         int `json:"p4"` // driver_id
    P5         int `json:"p5"` // driver_id
    // FastestLap int `json:"fastest_lap"` // driver_id
    VSC        bool `json:"vsc"`
    SC         bool `json:"sc"`
    DNF        int  `json:"dnf"`
}

// DTO para crear un pronóstico de sesión que no sea carrera
type CreateProdeSessionDTO struct {
    UserID    int `json:"user_id"`
    SessionID int `json:"session_id"` // Vinculado a la sesión
    P1        int `json:"p1"` // driver_id
    P2        int `json:"p2"` // driver_id
    P3        int `json:"p3"` // driver_id
}

// DTO para eliminar un pronóstico
type DeleteProdeDTO struct {
    ProdeID int `json:"prode_id"`
    UserID  int `json:"user_id"`
}

// DTO para obtener pronósticos por usuario o evento
type GetProdesDTO struct {
    UserID  int `json:"user_id,omitempty"`
    EventID int `json:"event_id,omitempty"`
}

// DTO de respuesta para un pronóstico de carrera
type ResponseProdeCarreraDTO struct {
    ID         int  `json:"id"`
    UserID     int  `json:"user_id"`
    SessionID  int  `json:"session_id"` // Cambiado a session_id
    P1         int  `json:"p1"` // driver_id
    P2         int  `json:"p2"` // driver_id
    P3         int  `json:"p3"` // driver_id
    P4         int  `json:"p4"` // driver_id
    P5         int  `json:"p5"` // driver_id
    // FastestLap int  `json:"fastest_lap"` // driver_id
    VSC        bool `json:"vsc"`
    SC         bool `json:"sc"`
    DNF        int  `json:"dnf"`
    Score      int  `json:"score"`
}

// DTO de respuesta para un pronóstico de sesión
type ResponseProdeSessionDTO struct {
    ID        int `json:"id"`
    UserID    int `json:"user_id"`
    SessionID int `json:"session_id"` // Cambiado a session_id
    P1        int `json:"p1"` // driver_id
    P2        int `json:"p2"` // driver_id
    P3        int `json:"p3"` // driver_id
    Score     int  `json:"score"`
}

// DTO para actualizar un pronóstico de carrera
type UpdateProdeCarreraDTO struct {
    ProdeID    int `json:"prode_id"`
    UserID     int `json:"user_id"`
    SessionID  int `json:"session_id"` // Cambiado a session_id
    P1         int `json:"p1"` // driver_id
    P2         int `json:"p2"` // driver_id
    P3         int `json:"p3"` // driver_id
    P4         int `json:"p4"` // driver_id
    P5         int `json:"p5"` // driver_id
    // FastestLap int `json:"fastest_lap"` // driver_id
    VSC        bool `json:"vsc"`
    SC         bool `json:"sc"`
    DNF        int  `json:"dnf"`
}

// DTO para actualizar un pronóstico de sesión que no sea carrera normal
type UpdateProdeSessionDTO struct {
    ProdeID    int `json:"prode_id"`
    UserID     int `json:"user_id"`
    SessionID  int `json:"session_id"` // Cambiado a session_id
    P1         int `json:"p1"` // driver_id
    P2         int `json:"p2"` // driver_id
    P3         int `json:"p3"` // driver_id
}

type SessionNameAndTypeDTO struct {
	SessionName string `json:"session_name"`
	SessionType string `json:"session_type"`
}

/*
¿Por qué esta solución?
Independencia de microservicios: Mantienes a prodes separado de sessions sin depender directamente de sus DTOs. 
Cada microservicio tiene su propia lógica y estructura de datos.
Uso de DTOs equivalentes: El cliente HTTP en prodes hace una solicitud a sessions, 
recibe una respuesta en formato JSON, y usamos el DTO definido en prodes para deserializar esa respuesta. 
Esto asegura que los dos servicios sigan siendo independientes.
*/

// SessionDetailsDTO define los atributos que queremos obtener de la sesión
type SessionDetailsDTO struct {
    CircuitShortName string `json:"circuit_short_name"`
    CountryCode      string `json:"country_code"`
    CountryName      string `json:"country_name"`
    DateStart        time.Time `json:"date_start"`
    DateEnd          time.Time `json:"date_end"`
    Location         string `json:"location"`
    SessionName      string `json:"session_name"`
    SessionType      string `json:"session_type"`
}

type DriverDTO struct {
    ID             int    `json:"id"`              // Identificador único del piloto
    CountryCode    string `json:"country_code"`    // Código de país del piloto
    DriverNumber   int    `json:"driver_number"`   // Número del piloto
    FirstName      string `json:"first_name"`      // Nombre del piloto
    LastName       string `json:"last_name"`       // Apellido del piloto
    FullName       string `json:"full_name"`       // Nombre completo del piloto
    NameAcronym    string `json:"name_acronym"`    // Acrónimo del nombre
    TeamName       string `json:"team_name"`       // Nombre del equipo
}

type TopDriverDTO struct {
    Position int `json:"position"`
    DriverID int `json:"driver_id"`
}