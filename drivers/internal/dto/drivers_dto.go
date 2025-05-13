package dto

type CreateDriverDTO struct {
    BroadcastName  string `json:"broadcast_name" binding:"required"`
    CountryCode    string `json:"country_code" binding:"required"`
    DriverNumber   int    `json:"driver_number" binding:"required"`
    FirstName      string `json:"first_name" binding:"required"`
    LastName       string `json:"last_name" binding:"required"`
    FullName       string `json:"full_name" binding:"required"`
    NameAcronym    string `json:"name_acronym" binding:"required"`
    HeadshotURL    string `json:"headshot_url,omitempty"`
    TeamName       string `json:"team_name" binding:"required"`
    Activo         bool   `json:"activo" binding:"required"`
}

type UpdateDriverDTO struct {
    BroadcastName  string `json:"broadcast_name,omitempty"`
    CountryCode    string `json:"country_code,omitempty"`
    DriverNumber   int    `json:"driver_number,omitempty"`
    FirstName      string `json:"first_name,omitempty"`
    LastName       string `json:"last_name,omitempty"`
    FullName       string `json:"full_name,omitempty"`
    NameAcronym    string `json:"name_acronym,omitempty"`
    HeadshotURL    string `json:"headshot_url,omitempty"`
    TeamName       string `json:"team_name,omitempty"`
    Activo         *bool   `json:"activo,omitempty"`
}

type ResponseDriverDTO struct {
    ID             int    `json:"id"`
    BroadcastName  string `json:"broadcast_name"`
    CountryCode    string `json:"country_code"`
    DriverNumber   int    `json:"driver_number"`
    FirstName      string `json:"first_name"`
    LastName       string `json:"last_name"`
    FullName       string `json:"full_name"`
    NameAcronym    string `json:"name_acronym"`
    HeadshotURL    string `json:"headshot_url"`
    TeamName       string `json:"team_name"`
    Activo         bool   `json:"activo"`
}

type DeleteDriverDTO struct {
    DriverID int `json:"driver_id"`
}
