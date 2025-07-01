package model

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	ID               int            `gorm:"primaryKey" json:"id"`
	WeekendID        int            `json:"weekend_id"`
	CircuitKey       int            `json:"circuit_key"`
	CircuitShortName string         `json:"circuit_short_name"`
	CountryCode      string         `json:"country_code"`
	CountryKey       int            `json:"country_key"`
	CountryName      string         `json:"country_name"`
	Location         string         `json:"location"`
	SessionKey       *int           `json:"session_key"`
	SessionName      string         `json:"session_name"`
	SessionType      string         `json:"session_type"`
	DateStart        time.Time      `json:"date_start" gorm:"type:timestamp"`
	DateEnd          time.Time      `json:"date_end" gorm:"type:timestamp"`
	Year             int            `json:"year"`
	VSC              *bool          `json:"vsc,omitempty"`
	SF               *bool          `json:"sf,omitempty"`
	DNF              *int           `json:"dnf,omitempty"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
