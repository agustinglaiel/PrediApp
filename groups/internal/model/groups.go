package model

import "time"

type Group struct {
	ID		  	int       	`gorm:"primaryKey" json:"id"`	
	GroupName 	string		`gorm:"size:255" json:"group_name"`
	Description string		`gorm:"size:255" json:"description"`
	GroupXUsers GroupXUsers `gorm:"foreignKey:GroupID" json:"group_x_users"`
	CreatedAt 	time.Time 	`gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt 	time.Time 	`gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt 	*time.Time 	`gorm:"autoDeleteTime" json:"deleted_at,omitempty"`
}

type GroupXUsers struct {
	ID			int		`gorm:"primaryKey" json:"id"`
	GroupID		int		`json:"group_id"`
	UserID		int		`json:"user_id"`
}