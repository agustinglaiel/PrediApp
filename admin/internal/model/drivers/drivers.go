package drivers

import (
	"time"

	"gorm.io/gorm"
)

type Driver struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	BroadcastName string         `json:"broadcast_name"`
	CountryCode   string         `json:"country_code"`
	DriverNumber  int            `json:"driver_number"`
	FirstName     string         `json:"first_name"`
	FullName      string         `json:"full_name"`
	HeadshotURL   string         `json:"headshot_url"`
	LastName      string         `json:"last_name"`
	NameAcronym   string         `json:"name_acronym"`
	SessionKey    int            `json:"session_key"` // ForeignKey to Session
	TeamName      string         `json:"team_name"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
