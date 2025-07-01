package model

import "time"

// group.go
type Group struct {
	ID          int           `gorm:"primaryKey" json:"id"`
	GroupName   string        `gorm:"size:255;not null" json:"group_name"`
	Description string        `gorm:"size:255" json:"description"`
	GroupCode   string        `gorm:"size:8;uniqueIndex;not null" json:"group_code"`
	GroupUsers  []GroupXUsers `gorm:"foreignKey:GroupID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"group_users"`
	CreatedAt   time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   *time.Time    `gorm:"index" json:"deleted_at,omitempty"`
}

// GroupXUsers representa la relación entre grupos y usuarios
type GroupXUsers struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	GroupID   int       `gorm:"index;not null" json:"group_id"`
	UserID    int       `gorm:"index;not null" json:"user_id"`
	GroupRole string    `gorm:"size:50;not null" json:"group_role"`
	User      *User     `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"` // Eliminamos la carga circular
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
