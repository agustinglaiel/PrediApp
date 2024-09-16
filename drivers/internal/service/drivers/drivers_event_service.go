package service

import (
	"context"
	dto "drivers/internal/dto/drivers"
	model "drivers/internal/model/drivers"
	repository "drivers/internal/repository/drivers"
	e "drivers/pkg/utils"
)

type driverEventService struct {
	driverEventRepo repository.DriverEventRepository
	driverRepo      repository.DriverRepository
}

type DriverEventServiceInterface interface {
	AddDriverToEvent(ctx context.Context, request dto.DriverEventDTO) (dto.ResponseDriverEventDTO, e.ApiError)
	RemoveDriverFromEvent(ctx context.Context, driverEventID uint) e.ApiError
	ListDriversByEvent(ctx context.Context, eventID uint) ([]dto.ResponseDriverDTO, e.ApiError)
	ListEventsByDriver(ctx context.Context, driverID uint) ([]dto.ResponseEventDTO, e.ApiError)
	CheckDriverInEvent(ctx context.Context, driverID, eventID uint) (bool, e.ApiError)
}

func NewDriverEventService(driverEventRepo repository.DriverEventRepository, driverRepo repository.DriverRepository) DriverEventServiceInterface {
	return &driverEventService{
		driverEventRepo: driverEventRepo,
		driverRepo:      driverRepo,
	}
}

func (s *driverEventService) AddDriverToEvent(ctx context.Context, request dto.DriverEventDTO) (dto.ResponseDriverEventDTO, e.ApiError) {
	// Verificar si el piloto ya está asignado a este evento
	exists, err := s.CheckDriverInEvent(ctx, request.DriverID, request.EventID)
	if err != nil {
		return dto.ResponseDriverEventDTO{}, err
	}
	if exists {
		return dto.ResponseDriverEventDTO{}, e.NewBadRequestApiError("El piloto ya está asignado a este evento")
	}

	// Crear la asignación piloto-evento
	driverEvent := &model.DriverEvent{
		EventID:  request.EventID,
		DriverID: request.DriverID,
	}

	if err := s.driverEventRepo.AddDriverToEvent(ctx, driverEvent); err != nil {
		return dto.ResponseDriverEventDTO{}, err
	}

	response := dto.ResponseDriverEventDTO{
		ID:       driverEvent.ID,
		EventID:  driverEvent.EventID,
		DriverID: driverEvent.DriverID,
	}

	return response, nil
}

func (s *driverEventService) RemoveDriverFromEvent(ctx context.Context, driverEventID uint) e.ApiError {
	// Eliminar la asignación piloto-evento
	if err := s.driverEventRepo.RemoveDriverFromEvent(ctx, driverEventID); err != nil {
		return e.NewInternalServerApiError("Error eliminando piloto del evento", err)
	}
	return nil
}

func (s *driverEventService) ListDriversByEvent(ctx context.Context, eventID uint) ([]dto.ResponseDriverDTO, e.ApiError) {
	drivers, err := s.driverEventRepo.ListDriversByEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}

	var response []dto.ResponseDriverDTO
	for _, driver := range drivers {
		response = append(response, dto.ResponseDriverDTO{
			ID:             driver.ID,
			BroadcastName:  driver.BroadcastName,
			CountryCode:    driver.CountryCode,
			DriverNumber:   driver.DriverNumber,
			FirstName:      driver.FirstName,
			LastName:       driver.LastName,
			FullName:       driver.FullName,
			NameAcronym:    driver.NameAcronym,
			TeamName:       driver.TeamName,
		})
	}

	return response, nil
}

func (s *driverEventService) ListEventsByDriver(ctx context.Context, driverID uint) ([]dto.ResponseEventDTO, e.ApiError) {
	// Obtener eventos por piloto
	events, err := s.driverEventRepo.ListEventsByDriver(ctx, driverID)
	if err != nil {
		return nil, err
	}

	var response []dto.ResponseEventDTO
	for _, event := range events {
		response = append(response, dto.ResponseEventDTO{
			ID:           event.ID,
			SessionID:    event.SessionID,
			Date:         event.Date,
		})
	}

	return response, nil
}

func (s *driverEventService) CheckDriverInEvent(ctx context.Context, driverID, eventID uint) (bool, e.ApiError) {
	drivers, err := s.driverEventRepo.ListDriversByEvent(ctx, eventID)
	if err != nil {
		return false, err
	}

	for _, driver := range drivers {
		if driver.ID == driverID {
			return true, nil
		}
	}

	return false, nil
}