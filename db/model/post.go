package model

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID           int            `gorm:"primaryKey" json:"id"`
	UserID       int            `gorm:"index" json:"user_id"`
	User         *User          `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE" json:"-"`
	ParentPostID *int           `gorm:"index;foreignKey:ParentPostID;references:ID" json:"parent_post_id"`
	Body         string         `gorm:"type:varchar(500);not null" json:"body"`
	CreatedAt    time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Children     []*Post        `gorm:"-" json:"children"`
}
