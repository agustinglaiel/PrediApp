package model

type GroupXUser struct {
	ID      int    `json:"id"`
	UserId  int    `json:"user_id" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserId;references:ID"`
	User    User   `json:"user"`
	GroupId int    `json:"group_id" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:GroupId;references:ID"`
	Group   Group  `json:"group"`
}

type User struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	FirstName string `gorm:"size:255" json:"first_name"`
	LastName  string `gorm:"size:255" json:"last_name"`
	Username  string `gorm:"size:255;uniqueIndex" json:"username"`
	Email     string `gorm:"size:255;uniqueIndex" json:"email"`
	Password  string `gorm:"size:255" json:"-"` // omitir en la respuesta JSON
	Role      string `gorm:"size:255" json:"role"`
	Score     int    `gorm:"default:0" json:"score"`
}

type Group struct {
	ID          int    `gorm:"primaryKey" json:"id"`
	GroupName   string `gorm:"size:255" json:"group_name"`
	Description string `gorm:"size:255" json:"description"`
}