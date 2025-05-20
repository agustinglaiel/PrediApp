package model

type Driver struct {
	ID            int    `gorm:"primaryKey;type:int" json:"id"`
	BroadcastName string `json:"broadcast_name" gorm:"type:varchar(100)"`
	CountryCode   string `json:"country_code" gorm:"type:varchar(10)"`
	DriverNumber  int    `json:"driver_number"`
	FirstName     string `json:"first_name" gorm:"type:varchar(50);index:idx_driver_name,priority:1"`
	LastName      string `json:"last_name" gorm:"type:varchar(50);index:idx_driver_name,priority:2"`
	FullName      string `json:"full_name" gorm:"type:varchar(100)"`
	NameAcronym   string `json:"name_acronym" gorm:"type:varchar(10)"`
	HeadshotURL   string `json:"headshot_url" gorm:"type:varchar(200)"`
	TeamName      string `json:"team_name" gorm:"type:varchar(100)"`
	Activo        bool   `json:"activo"`
}
