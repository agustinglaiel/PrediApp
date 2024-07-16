package dto

// GoogleOAuthRequestDTO representa la solicitud de OAuth de Google
type GoogleOAuthRequestDTO struct {
    GoogleToken string `json:"google_token" binding:"required"`
}

// GoogleOAuthResponseDTO representa la respuesta de un registro/inicio de sesión exitoso utilizando Google OAuth
type GoogleOAuthResponseDTO struct {
    ID            string `json:"id"`
    FirstName     string `json:"first_name"`
    LastName      string `json:"last_name"`
    Username      string `json:"username"`
    Email         string `json:"email"`
    Role          string `json:"role"`
    Token         string `json:"token"` // JWT token
    Provider      string `json:"provider"`
    ProviderID    string `json:"provider_id"`
    AvatarURL     string `json:"avatar_url,omitempty"`
    AccessToken   string `json:"access_token,omitempty"`
    RefreshToken  string `json:"refresh_token,omitempty"`
    ExpiresAt     string `json:"expires_at,omitempty"`
}

// UserSignUpRequestDTO representa la solicitud de registro de un usuario con email y contraseña
type UserSignUpRequestDTO struct {
    FirstName   string `json:"first_name" binding:"required"`
    LastName    string `json:"last_name" binding:"required"`
    Username    string `json:"username" binding:"required"`
    Email       string `json:"email" binding:"required,email"`
    Password    string `json:"password" binding:"required"`
    PhoneNumber string `json:"phone_number,omitempty"`
}

// UserSignUpResponseDTO representa la respuesta de un registro exitoso
type UserSignUpResponseDTO struct {
    ID          string `json:"id"`
    FirstName   string `json:"first_name"`
    LastName    string `json:"last_name"`
    Username    string `json:"username"`
    Email       string `json:"email"`
    Role        string `json:"role"`
    CreatedAt   string `json:"created_at"`
}

// UserLoginRequestDTO representa la solicitud de inicio de sesión de un usuario con email y contraseña
type UserLoginRequestDTO struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

// UserLoginResponseDTO representa la respuesta de un inicio de sesión exitoso
type UserLoginResponseDTO struct {
    ID        string `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Username  string `json:"username"`
    Email     string `json:"email"`
    Role      string `json:"role"`
    Token     string `json:"token"` // JWT token
}