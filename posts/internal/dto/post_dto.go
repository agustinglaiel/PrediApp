package dto

type PostCreateRequestDTO struct {
	UserID	     int    `json:"user_id" binding:"required"` // ID del usuario que crea el post
	ParentPostID *int   `json:"parent_post_id,omitempty"` // NULL para posts principales, ID para comentarios
	Body         string `json:"body" binding:"required,max=500"`
}

type PostResponseDTO struct {
	ID           int         `json:"id"`
	UserID       int         `json:"user_id"`
	ParentPostID *int        `json:"parent_post_id"`
	Body         string      `json:"body"`
	CreatedAt    string      `json:"created_at"`
	Children     []PostResponseDTO `json:"children,omitempty"`
}