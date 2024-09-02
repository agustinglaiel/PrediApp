package repository

import (
	"admin/internal/model/drivers"
	e "admin/pkg/utils"
	"context"

	"gorm.io/gorm"
)

type driverEventRepository struct {
	db *gorm.DB
}

type DriverEventRepository interface {
	AddDriverToEvent(ctx context.Context, driverEvent *drivers.DriverEvent) e.ApiError
	RemoveDriverFromEvent(ctx context.Context, driverEventID uint) e.ApiError
	ListDriversByEvent(ctx context.Context, eventID uint) ([]*drivers.Driver, e.ApiError)
}

func NewDriverEventRepository(db *gorm.DB) DriverEventRepository {
	return &driverEventRepository{db: db}
}

func (r *driverEventRepository) AddDriverToEvent(ctx context.Context, driverEvent *drivers.DriverEvent) e.ApiError {
	if err := r.db.WithContext(ctx).Create(driverEvent).Error; err != nil {
		return e.NewInternalServerApiError("Error añadiendo piloto al evento", err)
	}
	return nil
}

func (r *driverEventRepository) RemoveDriverFromEvent(ctx context.Context, driverEventID uint) e.ApiError {
	if err := r.db.WithContext(ctx).Where("id = ?", driverEventID).Delete(&drivers.DriverEvent{}).Error; err != nil {
		return e.NewInternalServerApiError("Error eliminando piloto del evento", err)
	}
	return nil
}

func (r *driverEventRepository) ListDriversByEvent(ctx context.Context, eventID uint) ([]*drivers.Driver, e.ApiError) {
	var drivers []*drivers.Driver
	if err := r.db.WithContext(ctx).Joins("JOIN drivers_events ON drivers_events.driver_id = drivers.id").Where("drivers_events.event_id = ?", eventID).Find(&drivers).Error; err != nil {
		return nil, e.NewInternalServerApiError("Error listando pilotos por evento", err)
	}
	return drivers, nil
}