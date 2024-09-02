package service

import (
	prodes "admin/internal/dto/prodes"
	model "admin/internal/model/prodes"
	repository "admin/internal/repository/prodes"
	e "admin/pkg/utils"
	"context"
)

type prodeService struct {
	prodeRepo repository.ProdeRepository
}

type ProdeServiceInterface interface {
	CreateProdeCarrera(ctx context.Context, request prodes.CreateProdeCarreraDTO) (prodes.ResponseProdeCarreraDTO, e.ApiError)
	CreateProdeSession(ctx context.Context, request prodes.CreateProdeSessionDTO) (prodes.ResponseProdeSessionDTO, e.ApiError)
	UpdateProdeCarrera(ctx context.Context, request prodes.UpdateProdeCarreraDTO) (prodes.ResponseProdeCarreraDTO, e.ApiError)
	UpdateProdeSession(ctx context.Context, request prodes.UpdateProdeSessionDTO) (prodes.ResponseProdeSessionDTO, e.ApiError)
	DeleteProdeCarrera(ctx context.Context, prodeID int) e.ApiError
	DeleteProdeSession(ctx context.Context, prodeID int) e.ApiError
	GetProdesByUserId(ctx context.Context, userID int) ([]prodes.ResponseProdeCarreraDTO, []prodes.ResponseProdeSessionDTO, e.ApiError)
}

func NewPrediService(prodeRepo repository.ProdeRepository) ProdeServiceInterface {
    return &prodeService{
        prodeRepo: prodeRepo,
    }
}

func (s *prodeService) CreateProdeCarrera(ctx context.Context, request prodes.CreateProdeCarreraDTO) (prodes.ResponseProdeCarreraDTO, e.ApiError) {
	// Convertir DTO a modelo
	prode := model.ProdeCarrera{
		UserID:     request.UserID,
		EventID:    request.EventID,
		P1:         request.P1,
		P2:         request.P2,
		P3:         request.P3,
		P4:         request.P4,
		P5:         request.P5,
		FastestLap: request.FastestLap,
		VSC:        request.VSC,
		SC:         request.SC,
		DNF:        request.DNF,
	}

	// Crear el pronóstico de carrera en la base de datos
	err := s.prodeRepo.CreateProdeCarrera(ctx, &prode)
	if err != nil {
		return prodes.ResponseProdeCarreraDTO{}, err
	}

	// Convertir el modelo a DTO de respuesta
	response := prodes.ResponseProdeCarreraDTO{
		ID:         prode.ID,
		UserID:     prode.UserID,
		EventID:    prode.EventID,
		P1:         prode.P1,
		P2:         prode.P2,
		P3:         prode.P3,
		P4:         prode.P4,
		P5:         prode.P5,
		FastestLap: prode.FastestLap,
		VSC:        prode.VSC,
		SC:         prode.SC,
		DNF:        prode.DNF,
	}

	return response, nil
}

func (s *prodeService) CreateProdeSession(ctx context.Context, request prodes.CreateProdeSessionDTO) (prodes.ResponseProdeSessionDTO, e.ApiError) {
	prode := model.ProdeSession{
		UserID:  request.UserID,
		EventID: request.EventID,
		P1:      request.P1,
		P2:      request.P2,
		P3:      request.P3,
	}
	err := s.prodeRepo.CreateProdeSession(ctx, &prode)
	if err != nil {
		return prodes.ResponseProdeSessionDTO{}, err
	}
	response := prodes.ResponseProdeSessionDTO{
		ID:      prode.ID,
		UserID:  prode.UserID,
		EventID: prode.EventID,
		P1:      prode.P1,
		P2:      prode.P2,
		P3:      prode.P3,
	}
	return response, nil
}

func (s *prodeService) UpdateProdeCarrera(ctx context.Context, request prodes.UpdateProdeCarreraDTO) (prodes.ResponseProdeCarreraDTO, e.ApiError) {
	prode := model.ProdeCarrera{
		ID:         request.ProdeID,
		UserID:     request.UserID,
		EventID:    request.EventID,
		P1:         request.P1,
		P2:         request.P2,
		P3:         request.P3,
		P4:         request.P4,
		P5:         request.P5,
		FastestLap: request.FastestLap,
		VSC:        request.VSC,
		SC:         request.SC,
		DNF:        request.DNF,
	}
	err := s.prodeRepo.UpdateProdeCarrera(ctx, &prode)
	if err != nil {
		return prodes.ResponseProdeCarreraDTO{}, err
	}
	response := prodes.ResponseProdeCarreraDTO{
		ID:         prode.ID,
		UserID:     prode.UserID,
		EventID:    prode.EventID,
		P1:         prode.P1,
		P2:         prode.P2,
		P3:         prode.P3,
		P4:         prode.P4,
		P5:         prode.P5,
		FastestLap: prode.FastestLap,
		VSC:        prode.VSC,
		SC:         prode.SC,
		DNF:        prode.DNF,
	}
	return response, nil
}

func (s *prodeService) UpdateProdeSession(ctx context.Context, request prodes.UpdateProdeSessionDTO) (prodes.ResponseProdeSessionDTO, e.ApiError) {
	prode := model.ProdeSession{
		ID:      request.ProdeID,
		UserID:  request.UserID,
		EventID: request.EventID,
		P1:      request.P1,
		P2:      request.P2,
		P3:      request.P3,
	}
	err := s.prodeRepo.UpdateProdeSession(ctx, &prode)
	if err != nil {
		return prodes.ResponseProdeSessionDTO{}, err
	}
	response := prodes.ResponseProdeSessionDTO{
		ID:      prode.ID,
		UserID:  prode.UserID,
		EventID: prode.EventID,
		P1:      prode.P1,
		P2:      prode.P2,
		P3:      prode.P3,
	}
	return response, nil
}

func (s *prodeService) DeleteProdeCarrera(ctx context.Context, prodeID int) e.ApiError {
	// Buscar el prode de carrera por ID
	prode, err := s.prodeRepo.GetProdeCarreraByID(ctx, prodeID)
	if err != nil {
		return err
	}

	// Eliminar el prode de carrera
	if err := s.prodeRepo.DeleteProdeCarreraByID(ctx, prode.ID, prode.UserID); err != nil {
		return err
	}

	return nil
}

func (s *prodeService) DeleteProdeSession(ctx context.Context, prodeID int) e.ApiError {
	// Buscar el prode de sesión por ID
	prode, err := s.prodeRepo.GetProdeSessionByID(ctx, prodeID)
	if err != nil {
		return err
	}

	// Eliminar el prode de sesión
	if err := s.prodeRepo.DeleteProdeSessionByID(ctx, prode.ID, prode.UserID); err != nil {
		return err
	}

	return nil
}

func (s *prodeService) GetProdesByUserId(ctx context.Context, userID int) ([]prodes.ResponseProdeCarreraDTO, []prodes.ResponseProdeSessionDTO, e.ApiError) {
	carreraProdes, sessionProdes, err := s.prodeRepo.GetProdesByUserID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	
	var carreraResponses []prodes.ResponseProdeCarreraDTO
	for _, prode := range carreraProdes {
		carreraResponses = append(carreraResponses, prodes.ResponseProdeCarreraDTO{
			ID:         prode.ID,
			UserID:     prode.UserID,
			EventID:    prode.EventID,
			P1:         prode.P1,
			P2:         prode.P2,
			P3:         prode.P3,
			P4:         prode.P4,
			P5:         prode.P5,
			FastestLap: prode.FastestLap,
			VSC:        prode.VSC,
			SC:         prode.SC,
			DNF:        prode.DNF,
		})
	}

	var sessionResponses []prodes.ResponseProdeSessionDTO
	for _, prode := range sessionProdes {
		sessionResponses = append(sessionResponses, prodes.ResponseProdeSessionDTO{
			ID:      prode.ID,
			UserID:  prode.UserID,
			EventID: prode.EventID,
			P1:      prode.P1,
			P2:      prode.P2,
			P3:      prode.P3,
		})
	}

	return carreraResponses, sessionResponses, nil
}

/*
// GetSessionNameAndType retrieves the session name and session type based on the event ID.
func (s *prodeService) GetSessionNameAndType(ctx context.Context, eventID uint) (string, string, e.ApiError) {
	session, apiErr := s.prodeRepo.GetSessionByEventID(ctx, eventID)
	if apiErr != nil {
		return "", "", apiErr
	}

	return session.SessionName, session.SessionType, nil
}
*/