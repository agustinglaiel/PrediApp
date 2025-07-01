package service

import (
	"context"
	"time"

	"prediapp.local/users/internal/dto"
	"prediapp.local/users/internal/model"
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
	GetUserByUsername(ctx context.Context, username string) (dto.UserResponseDTO, e.ApiError)
	GetUsers(ctx context.Context) ([]dto.UserResponseDTO, e.ApiError)
	UpdateUserById(ctx context.Context, id int, request dto.UserUpdateRequestDTO) (dto.UserResponseDTO, e.ApiError)
	UpdateUserByUsername(ctx context.Context, username string, request dto.UserUpdateRequestDTO) (dto.UserResponseDTO, e.ApiError)
	DeleteUserById(ctx context.Context, id int) e.ApiError
	DeleteUserByUsername(ctx context.Context, username string) e.ApiError
	UpdateRoleByUserId(ctx context.Context, id int, request dto.UserUpdateRoleRequestDTO) (dto.UserResponseDTO, e.ApiError)
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

	response := dto.UserResponseDTO{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Score:     user.Score,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		IsActive:  user.IsActive,
	}

	return response, nil
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (dto.UserResponseDTO, e.ApiError) {
	user, apiErr := s.userRepo.GetUserByUsername(ctx, username)
	if apiErr != nil {
		return dto.UserResponseDTO{}, apiErr
	}

	response := dto.UserResponseDTO{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Score:     user.Score,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		IsActive:  user.IsActive,
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
		response = append(response, dto.UserResponseDTO{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			Score:     user.Score,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			IsActive:  user.IsActive,
		})
	}

	return response, nil
}

func (s *userService) UpdateUserById(ctx context.Context, id int, request dto.UserUpdateRequestDTO) (dto.UserResponseDTO, e.ApiError) {
	user, apiErr := s.userRepo.GetUserByID(ctx, id)
	if apiErr != nil {
		return dto.UserResponseDTO{}, apiErr
	}

	// Actualizar solo los campos enviados en el request
	if request.FirstName != "" {
		user.FirstName = request.FirstName
	}
	if request.LastName != "" {
		user.LastName = request.LastName
	}
	if request.Username != "" {
		user.Username = request.Username
	}
	if request.Email != "" {
		user.Email = request.Email
	}
	if request.Role != "" {
		user.Role = request.Role
	}
	validRoles := map[string]bool{"user": true, "admin": true}
	if request.Role != "" && !validRoles[request.Role] {
		return dto.UserResponseDTO{}, e.NewBadRequestApiError("invalid role")
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
	// Manejar booleanos como `IsActive` explícitamente
	user.IsActive = request.IsActive

	// Actualizar el usuario en la base de datos
	if apiErr := s.userRepo.UpdateUserByID(ctx, id, user); apiErr != nil {
		return dto.UserResponseDTO{}, apiErr
	}

	// Crear la respuesta
	response := dto.UserResponseDTO{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Score:     user.Score,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		IsActive:  user.IsActive,
	}

	return response, nil
}

func (s *userService) UpdateUserByUsername(ctx context.Context, username string, request dto.UserUpdateRequestDTO) (dto.UserResponseDTO, e.ApiError) {
	user, apiErr := s.userRepo.GetUserByUsername(ctx, username)
	if apiErr != nil {
		return dto.UserResponseDTO{}, apiErr
	}

	// Actualizar solo los campos enviados en el request
	if request.FirstName != "" {
		user.FirstName = request.FirstName
	}
	if request.LastName != "" {
		user.LastName = request.LastName
	}
	if request.Username != "" {
		user.Username = request.Username
	}
	if request.Email != "" {
		user.Email = request.Email
	}
	if request.Role != "" {
		user.Role = request.Role
	}
	validRoles := map[string]bool{"user": true, "admin": true}
	if request.Role != "" && !validRoles[request.Role] {
		return dto.UserResponseDTO{}, e.NewBadRequestApiError("invalid role")
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
	// Manejar booleanos como `IsActive` explícitamente
	user.IsActive = request.IsActive

	if apiErr := s.userRepo.UpdateUserByUsername(ctx, username, user); apiErr != nil {
		return dto.UserResponseDTO{}, apiErr
	}

	response := dto.UserResponseDTO{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Score:     user.Score,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		IsActive:  user.IsActive,
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

func (s *userService) DeleteUserByUsername(ctx context.Context, username string) e.ApiError {
	if apiErr := s.userRepo.DeleteUserByUsername(ctx, username); apiErr != nil {
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
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Score:     user.Score,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		IsActive:  user.IsActive,
	}

	return response, nil
}
