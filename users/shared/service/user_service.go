package service

import (
	"context"
	"time"
	e "users/pkg/utils"
	"users/shared/dto"
	"users/shared/model"
	"users/shared/repository"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
    userRepo repository.UserRepository
}

type UserServiceInterface interface {
    SignUp(ctx context.Context, request dto.UserSignUpRequestDTO) (dto.UserSignUpResponseDTO, e.ApiError)
    Login(ctx context.Context, request dto.UserLoginRequestDTO) (dto.UserLoginResponseDTO, e.ApiError)
    // OAuthSignIn(ctx context.Context, request dto.GoogleOAuthRequestDTO) (dto.GoogleOAuthResponseDTO, e.ApiError)
    GetUserById(ctx context.Context, id int) (dto.UserResponseDTO, e.ApiError)
    GetUserByUsername(ctx context.Context, username string) (dto.UserResponseDTO, e.ApiError)
    GetUsers(ctx context.Context) ([]dto.UserResponseDTO, e.ApiError)
    UpdateUserById(ctx context.Context, id int, request dto.UserUpdateRequestDTO) (dto.UserResponseDTO, e.ApiError)
    UpdateUserByUsername(ctx context.Context, username string, request dto.UserUpdateRequestDTO) (dto.UserResponseDTO, e.ApiError)
    DeleteUserById(ctx context.Context, id int) e.ApiError
    DeleteUserByUsername(ctx context.Context, username string) e.ApiError
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

    // Genarar el token JWT para el usuario regittado
    // token, err := jwt.GenerateJWT(newUser.ID, newUser.Role)
    // if err != nil {
    //     return dto.UserSignUpResponseDTO{}, e.NewInternalServerApiError("error generating token", err)
    // }

	response := dto.UserSignUpResponseDTO{
		ID:        newUser.ID,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Username:  newUser.Username,
		Email:     newUser.Email,
		Role:      newUser.Role,
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
	}

	return response, nil
}

// func (s *userService) OAuthSignIn(ctx context.Context, request dto.GoogleOAuthRequestDTO) (dto.GoogleOAuthResponseDTO, e.ApiError) {
// 	googleUser, err := VerifyGoogleToken(request.GoogleToken)
// 	if err != nil {
// 		return dto.GoogleOAuthResponseDTO{}, e.NewBadRequestApiError("invalid Google token")
// 	}
// 	user, apiErr := s.userRepo.GetUserByEmail(ctx, googleUser.Email)
// 	if apiErr != nil {
// 		if apiErr.Status() != 404 {
// 			return dto.GoogleOAuthResponseDTO{}, apiErr
// 		}
// 		user = &model.User{
// 			FirstName:       googleUser.FirstName,
// 			LastName:        googleUser.LastName,
// 			Username:        googleUser.NickName,
// 			Email:           googleUser.Email,
// 			Role:            "user",
// 			Provider:        "google",
// 			ProviderID:      googleUser.UserID,
// 			AvatarURL:       googleUser.AvatarURL,
// 			IsActive:        true,
// 			IsEmailVerified: true,
// 			CreatedAt:       time.Now(),
// 		}
// 		if apiErr := s.userRepo.CreateUser(ctx, user); apiErr != nil {
// 			return dto.GoogleOAuthResponseDTO{}, apiErr
// 		}
// 	} else {
// 		user.Provider = "google"
// 		user.ProviderID = googleUser.UserID
// 		user.AvatarURL = googleUser.AvatarURL
// 		user.IsActive = true
// 		now := time.Now()
// 		user.LastLoginAt = &now
// 		if apiErr := s.userRepo.UpdateUserByID(ctx, user.ID, user); apiErr != nil {
// 			return dto.GoogleOAuthResponseDTO{}, apiErr
// 		}
// 	}
// 	response := dto.GoogleOAuthResponseDTO{
// 		ID:          user.ID,
// 		FirstName:   user.FirstName,
// 		LastName:    user.LastName,
// 		Username:    user.Username,
// 		Email:       user.Email,
// 		Role:        user.Role,
// 		Token:       "dummy-jwt-token", // Reemplazar con el token real
// 		Provider:    user.Provider,
// 		ProviderID:  user.ProviderID,
// 		AvatarURL:   user.AvatarURL,
// 	}
// 	return response, nil
// }

// func VerifyGoogleToken(googleToken string) (*goth.User, e.ApiError) {
// 	provider := google.New("client-id", "client-secret", "redirect-url", "profile", "email")
// 	goth.UseProviders(provider)
// 	session, err := provider.UnmarshalSession(`{"AuthURL":"","AccessToken":"` + googleToken + `"}`)
// 	if err != nil {
// 		return nil, e.NewInternalServerApiError("failed to unmarshal session", err)
// 	}
// 	user, err := provider.FetchUser(session)
// 	if err != nil {
// 		return nil, e.NewInternalServerApiError("failed to verify Google token", err)
// 	}
// 	return &user, nil
// }

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

    user.FirstName = request.FirstName
    user.LastName = request.LastName
    user.Username = request.Username
    user.Email = request.Email
    user.IsActive = request.IsActive

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

func (s *userService) UpdateUserByUsername(ctx context.Context, username string, request dto.UserUpdateRequestDTO) (dto.UserResponseDTO, e.ApiError) {
    user, apiErr := s.userRepo.GetUserByUsername(ctx, username)
    if apiErr != nil {
        return dto.UserResponseDTO{}, apiErr
    }

    user.FirstName = request.FirstName
    user.LastName = request.LastName
    user.Username = request.Username
    user.Email = request.Email
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
    _ , err := s.userRepo.GetUserByID(ctx, id)
    if err != nil {
        return  e.NewBadRequestApiError("user not found")
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