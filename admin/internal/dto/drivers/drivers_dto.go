package drivers

type CreateDriverDTO struct {
    FirstName      string `json:"first_name" binding:"required"`
    LastName       string `json:"last_name" binding:"required"`
    FullName       string `json:"full_name" binding:"required"`
    CountryCode    string `json:"country_code" binding:"required"`
    DriverNumber   int    `json:"driver_number" binding:"required"`
    NameAcronym    string `json:"name_acronym" binding:"required"`
    TeamName       string `json:"team_name" binding:"required"`
    HeadshotURL    string `json:"headshot_url,omitempty"`
    SessionKey     int    `json:"session_key,omitempty"`
}


type UpdateDriverDTO struct {
    FirstName      string `json:"first_name,omitempty"`
    LastName       string `json:"last_name,omitempty"`
    FullName       string `json:"full_name,omitempty"`
    CountryCode    string `json:"country_code,omitempty"`
    DriverNumber   int    `json:"driver_number,omitempty"`
    NameAcronym    string `json:"name_acronym,omitempty"`
    TeamName       string `json:"team_name,omitempty"`
    HeadshotURL    string `json:"headshot_url,omitempty"`
    SessionKey     int    `json:"session_key,omitempty"`
}

type ResponseDriverDTO struct {
    ID             uint   `json:"id"`
    FirstName      string `json:"first_name"`
    LastName       string `json:"last_name"`
    FullName       string `json:"full_name"`
    CountryCode    string `json:"country_code"`
    DriverNumber   int    `json:"driver_number"`
    NameAcronym    string `json:"name_acronym"`
    TeamName       string `json:"team_name"`
    HeadshotURL    string `json:"headshot_url,omitempty"`
    SessionKey     int    `json:"session_key,omitempty"`
}

type DeleteDriverDTO struct {
    DriverID uint `json:"driver_id"`
}