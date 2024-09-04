package repository

import (
	"context"
	"events/internal/model"
	e "events/pkg/utils"

	"gorm.io/gorm"
)

type eventRepository struct {
    db *gorm.DB
}

type EventRepository interface {
	CreateEvent(ctx context.Context, event *model.Event) e.ApiError
	CreateRaceResult(ctx context.Context, raceResult *model.RaceResult) e.ApiError
	CreateSprintQualyResult(ctx context.Context, sprintQualyResult *model.SprintQualyResult) e.ApiError
	CreateFPResult(ctx context.Context, fpResult *model.FPResult) e.ApiError
	CreateDriversEvent(ctx context.Context, driversEvent *model.DriversEvent) e.ApiError
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

// CreateEvent inserta un nuevo evento en la base de datos
func (r *eventRepository) CreateEvent(ctx context.Context, event *model.Event) e.ApiError {
	// Usar GORM para insertar el evento en la base de datos
	if err := r.db.WithContext(ctx).Create(event).Error; err != nil {
		return e.NewInternalServerApiError("Error al crear el evento", err)
	}
	return nil
}

func (r *eventRepository) CreateRaceResult(ctx context.Context, raceResult *model.RaceResult) e.ApiError {
	if err := r.db.WithContext(ctx).Create(raceResult).Error; err != nil {
		return e.NewInternalServerApiError("Error creating race result", err)
	}
	return nil
}

func (r *eventRepository) CreateSprintQualyResult(ctx context.Context, sprintQualyResult *model.SprintQualyResult) e.ApiError {
	if err := r.db.WithContext(ctx).Create(sprintQualyResult).Error; err != nil {
		return e.NewInternalServerApiError("Error creating sprint qualy result", err)
	}
	return nil
}

func (r *eventRepository) CreateFPResult(ctx context.Context, fpResult *model.FPResult) e.ApiError {
	if err := r.db.WithContext(ctx).Create(fpResult).Error; err != nil {
		return e.NewInternalServerApiError("Error creating fp result", err)
	}
	return nil
}

func (r *eventRepository) CreateDriversEvent(ctx context.Context, driversEvent *model.DriversEvent) e.ApiError {
	if err := r.db.WithContext(ctx).Create(driversEvent).Error; err != nil {
		return e.NewInternalServerApiError("Error creating drivers event", err)
	}
	return nil
}
