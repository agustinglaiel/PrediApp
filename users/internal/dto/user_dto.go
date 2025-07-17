package dto

// UserSignUpRequestDTO representa la solicitud de registro de un usuario con email y contraseña
type UserSignUpRequestDTO struct {
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required"`
	Role        string `json:"role"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

// UserSignUpResponseDTO representa la respuesta de un registro exitoso
type UserSignUpResponseDTO struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Score     int    `json:"score"`
	Token     string `json:"token"`
	CreatedAt string `json:"created_at"`
}

// UserLoginRequestDTO representa la solicitud de inicio de sesión de un usuario con email y contraseña
type UserLoginRequestDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserLoginResponseDTO representa la respuesta de un inicio de sesión exitoso
type UserLoginResponseDTO struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Score     int    `json:"score"`
	Token     string `json:"token"`
}

// UserResponseDTO representa la respuesta general de un usuario
type UserResponseDTO struct {
	ID             int    `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	Score          int    `json:"score"`
	PhoneNumber    string `json:"phone_number,omitempty"`
	CreatedAt      string `json:"created_at"`
	IsActive       bool   `json:"is_active"`
	ImagenPerfil   string `json:"imagen_perfil,omitempty"`
	ImagenMimeType string `json:"imagen_mime_type,omitempty"`
}

// UserUpdateRequestDTO representa la solicitud de actualización de un usuario
type UserUpdateRequestDTO struct {
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

type UserUpdateRoleRequestDTO struct {
	Role string `json:"role" binding:"required"`
}

type UserScoreDtoSimplified struct {
	Username string `json:"username"`
	Score    int    `json:"score"`
}
