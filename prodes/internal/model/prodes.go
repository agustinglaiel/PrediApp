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
	Score      int            `gorm:"default:0" json:"score"`
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
	Score      int            `gorm:"default:0" json:"score"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}


// Session es un modelo simplificado para representar la sesión que se relaciona con los prodes
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

type User struct {
	ID              int           `gorm:"primaryKey" json:"id"`
	FirstName       string         `gorm:"size:255" json:"first_name"`
	LastName        string         `gorm:"size:255" json:"last_name"`
	Username        string         `gorm:"size:255;uniqueIndex" json:"username"`
	Email           string         `gorm:"size:255;uniqueIndex" json:"email"`
	Password        string         `gorm:"size:255" json:"-"` // omitir en la respuesta JSON
	Role            string         `gorm:"size:255" json:"role"`
	Score           int            `gorm:"default:0" json:"score"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	IsEmailVerified bool           `gorm:"default:false" json:"is_email_verified"`
	LastLoginAt     *time.Time     `json:"last_login_at,omitempty"`
	PhoneNumber     string         `gorm:"size:20" json:"phone_number,omitempty"`
	Provider        string         `gorm:"size:255" json:"provider,omitempty"`
	ProviderID      string         `gorm:"size:255" json:"provider_id,omitempty"`
	AvatarURL       string         `gorm:"size:255" json:"avatar_url,omitempty"`
	RefreshTokens   []RefreshToken `gorm:"foreignKey:UserID" json:"-"`
	Posts           []*Post        `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
}

type Driver struct {
	ID            int    `gorm:"primaryKey;type:int" json:"id"`
	BroadcastName string `json:"broadcast_name" gorm:"type:varchar(100)"`
	CountryCode   string `json:"country_code" gorm:"type:varchar(10)"`
	DriverNumber  int    `json:"driver_number"`
	FirstName     string `json:"first_name" gorm:"type:varchar(50);index:idx_driver_name,priority:1"`
    LastName      string `json:"last_name" gorm:"type:varchar(50);index:idx_driver_name,priority:2"`
	FullName      string `json:"full_name" gorm:"type:varchar(100)"`
	NameAcronym   string `json:"name_acronym" gorm:"type:varchar(10)"`
	HeadshotURL   string `json:"headshot_url" gorm:"type:varchar(200)"` // Añadimos el campo de la foto
	TeamName      string `json:"team_name" gorm:"type:varchar(100)"`
	Activo        bool   `json:"activo"`
}

type RefreshToken struct {
    ID        int       `gorm:"primaryKey" json:"id"`
    UserID    int       `gorm:"index;foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user_id"`
    Token     string    `gorm:"size:255;uniqueIndex" json:"token"`
    ExpiresAt time.Time `json:"expires_at"`
}

type Post struct {
	ID           int          `gorm:"primaryKey"`
	UserID       int          `gorm:"index;not null"`
	ParentPostID *int         `gorm:"index"`
	Body         string       `gorm:"type:text;not null"`
	CreatedAt    time.Time    `gorm:"autoCreateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}