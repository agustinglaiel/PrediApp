package model

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID           int            `gorm:"primaryKey" json:"id"`
	UserID       int            `gorm:"index;foreignKey:UserID;references:ID" json:"user_id"`
	User         *User          `gorm:"foreignKey:UserID;references:ID" json:"-"`
	ParentPostID *int           `gorm:"index;foreignKey:ParentPostID;references:ID" json:"parent_post_id"` // NULL si es un post principal
	Body         string         `gorm:"type:varchar(500);not null" json:"body"`
	CreatedAt    time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Children     []*Post        `gorm:"-" json:"children"` // Para manejar hilos, ignorado por GORM
}

// El campo Children es un campo virtual (no se almacena en la base de datos,
// por eso tiene la etiqueta gorm:"-") que se usa para construir la estructura
// de hilos (es decir, un árbol de posts y comentarios) al devolver datos
// al cliente.

type User struct {
	ID       int    `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique;not null" json:"username"`
}
