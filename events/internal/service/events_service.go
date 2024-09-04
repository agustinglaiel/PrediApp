package service

/*
import (
	"admin/internal/dto/events"
	"admin/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
)

type eventService struct {
	eventRepo EventRepository
	httpClient *utils.HttpClient // Inyectamos el cliente HTTP
}

type EventService interface {
	GetSessionById(ctx context.Context, sessionId int) (events.SessionDTO, error)
	// Otros métodos del servicio
}

func NewEventService(eventRepo EventRepository, httpClient *utils.HttpClient) EventService {
	return &eventService{
		eventRepo: eventRepo,
		httpClient: httpClient,
	}
}

func (s *eventService) GetSessionById(ctx context.Context, sessionId int) (events.SessionDTO, error) {
	endpoint := fmt.Sprintf("/sessions/%d", sessionId)
	responseData, err := s.httpClient.Get(endpoint)
	if err != nil {
		return events.SessionDTO{}, err
	}

	var sessionDTO events.SessionDTO
	err = json.Unmarshal(responseData, &sessionDTO)
	if err != nil {
		return events.SessionDTO{}, err
	}

	return sessionDTO, nil
}
*/