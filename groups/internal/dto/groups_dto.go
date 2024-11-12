package dto

type CreateGroupRequestDTO struct {
	GroupName string `json:"group_name" binding:"required"`
	Description string `json:"description"`
}

type GroupResponseDTO struct {
	ID        int    `json:"id"`
	GroupName string `json:"group_name"`
	Description string `json:"description"`
}