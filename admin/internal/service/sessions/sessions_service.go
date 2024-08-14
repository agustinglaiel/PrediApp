package service

import (
	dto "admin/internal/dto/sessions"
	model "admin/internal/model/sessions"
	repository "admin/internal/repository/sessions"
	e "admin/pkg/utils"
	"context"
	"time"
)

type sessionService struct{
	sessionsRepo repository.SessionRepository
}

type sessionServiceInterface interface{
	CreateSession(ctx context.Context, request dto.CreateSessionDTO) (dto.ResponseSessionDTO, e.ApiError)
	GetSessionById(ctx context.Context, sessionID uint) (dto.ResponseSessionDTO, e.ApiError)
}

func NewSessionService(sessionsRepo repository.SessionRepository) sessionServiceInterface{
	return &sessionService{
		sessionsRepo: sessionsRepo,
	}
}

func (s *sessionService) CreateSession(ctx context.Context, request dto.CreateSessionDTO) (dto.ResponseSessionDTO, e.ApiError) {
	// Convert DTO to Model
	newSession := &model.Session{
		CircuitKey:       request.CircuitKey,
		CircuitShortName: request.CircuitShortName,
		CountryCode:      request.CountryCode,
		CountryName:      request.CountryName,
		DateStart:        request.DateStart,
		DateEnd:          request.DateEnd,
		Location:         request.Location,
		SessionKey:       request.SessionKey,
		SessionName:      request.SessionName,
		SessionType:      request.SessionType,
		Year:             request.Year,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.sessionsRepo.CreateSession(ctx, newSession); err != nil {
		return dto.ResponseSessionDTO{}, e.NewInternalServerApiError("Error creando la sesión", err)
	}

	// Convert Model to Response DTO
	response := dto.ResponseSessionDTO{
		ID:               uint(newSession.SessionKey),
		CircuitKey:       newSession.CircuitKey,
		CircuitShortName: newSession.CircuitShortName,
		CountryCode:      newSession.CountryCode,
		CountryName:      newSession.CountryName,
		DateStart:        newSession.DateStart,
		DateEnd:          newSession.DateEnd,
		Location:         newSession.Location,
		SessionName:      newSession.SessionName,
		SessionType:      newSession.SessionType,
		Year:             newSession.Year,
	}

	return response, nil
}

func (s *sessionService) GetSessionById(ctx context.Context, sessionID uint) (dto.ResponseSessionDTO, e.ApiError) {
	session, err := s.sessionsRepo.GetSessionById(ctx, sessionID)
	if err != nil {
		return dto.ResponseSessionDTO{}, err
	}

	// Convert Model to Response DTO
	response := dto.ResponseSessionDTO{
		ID:               uint(session.SessionKey),
		CircuitKey:       session.CircuitKey,
		CircuitShortName: session.CircuitShortName,
		CountryCode:      session.CountryCode,
		CountryName:      session.CountryName,
		DateStart:        session.DateStart,
		DateEnd:          session.DateEnd,
		Location:         session.Location,
		SessionName:      session.SessionName,
		SessionType:      session.SessionType,
		Year:             session.Year,
	}

	return response, nil
}

func (s *sessionService) UpdateSessionById(ctx context.Context, sessionID uint, request dto.UpdateSessionDTO) (dto.ResponseSessionDTO, e.ApiError){
	session, apiErr := s.sessionsRepo.GetSessionById(ctx, sessionID)
	if apiErr != nil {
		return dto.ResponseSessionDTO{}, apiErr
	}

	session.CircuitKey = request.CircuitKey
	session.CircuitShortName = request.CircuitShortName
	session.CountryCode = request.CountryCode
	session.CountryKey = request.CountryKey
	session.CountryName = request.CountryName
	session.DateStart = request.DateStart
	session.DateEnd = request.DateEnd
	session.Location = request.Location
	session.SessionKey = request.SessionKey
	session.SessionName = request.SessionName
	session.SessionType = request.SessionType
	session.Year = request.Year

	if apiErr := s.sessionsRepo.UpdateSessionById(ctx, session); apiErr != nil {

	}
}