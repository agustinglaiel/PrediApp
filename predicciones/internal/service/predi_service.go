package service

import (
	"context"
	"predicciones/internal/dto"
	"predicciones/internal/model"
	"predicciones/internal/repository"
	e "predicciones/pkg/utils"
)

type prediService struct {
    prodeRepo repository.ProdeRepository
}

type PrediServiceInterface interface {
    CreateProdeCarrera(ctx context.Context, dto dto.CreateProdeCarreraDTO) (dto.ResponseProdeCarreraDTO, e.ApiError)
    CreateProdeSession(ctx context.Context, dto dto.CreateProdeSessionDTO) (dto.ResponseProdeSessionDTO, e.ApiError)
    GetProdeCarreraByID(ctx context.Context, prodeID uint) (dto.ResponseProdeCarreraDTO, e.ApiError)
    GetProdeSessionByID(ctx context.Context, prodeID uint) (dto.ResponseProdeSessionDTO, e.ApiError)
    GetProdesByUserID(ctx context.Context, userID uint) ([]dto.ResponseProdeCarreraDTO, []dto.ResponseProdeSessionDTO, e.ApiError)
    UpdateProdeCarrera(ctx context.Context, dto dto.UpdateProdeCarreraDTO) (dto.ResponseProdeCarreraDTO, e.ApiError)
    UpdateProdeSession(ctx context.Context, dto dto.UpdateProdeSessionDTO) (dto.ResponseProdeSessionDTO, e.ApiError)
    DeleteProdeByID(ctx context.Context, prodeID uint, userID uint) e.ApiError
}

func NewPrediService(prodeRepo repository.ProdeRepository) PrediServiceInterface {
    return &prediService{
        prodeRepo: prodeRepo,
    }
}

func (s *prediService) CreateProdeCarrera(ctx context.Context, request dto.CreateProdeCarreraDTO) (dto.ResponseProdeCarreraDTO, e.ApiError) {
    prode := model.ProdeCarrera{
        UserID:     uint(request.UserID),
        EventID:    uint(request.EventID),
        P1:         uint(request.P1),
        P2:         uint(request.P2),
        P3:         uint(request.P3),
        P4:         uint(request.P4),
        P5:         uint(request.P5),
        FastestLap: uint(request.FastestLap),
        VSC:        request.VSC,
        SC:         request.SC,
        DNF:        request.DNF,
    }

    if err := s.prodeRepo.CreateProdeCarrera(ctx, &prode); err != nil {
        return dto.ResponseProdeCarreraDTO{}, err
    }

    response := dto.ResponseProdeCarreraDTO{
        ID:         uint(prode.ID),
        UserID:     uint(prode.UserID),
        EventID:    uint(prode.EventID),
        P1:         uint(prode.P1),
        P2:         uint(prode.P2),
        P3:         uint(prode.P3),
        P4:         uint(prode.P4),
        P5:         uint(prode.P5),
        FastestLap: uint(prode.FastestLap),
        VSC:        prode.VSC,
        SC:         prode.SC,
        DNF:        prode.DNF,
    }

    return response, nil
}

func (s *prediService) CreateProdeSession(ctx context.Context, request dto.CreateProdeSessionDTO) (dto.ResponseProdeSessionDTO, e.ApiError) {
    // Crear una nueva instancia de ProdeSession
    prode := model.ProdeSession{
        UserID:  uint(request.UserID),
        EventID: uint(request.EventID),
        P1:      uint(request.P1),
        P2:      uint(request.P2),
        P3:      uint(request.P3),
    }

    // Insertar el nuevo pronóstico de sesión en la base de datos
    if err := s.prodeRepo.CreateProdeSession(ctx, &prode); err != nil {
        return dto.ResponseProdeSessionDTO{}, err
    }

    // Construir la respuesta DTO con los datos insertados
    response := dto.ResponseProdeSessionDTO{
        ID:       uint(prode.ID),
        UserID:   uint(prode.UserID),
        EventID:  uint(prode.EventID),
        P1:       uint(prode.P1),
        P2:       uint(prode.P2),
        P3:       uint(prode.P3),
    }

    return response, nil
}

func (s *prediService) GetProdeCarreraByID(ctx context.Context, id uint) (dto.ResponseProdeCarreraDTO, e.ApiError) {
    prode, err := s.prodeRepo.GetProdeCarreraByID(ctx, uint(id))
    if err != nil {
        return dto.ResponseProdeCarreraDTO{}, e.NewNotFoundApiError("prode not found")
    }

    response := dto.ResponseProdeCarreraDTO{
        ID:         uint(prode.ID),
        UserID:     uint(prode.UserID),
        EventID:    uint(prode.EventID),
        P1:         uint(prode.P1),
        P2:         uint(prode.P2),
        P3:         uint(prode.P3),
        P4:         uint(prode.P4),
        P5:         uint(prode.P5),
        FastestLap: uint(prode.FastestLap),
        VSC:        prode.VSC,
        SC:         prode.SC,
        DNF:        prode.DNF,
    }

    return response, nil
}

func (s *prediService) GetProdeSessionByID(ctx context.Context, prodeID uint) (dto.ResponseProdeSessionDTO, e.ApiError) {
    prode, err := s.prodeRepo.GetProdeSessionByID(ctx, prodeID)
    if err != nil {
        return dto.ResponseProdeSessionDTO{}, err
    }

    response := dto.ResponseProdeSessionDTO{
        ID:      prode.ID,
        UserID:  prode.UserID,
        EventID: prode.EventID,
        P1:      prode.P1,
        P2:      prode.P2,
        P3:      prode.P3,
    }

    return response, nil
}

func (s *prediService) GetProdesByUserID(ctx context.Context, userID uint) ([]dto.ResponseProdeCarreraDTO, []dto.ResponseProdeSessionDTO, e.ApiError) {
    prodesCarrera, prodesSession, err := s.prodeRepo.GetProdesByUserID(ctx, userID)
    if err != nil {
        return nil, nil, err
    }

    var responseCarrera []dto.ResponseProdeCarreraDTO
    for _, prode := range prodesCarrera {
        responseCarrera = append(responseCarrera, dto.ResponseProdeCarreraDTO{
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

    var responseSession []dto.ResponseProdeSessionDTO
    for _, prode := range prodesSession {
        responseSession = append(responseSession, dto.ResponseProdeSessionDTO{
            ID:      prode.ID,
            UserID:  prode.UserID,
            EventID: prode.EventID,
            P1:      prode.P1,
            P2:      prode.P2,
            P3:      prode.P3,
        })
    }

    return responseCarrera, responseSession, nil
}

func (s *prediService) UpdateProdeCarrera(ctx context.Context, request dto.UpdateProdeCarreraDTO) (dto.ResponseProdeCarreraDTO, e.ApiError) {
    // Convertir ProdeID a uint
    prode, err := s.prodeRepo.GetProdeCarreraByID(ctx, uint(request.ProdeID))
    if err != nil {
        return dto.ResponseProdeCarreraDTO{}, e.NewNotFoundApiError("prode not found")
    }

    // Actualizar los campos del pronóstico
    prode.P1 = uint(request.P1)
    prode.P2 = uint(request.P2)
    prode.P3 = uint(request.P3)
    prode.P4 = uint(request.P4)
    prode.P5 = uint(request.P5)
    prode.FastestLap = uint(request.FastestLap)
    prode.VSC = request.VSC
    prode.SC = request.SC
    prode.DNF = request.DNF

    if err := s.prodeRepo.UpdateProdeCarrera(ctx, prode); err != nil {
        return dto.ResponseProdeCarreraDTO{}, e.NewInternalServerApiError("error updating prode", err)
    }

    response := dto.ResponseProdeCarreraDTO{
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

func (s *prediService) UpdateProdeSession(ctx context.Context, request dto.UpdateProdeSessionDTO) (dto.ResponseProdeSessionDTO, e.ApiError) {
    // Obtener el pronóstico de la sesión por ID
    prode, err := s.prodeRepo.GetProdeSessionByID(ctx, uint(request.ProdeID))
    if err != nil {
        return dto.ResponseProdeSessionDTO{}, e.NewNotFoundApiError("prode not found")
    }

    // Actualizar los campos del pronóstico de la sesión
    prode.P1 = uint(request.P1)
    prode.P2 = uint(request.P2)
    prode.P3 = uint(request.P3)

    // Actualizar el pronóstico en la base de datos
    if err := s.prodeRepo.UpdateProdeSession(ctx, prode); err != nil {
        return dto.ResponseProdeSessionDTO{}, e.NewInternalServerApiError("error updating prode session", err)
    }

    // Preparar la respuesta
    response := dto.ResponseProdeSessionDTO{
        ID:      prode.ID,
        UserID:  prode.UserID,
        EventID: prode.EventID,
        P1:      prode.P1,
        P2:      prode.P2,
        P3:      prode.P3,
    }

    return response, nil
}

func (s *prediService) DeleteProdeByID(ctx context.Context, prodeID uint, userID uint) e.ApiError {
    prodeCarrera, err := s.prodeRepo.GetProdeCarreraByID(ctx, prodeID)
    if err == nil && prodeCarrera.UserID == userID {
        if apiErr := s.prodeRepo.DeleteProdeCarreraByID(ctx, prodeID); apiErr != nil {
            return apiErr
        }
        return nil
    }

    prodeSession, err := s.prodeRepo.GetProdeSessionByID(ctx, prodeID)
    if err == nil && prodeSession.UserID == userID {
        if apiErr := s.prodeRepo.DeleteProdeSessionByID(ctx, prodeID); apiErr != nil {
            return apiErr
        }
        return nil
    }

    return e.NewNotFoundApiError("prode not found")
}
