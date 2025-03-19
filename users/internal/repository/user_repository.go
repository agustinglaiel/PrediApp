package repository

import (
	"context"
	"log"
	"users/internal/model"
	e "users/pkg/utils"

	"gorm.io/gorm"
)

// userRepository es una estructura vacía que implementa la interfaz UserRepository
type userRepository struct {
    db *gorm.DB
}

// UserRepository define los métodos que deben ser implementados por el repositorio de usuarios
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) e.ApiError
	GetUserByEmail(ctx context.Context, email string) (*model.User, e.ApiError)
	GetUserByUsername(ctx context.Context, username string) (*model.User, e.ApiError)
	GetUserByID(ctx context.Context, id int) (*model.User, e.ApiError)
	GetUsers(ctx context.Context) ([]*model.User, e.ApiError)
	UpdateUserByID(ctx context.Context, id int, user *model.User) e.ApiError
	UpdateUserByUsername(ctx context.Context, username string, user *model.User) e.ApiError
	DeleteUserByID(ctx context.Context, id int) e.ApiError
	DeleteUserByUsername(ctx context.Context, username string) e.ApiError

	// Métodos para el Refresh Token
	CreateRefreshToken(ctx context.Context, refreshToken *model.RefreshToken) e.ApiError
	GetRefreshToken(ctx context.Context, token string) (*model.RefreshToken, e.ApiError)
	UpdateRefreshToken(ctx context.Context, refreshToken *model.RefreshToken) e.ApiError
	DeleteRefreshToken(ctx context.Context, token string) e.ApiError 
}

// NewUserRepository crea una nueva instancia de userRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// CreateUser inserta un nuevo usuario en la base de datos
func (r *userRepository) CreateUser(ctx context.Context, user *model.User) e.ApiError {
	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		return e.NewInternalServerApiError("error creating user", err)
	}
	return nil
}

// GetUserByEmail obtiene un usuario por su email
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, e.ApiError) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("User not found: %s", email)
			return nil, e.NewNotFoundApiError("user not found")
		}
		log.Printf("Error finding user by email: %v", err)
		return nil, e.NewInternalServerApiError("error finding user by email", err)
	}
	return &user, nil
}

// GetUserByUsername obtiene un usuario por su nombre de usuario
func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, e.ApiError) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, e.NewNotFoundApiError("user not found")
		}
		return nil, e.NewInternalServerApiError("error finding user by username", err)
	}
	return &user, nil
}

// GetUserByID obtiene un usuario por su ID
func (r *userRepository) GetUserByID(ctx context.Context, id int) (*model.User, e.ApiError) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, e.NewNotFoundApiError("user not found")
		}
		return nil, e.NewInternalServerApiError("error finding user by ID", err)
	}
	return &user, nil
}

// GetUsers obtiene todos los usuarios
func (r *userRepository) GetUsers(ctx context.Context) ([]*model.User, e.ApiError) {
	var users []*model.User
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, e.NewInternalServerApiError("error finding users", err)
	}
	return users, nil
}

// UpdateUserByID actualiza un usuario por su ID en la base de datos
func (r *userRepository) UpdateUserByID(ctx context.Context, id int, user *model.User) e.ApiError {
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(user).Error; err != nil {
		return e.NewInternalServerApiError("error updating user by ID", err)
	}
	return nil
}

// UpdateUserByUsername actualiza un usuario por su nombre de usuario en la base de datos
func (r *userRepository) UpdateUserByUsername(ctx context.Context, username string, user *model.User) e.ApiError {
	if err := r.db.WithContext(ctx).Model(&model.User{}).Where("username = ?", username).Updates(user).Error; err != nil {
		return e.NewInternalServerApiError("error updating user by username", err)
	}
	return nil
}

// DeleteUserByID elimina un usuario por su ID de la base de datos
func (r *userRepository) DeleteUserByID(ctx context.Context, id int) e.ApiError {
    var user model.User
    if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return e.NewNotFoundApiError("user not found")
        }
        return e.NewInternalServerApiError("error finding user by ID", err)
    }
    if err := r.db.WithContext(ctx).Unscoped().Delete(&user).Error; err != nil { // Unscoped para eliminación física
        return e.NewInternalServerApiError("error deleting user by ID", err)
    }
    return nil
}

// DeleteUserByUsername elimina un usuario por su nombre de usuario de la base de datos
func (r *userRepository) DeleteUserByUsername(ctx context.Context, username string) e.ApiError {
	if err := r.db.WithContext(ctx).Where("username = ?", username).Delete(&model.User{}).Error; err != nil {
		return e.NewInternalServerApiError("error deleting user by username", err)
	}
	return nil
}

// Metodos para el Refresh Token
func (r *userRepository) CreateRefreshToken(ctx context.Context, refreshToken *model.RefreshToken) e.ApiError {
	if err := r.db.WithContext(ctx).Create(refreshToken).Error; err != nil {
		log.Printf("Error creating refresh token: %v", err)
		return e.NewInternalServerApiError("error creating refresh token", err)
	}
	return nil
}

// GetRefreshToken obtiene un refresh token por su valor
func (r *userRepository) GetRefreshToken(ctx context.Context, token string) (*model.RefreshToken, e.ApiError) {
	var refreshToken model.RefreshToken
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&refreshToken).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, e.NewNotFoundApiError("refresh token not found")
		}
		log.Printf("Error finding refresh token: %v", err)
		return nil, e.NewInternalServerApiError("error finding refresh token", err)
	}
	return &refreshToken, nil
}

// UpdateRefreshToken actualiza un refresh token en la base de datos
func (r *userRepository) UpdateRefreshToken(ctx context.Context, refreshToken *model.RefreshToken) e.ApiError {
	if err := r.db.WithContext(ctx).Save(refreshToken).Error; err != nil {
		log.Printf("Error updating refresh token: %v", err)
		return e.NewInternalServerApiError("error updating refresh token", err)
	}
	return nil
}

// DeleteRefreshToken elimina un refresh token de la base de datos
func (r *userRepository) DeleteRefreshToken(ctx context.Context, token string) e.ApiError {
    if err := r.db.WithContext(ctx).Where("token = ?", token).Delete(&model.RefreshToken{}).Error; err != nil {
        log.Printf("Error deleting refresh token: %v", err)
        return e.NewInternalServerApiError("error deleting refresh token", err)
    }
    return nil
}