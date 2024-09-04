package service

import (
	"context"
	"encoding/json"
	client "events/internal/client"
	dto "events/internal/dto"
	"events/internal/model"
	event "events/internal/repository"
	e "events/pkg/utils"
	"fmt"
)

type eventService struct {
	eventRepo event.EventRepository
	httpClient *client.HttpClient // Inyectamos el cliente HTTP
}

type EventService interface {
	GetSessionById(ctx context.Context, sessionId int) (dto.SessionDTO, e.ApiError)
	CreateEvent(ctx context.Context, request dto.CreateEventDTO) (dto.EventResponseDTO, e.ApiError)
	UpdateEvent(ctx context.Context, eventID int, request dto.UpdateEventDTO) (dto.EventResponseDTO, e.ApiError)
	DeleteEvent(ctx context.Context, eventID int) e.ApiError
	ListEvents(ctx context.Context) ([]dto.ListEventDTO, e.ApiError)
}

func NewEventService(eventRepo event.EventRepository, httpClient *client.HttpClient) EventService {
	return &eventService{
		eventRepo: eventRepo,
		httpClient: httpClient,
	}
}

func (s *eventService) GetSessionById(ctx context.Context, sessionId int) (dto.SessionDTO, e.ApiError) {
	endpoint := fmt.Sprintf("/sessions/%d", sessionId)
	responseData, err := s.httpClient.Get(endpoint)
	if err != nil {
		return dto.SessionDTO{}, e.NewInternalServerApiError("Error al obtener la sesión", err)
	}

	var sessionDTO dto.SessionDTO
	err = json.Unmarshal(responseData, &sessionDTO)
	if err != nil {
		return dto.SessionDTO{}, e.NewInternalServerApiError("Error al deserializar los datos de la sesión", err)
	}

	return sessionDTO, nil
}

func (s *eventService) CreateEvent(ctx context.Context, request dto.CreateEventDTO) (dto.EventResponseDTO, e.ApiError) {
	// Lógica para crear un nuevo evento
	event := model.Event{
		SessionID:           request.SessionID,
		Date:                request.Date,
	}

	// Guardar el evento en el repositorio
	if err := s.eventRepo.CreateEvent(ctx, &event); err != nil {
		return dto.EventResponseDTO{}, err
	}

	// Crear registros relacionados automáticamente
	if err := s.createRelatedEntities(ctx, event.ID); err != nil {
		return dto.EventResponseDTO{}, err
	}

	// Convertir a DTO para la respuesta
	response := dto.EventResponseDTO{
		ID:      event.ID,
		Session: dto.SessionDTO{ID: event.SessionID},
		Date:    event.Date,
	}

	return response, nil
}

// Función para crear automáticamente las entidades relacionadas
func (s *eventService) createRelatedEntities(ctx context.Context, eventID int) e.ApiError {
	// Crear un registro vacío en race_result
	raceResult := model.RaceResult{EventID: eventID}
	if err := s.eventRepo.CreateRaceResult(ctx, &raceResult); err != nil {
		return err
	}

	// Crear un registro vacío en sprint_qualy_results
	sprintQualyResult := model.SprintQualyResult{EventID: eventID}
	if err := s.eventRepo.CreateSprintQualyResult(ctx, &sprintQualyResult); err != nil {
		return err
	}

	// Crear un registro vacío en fp_results
	fpResult := model.FPResult{EventID: eventID}
	if err := s.eventRepo.CreateFPResult(ctx, &fpResult); err != nil {
		return err
	}

	// Crear un registro vacío en drivers_event
	driversEvent := model.DriversEvent{EventID: eventID}
	if err := s.eventRepo.CreateDriversEvent(ctx, &driversEvent); err != nil {
		return err
	}

	return nil
}


func (s *eventService) UpdateEvent(ctx context.Context, eventID int, request dto.UpdateEventDTO) (dto.EventResponseDTO, e.ApiError) {
	// Buscar el evento por ID
	event, err := s.eventRepo.GetEventByID(ctx, eventID)
	if err != nil {
		return dto.EventResponseDTO{}, err
	}

	// Actualizar solo los campos proporcionados
	if request.Date != nil {
		event.Date = *request.Date
	}
	if request.RaceResultID != nil {
		event.RaceResultID = request.RaceResultID
	}
	if request.SprintRaceResultID != nil {
		event.SprintRaceResultID = request.SprintRaceResultID
	}
	if request.QualyResultID != nil {
		event.QualyResultID = request.QualyResultID
	}
	// Actualizar otros campos...

	// Guardar los cambios
	if err := s.eventRepo.UpdateEvent(ctx, event); err != nil {
		return dto.EventResponseDTO{}, err
	}

	// Convertir a DTO para la respuesta
	response := dto.EventResponseDTO{
		ID:          event.ID,
		Session:     dto.SessionDTO{ID: event.SessionID},
		Date:        event.Date,
		// Otros campos actualizados...
	}

	return response, nil
}

func (s *eventService) DeleteEvent(ctx context.Context, eventID int) e.ApiError {
	// Intentar eliminar el evento
	if err := s.eventRepo.DeleteEventById(ctx, eventID); err != nil {
		return err
	}

	return nil
}

func (s *eventService) ListEvents(ctx context.Context) ([]dto.ListEventDTO, e.ApiError) {
	// Obtener todos los eventos del repositorio
	events, err := s.eventRepo.ListAllEvents(ctx)
	if err != nil {
		return nil, err
	}

	// Convertir los eventos a DTO
	var response []dto.ListEventDTO
	for _, event := range events {
		response = append(response, dto.ListEventDTO{
			ID:       event.ID,
			Session:  fmt.Sprintf("Session %d", event.SessionID),  // Podemos hacer esta lógica más compleja
			Date:     event.Date,
		})
	}

	return response, nil
}