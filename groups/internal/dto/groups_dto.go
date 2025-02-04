package dto

type CreateGroupRequestDTO struct {
	GroupName   string `json:"group_name" binding:"required"`
	Description string `json:"description"`
	UserID      int    `json:"user_id" binding:"required"` // ID del usuario que crea el grupo
}

type GroupResponseDTO struct {
	ID          int                     `json:"id"`
	GroupName   string                  `json:"group_name"`
	Description string                  `json:"description"`
	GroupCode   string                  `json:"group_code"` // Código único del grupo
	Users       []GroupUserResponseDTO  `json:"users"`      // Usuarios en el grupo y sus roles
	CreatedAt   string                  `json:"created_at"` // Fecha de creación del grupo
	UpdatedAt   string                  `json:"updated_at"` // Fecha de última actualización
}

type GroupUserResponseDTO struct {
	UserID int    `json:"user_id"` // ID del usuario
	Role   string `json:"role"`    // Rol del usuario: "creator" o "invited"
	Score  *int	  `json:"score,omitempty"`   // Puntuación del usuario en el grupo
}

type GroupXUsersRequestDTO struct {
	GroupID int `json:"group_id" binding:"required"`
	UserID  int `json:"user_id" binding:"required"`
	Role    string `json:"role" binding:"required"` // Rol en el grupo
}

type GroupListResponseDTO struct {
	ID        int    `json:"id"`
	GroupName string `json:"group_name"`
	GroupCode string `json:"group_code"` // Código único del grupo
}

type RequestJoinGroupDTO struct {
	GroupCode string `json:"group_code" binding:"required"` // Código del grupo al que quiere unirse
	UserID    int    `json:"user_id" binding:"required"`    // ID del usuario que solicita unirse
}

type ManageGroupInvitationDTO struct {
	GroupID  int    `json:"group_id" binding:"required"`  // ID del grupo
	UserID   int    `json:"user_id" binding:"required"`   // ID del usuario a aceptar/rechazar
	Action   string `json:"action" binding:"required"`    // Acción: "accept" o "reject"
}
