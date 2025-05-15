package model

type Driver struct {
	ID            int    `gorm:"primaryKey;type:int" json:"id"`
	BroadcastName string `json:"broadcast_name"`
	CountryCode   string `json:"country_code"`
	DriverNumber  int    `json:"driver_number"`
	FirstName     string `json:"first_name" gorm:"index:idx_driver_name"`
	LastName      string `json:"last_name" gorm:"index:idx_driver_name"`
	FullName      string `json:"full_name"`
	NameAcronym   string `json:"name_acronym"`
	HeadshotURL   string `json:"headshot_url"` // Añadimos el campo de la foto
	TeamName      string `json:"team_name"`
	Activo        bool   `json:"activo"`
}
