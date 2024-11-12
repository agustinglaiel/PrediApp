package dto

type GroupXUserCreateDTO struct {
	UserId  int `json:"user_id" binding:"required"`
	GroupId int `json:"group_id" binding:"required"`
}

type GroupXUserResponseDTO struct {
	ID      int `json:"id"`
	UserId  int `json:"user_id"`
	GroupId int `json:"group_id"`
}

type GroupXUserResponseWithUserAndGroupDTO struct {
	ID      int    `json:"id"`
	UserId  int    `json:"user_id"`
	User    User   `json:"user"`
	GroupId int    `json:"group_id"`
	Group   Group  `json:"group"`
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type Group struct {
	ID          int    `json:"id"`
	GroupName   string `json:"group_name"`
}