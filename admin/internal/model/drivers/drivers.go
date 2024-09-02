package model

type Driver struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	BroadcastName  string `json:"broadcast_name"`
	CountryCode    string `json:"country_code"`
	DriverNumber   int    `json:"driver_number"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	FullName       string `json:"full_name"`
	NameAcronym    string `json:"name_acronym"`
	TeamName       string `json:"team_name"`
}

type DriverEvent struct {
	ID       uint `gorm:"primaryKey" json:"id"`
	EventID  uint `json:"event_id"`
	DriverID uint `json:"driver_id"`
}