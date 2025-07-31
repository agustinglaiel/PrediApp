package service

import (
	"context"
	"encoding/base64"
	"time"

	"prediapp.local/db/model"
	"prediapp.local/users/internal/dto"
	"prediapp.local/users/internal/repository"
	e "prediapp.local/users/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo repository.UserRepository
}

type UserServiceInterface interface {
	SignUp(ctx context.Context, request dto.UserSignUpRequestDTO) (dto.UserSignUpResponseDTO, e.ApiError)
	Login(ctx context.Context, request dto.UserLoginRequestDTO) (dto.UserLoginResponseDTO, e.ApiError)
	GetUserById(ctx context.Context, id int) (dto.UserResponseDTO, e.ApiError)
	GetUsers(ctx context.Context) ([]dto.UserResponseDTO, e.ApiError)
	UpdateUserById(ctx context.Context, id int, request dto.UserUpdateRequestDTO) (dto.UserResponseDTO, e.ApiError)
	DeleteUserById(ctx context.Context, id int) e.ApiError
	UpdateRoleByUserId(ctx context.Context, id int, request dto.UserUpdateRoleRequestDTO) (dto.UserResponseDTO, e.ApiError)
	GetUserScoreByUserId(ctx context.Context, id int) (dto.UserScoreDtoSimplified, e.ApiError)
	UploadProfilePicture(ctx context.Context, id int, image []byte, mimeType string) e.ApiError
	GetScoreboard(ctx context.Context) ([]dto.UserScoreDtoSimplified, e.ApiError)
}

func NewUserService(userRepo repository.UserRepository) UserServiceInterface {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) SignUp(ctx context.Context, request dto.UserSignUpRequestDTO) (dto.UserSignUpResponseDTO, e.ApiError) {
	if _, err := s.userRepo.GetUserByEmail(ctx, request.Email); err == nil {
		return dto.UserSignUpResponseDTO{}, e.NewBadRequestApiError("email already exists")
	}

	if _, err := s.userRepo.GetUserByUsername(ctx, request.Username); err == nil {
		return dto.UserSignUpResponseDTO{}, e.NewBadRequestApiError("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.UserSignUpResponseDTO{}, e.NewInternalServerApiError("error hashing password", err)
	}

	if request.Role == "" {
		request.Role = "user"
	}

	newUser := &model.User{
		FirstName:       request.FirstName,
		LastName:        request.LastName,
		Username:        request.Username,
		Email:           request.Email,
		Password:        string(hashedPassword),
		Role:            request.Role,
		Score:           0,
		CreatedAt:       time.Now(),
		IsActive:        true,
		IsEmailVerified: false,
	}

	if err := s.userRepo.CreateUser(ctx, newUser); err != nil {
		return dto.UserSignUpResponseDTO{}, e.NewInternalServerApiError("error creating user", err)
	}

	response := dto.UserSignUpResponseDTO{
		ID:        newUser.ID,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Username:  newUser.Username,
		Email:     newUser.Email,
		Role:      newUser.Role,
		Score:     newUser.Score,
		// Token:     token,
		CreatedAt: newUser.CreatedAt.Format(time.RFC3339),
	}

	return response, nil
}

func (s *userService) Login(ctx context.Context, request dto.UserLoginRequestDTO) (dto.UserLoginResponseDTO, e.ApiError) {
	user, err := s.userRepo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return dto.UserLoginResponseDTO{}, e.NewBadRequestApiError("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return dto.UserLoginResponseDTO{}, e.NewBadRequestApiError("invalid credentials")
	}

	response := dto.UserLoginResponseDTO{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Score:     user.Score,
	}

	return response, nil
}

func (s *userService) GetUserById(ctx context.Context, id int) (dto.UserResponseDTO, e.ApiError) {
	user, apiErr := s.userRepo.GetUserByID(ctx, id)
	if apiErr != nil {
		return dto.UserResponseDTO{}, apiErr
	}

	// codificar imagen a base64 si existe
	var imgB64 string
	if len(user.ImagenPerfil) > 0 {
		imgB64 = base64.StdEncoding.EncodeToString(user.ImagenPerfil)
	}

	response := dto.UserResponseDTO{
		ID:             user.ID,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Username:       user.Username,
		Email:          user.Email,
		Role:           user.Role,
		Score:          user.Score,
		PhoneNumber:    user.PhoneNumber,
		CreatedAt:      user.CreatedAt.Format(time.RFC3339),
		IsActive:       user.IsActive,
		ImagenPerfil:   imgB64,
		ImagenMimeType: user.ImagenMimeType,
	}

	return response, nil
}

func (s *userService) GetUsers(ctx context.Context) ([]dto.UserResponseDTO, e.ApiError) {
	users, apiErr := s.userRepo.GetUsers(ctx)
	if apiErr != nil {
		return nil, apiErr
	}

	var response []dto.UserResponseDTO
	for _, user := range users {
		var imgB64 string
		if len(user.ImagenPerfil) > 0 {
			imgB64 = base64.StdEncoding.EncodeToString(user.ImagenPerfil)
		}
		response = append(response, dto.UserResponseDTO{
			ID:             user.ID,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			Username:       user.Username,
			Email:          user.Email,
			Role:           user.Role,
			Score:          user.Score,
			PhoneNumber:    user.PhoneNumber,
			CreatedAt:      user.CreatedAt.Format(time.RFC3339),
			IsActive:       user.IsActive,
			ImagenPerfil:   imgB64,
			ImagenMimeType: user.ImagenMimeType,
		})
	}

	return response, nil
}

func (s *userService) UpdateUserById(ctx context.Context, id int, request dto.UserUpdateRequestDTO) (dto.UserResponseDTO, e.ApiError) {
	user, apiErr := s.userRepo.GetUserByID(ctx, id)
	if apiErr != nil {
		return dto.UserResponseDTO{}, apiErr
	}

	if request.Email != "" {
		if existing, err := s.userRepo.GetUserByEmail(ctx, request.Email); err == nil && existing.ID != id {
			return dto.UserResponseDTO{}, e.NewBadRequestApiError("email actualmente en uso")
		}
		user.Email = request.Email
	}

	if request.Username != "" {
		if existing, err := s.userRepo.GetUserByUsername(ctx, request.Username); err == nil && existing.ID != id {
			return dto.UserResponseDTO{}, e.NewBadRequestApiError("usuario actualmente en uso")
		}
		user.Username = request.Username
	}
	// Actualizar solo los campos enviados en el request
	if request.FirstName != "" {
		user.FirstName = request.FirstName
	}
	if request.LastName != "" {
		user.LastName = request.LastName
	}
	if request.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			return dto.UserResponseDTO{}, e.NewInternalServerApiError("failed to hash password", err)
		}
		user.Password = string(hashedPassword)
	}
	if request.PhoneNumber != "" {
		user.PhoneNumber = request.PhoneNumber
	}

	// Actualizar el usuario en la base de datos
	if apiErr := s.userRepo.UpdateUserByID(ctx, id, user); apiErr != nil {
		return dto.UserResponseDTO{}, apiErr
	}

	// Crear la respuesta
	response := dto.UserResponseDTO{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Username:    user.Username,
		Email:       user.Email,
		Score:       user.Score,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		IsActive:    user.IsActive,
	}

	return response, nil
}

func (s *userService) DeleteUserById(ctx context.Context, id int) e.ApiError {
	// Verificar si el usuario existe
	_, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return e.NewBadRequestApiError("user not found")
	}

	if apiErr := s.userRepo.DeleteUserByID(ctx, id); apiErr != nil {
		return apiErr
	}

	return nil
}

func (s *userService) UpdateRoleByUserId(ctx context.Context, id int, request dto.UserUpdateRoleRequestDTO) (dto.UserResponseDTO, e.ApiError) {
	user, apiErr := s.userRepo.GetUserByID(ctx, id)
	if apiErr != nil {
		return dto.UserResponseDTO{}, apiErr
	}

	// Solo actualiza el rol, los demás campos permanecen iguales
	user.Role = request.Role
	if apiErr := s.userRepo.UpdateUserByID(ctx, id, user); apiErr != nil {
		return dto.UserResponseDTO{}, apiErr
	}

	response := dto.UserResponseDTO{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Username:    user.Username,
		Email:       user.Email,
		Role:        user.Role,
		Score:       user.Score,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		IsActive:    user.IsActive,
	}

	return response, nil
}

func (s *userService) GetUserScoreByUserId(ctx context.Context, id int) (dto.UserScoreDtoSimplified, e.ApiError) {
	user, apiErr := s.userRepo.GetUserByID(ctx, id)
	if apiErr != nil {
		return dto.UserScoreDtoSimplified{}, apiErr
	}
	response := dto.UserScoreDtoSimplified{
		Username: user.Username,
		Score:    user.Score,
	}
	return response, nil
}

func (s *userService) UploadProfilePicture(ctx context.Context, userID int, data []byte, mimeType string) e.ApiError {
	// Verificar que el usuario exista
	if _, apiErr := s.userRepo.GetUserByID(ctx, userID); apiErr != nil {
		return apiErr
	}

	// Actualizar únicamente la imagen de perfil mediante el repositorio especializado
	if apiErr := s.userRepo.UpdateProfileImage(ctx, userID, data, mimeType); apiErr != nil {
		return apiErr
	}

	return nil
}

func (s *userService) GetScoreboard(ctx context.Context) ([]dto.UserScoreDtoSimplified, e.ApiError) {
	users, apiErr := s.userRepo.GetScoreboard(ctx)
	if apiErr != nil {
		return nil, apiErr
	}

	var scoreboard []dto.UserScoreDtoSimplified
	for _, user := range users {
		scoreboard = append(scoreboard, dto.UserScoreDtoSimplified{
			Username: user.Username,
			Score:    user.Score,
		})
	}

	return scoreboard, nil
}
