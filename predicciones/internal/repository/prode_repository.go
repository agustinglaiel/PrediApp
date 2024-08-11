package repository

import (
	"context"
	"predicciones/internal/model"
	e "predicciones/pkg/utils"

	"gorm.io/gorm"
)

// prodeRepository es una estructura que implementa la interfaz ProdeRepository
type prodeRepository struct {
	db *gorm.DB
}

// ProdeRepository define los métodos que deben ser implementados por el repositorio de predicciones
type ProdeRepository interface {
	CreateProdeCarrera(ctx context.Context, prode *model.ProdeCarrera) e.ApiError
	CreateProdeSession(ctx context.Context, prode *model.ProdeSession) e.ApiError
	GetProdeCarreraByID(ctx context.Context, id uint) (*model.ProdeCarrera, e.ApiError)
	GetProdeSessionByID(ctx context.Context, id uint) (*model.ProdeSession, e.ApiError)
	GetProdesByUserID(ctx context.Context, userID uint) ([]model.ProdeCarrera, []model.ProdeSession, e.ApiError) // Lista de predicciones de carrera o sesión
	GetProdesByEventID(ctx context.Context, eventID uint) ([]model.ProdeCarrera, []model.ProdeSession, e.ApiError) // Lista de predicciones de carrera o sesión
	UpdateProdeCarrera(ctx context.Context, prode *model.ProdeCarrera) e.ApiError
	UpdateProdeSession(ctx context.Context, prode *model.ProdeSession) e.ApiError
	DeleteProdeCarreraByID(ctx context.Context, id uint) e.ApiError
	DeleteProdeSessionByID(ctx context.Context, id uint) e.ApiError
}

// NewProdeRepository crea una nueva instancia de prodeRepository
func NewProdeRepository(db *gorm.DB) ProdeRepository {
	return &prodeRepository{db: db}
}

func (r *prodeRepository) CreateProdeCarrera(ctx context.Context, prode *model.ProdeCarrera) e.ApiError {
	if err := r.db.WithContext(ctx).Create(prode).Error; err != nil {
		return e.NewInternalServerApiError("error creating prode carrera", err)
	}
	return nil
}

func (r *prodeRepository) CreateProdeSession(ctx context.Context, prode *model.ProdeSession) e.ApiError {
    if err := r.db.Create(prode).Error; err != nil {
        return e.NewInternalServerApiError("error creating prode session", err)
    }
    return nil
}

func (r *prodeRepository) GetProdeCarreraByID(ctx context.Context, id uint) (*model.ProdeCarrera, e.ApiError) {
	var prode model.ProdeCarrera
	if err := r.db.WithContext(ctx).First(&prode, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, e.NewNotFoundApiError("prode carrera not found")
		}
		return nil, e.NewInternalServerApiError("error finding prode carrera", err)
	}
	return &prode, nil
}

func (r *prodeRepository) GetProdeSessionByID(ctx context.Context, id uint) (*model.ProdeSession, e.ApiError) {
	var prode model.ProdeSession
	if err := r.db.WithContext(ctx).First(&prode, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, e.NewNotFoundApiError("prode session not found")
		}
		return nil, e.NewInternalServerApiError("error finding prode session", err)
	}
	return &prode, nil
}

func (r *prodeRepository) GetProdesByUserID(ctx context.Context, userID uint) ([]model.ProdeCarrera, []model.ProdeSession, e.ApiError) {
	var prodesCarrera []model.ProdeCarrera
	var prodesSession []model.ProdeSession
	
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&prodesCarrera).Error; err != nil {
		return nil, nil, e.NewInternalServerApiError("error finding prode carrera", err)
	}
	
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&prodesSession).Error; err != nil {
		return nil, nil, e.NewInternalServerApiError("error finding prode session", err)
	}

	return prodesCarrera, prodesSession, nil
}

func (r *prodeRepository) GetProdesByEventID(ctx context.Context, eventID uint) ([]model.ProdeCarrera, []model.ProdeSession, e.ApiError) {
	var prodesCarrera []model.ProdeCarrera
	var prodesSession []model.ProdeSession
	
	if err := r.db.WithContext(ctx).Where("event_id = ?", eventID).Find(&prodesCarrera).Error; err != nil {
		return nil, nil, e.NewInternalServerApiError("error finding prode carrera", err)
	}
	
	if err := r.db.WithContext(ctx).Where("event_id = ?", eventID).Find(&prodesSession).Error; err != nil {
		return nil, nil, e.NewInternalServerApiError("error finding prode session", err)
	}

	return prodesCarrera, prodesSession, nil
}


func (r *prodeRepository) UpdateProdeCarrera(ctx context.Context, prode *model.ProdeCarrera) e.ApiError {
	if err := r.db.WithContext(ctx).Save(prode).Error; err != nil {
		return e.NewInternalServerApiError("error updating prode carrera", err)
	}
	return nil
}

func (r *prodeRepository) UpdateProdeSession(ctx context.Context, prode *model.ProdeSession) e.ApiError {
	if err := r.db.WithContext(ctx).Save(prode).Error; err != nil {
		return e.NewInternalServerApiError("error updating prode session", err)
	}
	return nil
}

func (r *prodeRepository) DeleteProdeCarreraByID(ctx context.Context, id uint) e.ApiError {
	if err := r.db.WithContext(ctx).Delete(&model.ProdeCarrera{}, id).Error; err != nil {
		return e.NewInternalServerApiError("error deleting prode carrera", err)
	}
	return nil
}

func (r *prodeRepository) DeleteProdeSessionByID(ctx context.Context, id uint) e.ApiError {
	if err := r.db.WithContext(ctx).Delete(&model.ProdeSession{}, id).Error; err != nil {
		return e.NewInternalServerApiError("error deleting prode session", err)
	}
	return nil
}