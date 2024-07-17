package service

import (
	"context"
	"log"
	"time"
	"users/internal/dto"
	"users/internal/model"
	"users/internal/repository"
	e "users/pkg/utils"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
    userRepo repository.UserRepository
}

type UserServiceInterface interface {
    SignUp(ctx context.Context, request dto.UserSignUpRequestDTO) (dto.UserSignUpResponseDTO, e.ApiError)
    Login(ctx context.Context, request dto.UserLoginRequestDTO) (dto.UserLoginResponseDTO, e.ApiError)
    OAuthSignIn(ctx context.Context, request dto.GoogleOAuthRequestDTO) (dto.GoogleOAuthResponseDTO, e.ApiError)
}

func NewUserService(userRepo repository.UserRepository) UserServiceInterface {
    return &userService{
        userRepo: userRepo,
    }
}

func (s *userService) SignUp(ctx context.Context, request dto.UserSignUpRequestDTO) (dto.UserSignUpResponseDTO, e.ApiError) {
    log.Printf("Checking if email already exists: %s", request.Email)
    if _, err := s.userRepo.GetUserByEmail(ctx, request.Email); err == nil {
        return dto.UserSignUpResponseDTO{}, e.NewBadRequestApiError("email already exists")
    }

    log.Printf("Checking if username already exists: %s", request.Username)
    if _, err := s.userRepo.GetUserByUsername(ctx, request.Username); err == nil {
        return dto.UserSignUpResponseDTO{}, e.NewBadRequestApiError("username already exists")
    }

    // Hash de la contraseña
    log.Printf("Hashing password for user: %s", request.Username)
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
    if err != nil {
        return dto.UserSignUpResponseDTO{}, e.NewInternalServerApiError("error hashing password", err)
    }

    // Crear el nuevo usuario
    newUser := &model.User{
        ID:              primitive.NewObjectID(),
        FirstName:       request.FirstName,
        LastName:        request.LastName,
        Username:        request.Username,
        Email:           request.Email,
        Password:        string(hashedPassword),
        Role:            "user",
        Score:           0,
        CreatedAt:       time.Now(),
        IsActive:        true,
        IsEmailVerified: false,
    }

    log.Printf("Creating user: %+v", newUser)
    if err := s.userRepo.CreateUser(ctx, newUser); err != nil {
        return dto.UserSignUpResponseDTO{}, e.NewInternalServerApiError("error creating user", err)
    }

    response := dto.UserSignUpResponseDTO{
        ID:        newUser.ID.Hex(),
        FirstName: newUser.FirstName,
        LastName:  newUser.LastName,
        Username:  newUser.Username,
        Email:     newUser.Email,
        Role:      newUser.Role,
        CreatedAt: newUser.CreatedAt.Format(time.RFC3339),
    }

    log.Printf("User created successfully: %+v", response)
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
        ID:        user.ID.Hex(),
        FirstName: user.FirstName,
        LastName:  user.LastName,
        Username:  user.Username,
        Email:     user.Email,
        Role:      user.Role,
        Token:     "dummy-jwt-token", // Reemplazar con el token real
    }

    return response, nil
}

func (s *userService) OAuthSignIn(ctx context.Context, request dto.GoogleOAuthRequestDTO) (dto.GoogleOAuthResponseDTO, e.ApiError) {
    googleUser, err := VerifyGoogleToken(request.GoogleToken)
    if err != nil {
        return dto.GoogleOAuthResponseDTO{}, e.NewBadRequestApiError("invalid Google token")
    }

    user, apiErr := s.userRepo.GetUserByEmail(ctx, googleUser.Email)
    if apiErr != nil {
        if apiErr.Status() != 404 {
            return dto.GoogleOAuthResponseDTO{}, apiErr
        }

        user = &model.User{
            ID:              primitive.NewObjectID(),
            FirstName:       googleUser.FirstName,
            LastName:        googleUser.LastName,
            Username:        googleUser.NickName,
            Email:           googleUser.Email,
            Role:            "user",
            Provider:        "google",
            ProviderID:      googleUser.UserID,
            AvatarURL:       googleUser.AvatarURL,
            IsActive:        true,
            IsEmailVerified: true,
            CreatedAt:       time.Now(),
        }

        if apiErr := s.userRepo.CreateUser(ctx, user); apiErr != nil {
            return dto.GoogleOAuthResponseDTO{}, apiErr
        }
    } else {
        user.Provider = "google"
        user.ProviderID = googleUser.UserID
        user.AvatarURL = googleUser.AvatarURL
        user.IsActive = true
        now := time.Now()
        user.LastLoginAt = &now

        if apiErr := s.userRepo.UpdateUser(ctx, user); apiErr != nil {
            return dto.GoogleOAuthResponseDTO{}, apiErr
        }
    }

    response := dto.GoogleOAuthResponseDTO{
        ID:          user.ID.Hex(),
        FirstName:   user.FirstName,
        LastName:    user.LastName,
        Username:    user.Username,
        Email:       user.Email,
        Role:        user.Role,
        Token:       "dummy-jwt-token", // Reemplazar con el token real
        Provider:    user.Provider,
        ProviderID:  user.ProviderID,
        AvatarURL:   user.AvatarURL,
    }

    return response, nil
}

func VerifyGoogleToken(googleToken string) (*goth.User, e.ApiError) {
    provider := google.New("client-id", "client-secret", "redirect-url", "profile", "email")
    goth.UseProviders(provider)

    session, err := provider.UnmarshalSession(`{"AuthURL":"","AccessToken":"` + googleToken + `"}`)
    if err != nil {
        return nil, e.NewInternalServerApiError("failed to unmarshal session", err)
    }

    user, err := provider.FetchUser(session)
    if err != nil {
        return nil, e.NewInternalServerApiError("failed to verify Google token", err)
    }

    return &user, nil
}