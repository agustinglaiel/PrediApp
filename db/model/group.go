package model

import "time"

// group.go
type Group struct {
	ID          int           `gorm:"primaryKey" json:"id"`
	GroupName   string        `json:"group_name"`
	Description string        `json:"description"`
	GroupCode   string        `gorm:"uniqueIndex" json:"group_code"`
	CreatedAt   time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   *time.Time    `gorm:"index" json:"deleted_at,omitempty"`
	GroupUsers  []GroupXUsers `gorm:"foreignKey:GroupID" json:"group_users"`
}

// group_x_users.go
type GroupXUsers struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	GroupID   int       `json:"group_id"`
	UserID    int       `json:"user_id"`
	GroupRole string    `json:"group_role"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
