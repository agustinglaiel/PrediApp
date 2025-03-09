package model

import (
	"time"

	"gorm.io/gorm"
)

type ProdeCarrera struct {
	ID         int            `gorm:"primaryKey" json:"id"`
	UserID     int            `gorm:"index;not null" json:"user_id"`
	User       User           `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SessionID  int            `gorm:"index;not null" json:"session_id"`
	Session    Session        `gorm:"foreignKey:SessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	P1         int            `json:"p1"` // driver_id
	DriverP1   Driver         `gorm:"foreignKey:P1;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	P2         int            `json:"p2"` // driver_id
	DriverP2   Driver         `gorm:"foreignKey:P2;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	P3         int            `json:"p3"` // driver_id
	DriverP3   Driver         `gorm:"foreignKey:P3;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	P4         int            `json:"p4"` // driver_id
	DriverP4   Driver         `gorm:"foreignKey:P4;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	P5         int            `json:"p5"` // driver_id
	DriverP5   Driver         `gorm:"foreignKey:P5;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	// FastestLap int            `json:"fastest_lap"` // driver_id
	// DriverFastestLap Driver   `gorm:"foreignKey:FastestLap;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	VSC        bool           `json:"vsc"`
	SC         bool           `json:"sc"`
	DNF        int            `json:"dnf"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}


type ProdeSession struct {
	ID        int            `gorm:"primaryKey" json:"id"`
	UserID    int            `gorm:"index;not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SessionID int            `gorm:"index;not null" json:"session_id"`
	Session   Session        `gorm:"foreignKey:SessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	P1        int            `json:"p1"` // driver_id
	DriverP1  Driver         `gorm:"foreignKey:P1;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	P2        int            `json:"p2"` // driver_id
	DriverP2  Driver         `gorm:"foreignKey:P2;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	P3        int            `json:"p3"` // driver_id
	DriverP3  Driver         `gorm:"foreignKey:P3;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}


// Session es un modelo simplificado para representar la sesión que se relaciona con los prodes
type Session struct {
	ID               int       `gorm:"primaryKey" json:"id"`
	CircuitShortName string    `json:"circuit_short_name"`
	CountryCode      string    `json:"country_code"`
	CountryName      string    `json:"country_name"`
	DateStart        time.Time `json:"date_start"`
	DateEnd          time.Time `json:"date_end"`
	Location         string    `json:"location"`
}

type User struct {
	ID              int         `gorm:"primaryKey" json:"id"`
	FirstName       string      `gorm:"size:255" json:"first_name"`
	LastName		string      `gorm:"size:255" json:"last_name"`
	Username        string      `gorm:"size:255;uniqueIndex" json:"username"`
	Email           string      `gorm:"size:255;uniqueIndex" json:"email"`
	Password        string      `gorm:"size:255" json:"-"` // omitir en la respuesta JSON
	Role            string      `gorm:"size:255" json:"role"`
	Score           int         `gorm:"default:0" json:"score"`
}

type Driver struct {
    ID          int    `json:"id"`
    FirstName   string `json:"first_name"`
    LastName    string `json:"last_name"`
    FullName    string `json:"full_name"`
    NameAcronym string `json:"name_acronym"`
    TeamName    string `json:"team_name"`
}