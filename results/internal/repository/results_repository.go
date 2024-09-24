package repository

import (
	"context"
	"results/internal/model"
	e "results/pkg/utils"

	"gorm.io/gorm"
)

type resultRepository struct {
	db *gorm.DB
}

type ResultRepository interface {
	CreateResult(ctx context.Context, result *model.Result) e.ApiError
	GetResultByID(ctx context.Context, resultID uint) (*model.Result, e.ApiError)
	UpdateResult(ctx context.Context, result *model.Result) e.ApiError
	DeleteResult(ctx context.Context, resultID uint) e.ApiError
	GetResultsBySessionID(ctx context.Context, sessionID uint) ([]*model.Result, e.ApiError)
	GetResultsByDriverID(ctx context.Context, driverID uint) ([]*model.Result, e.ApiError)
	GetAllResults(ctx context.Context) ([]*model.Result, e.ApiError)
	GetFastestLapInSession(ctx context.Context, sessionID uint) (*model.Result, e.ApiError)
	GetDriverPositionInSession(ctx context.Context, driverID uint, sessionID uint) (int, e.ApiError)
	GetResultsOrderedByPosition(ctx context.Context, sessionID uint) ([]*model.Result, e.ApiError)
	ExistsSessionInResults(ctx context.Context, sessionID uint) (bool, e.ApiError)
}

func NewResultRepository(db *gorm.DB) ResultRepository {
	return &resultRepository{db: db}
}

// CreateResult crea un nuevo resultado en la base de datos
func (r *resultRepository) CreateResult(ctx context.Context, result *model.Result) e.ApiError {
	if err := r.db.WithContext(ctx).Create(result).Error; err != nil {
		return e.NewInternalServerApiError("Error creating result", err)
	}
	return nil
}

// GetResultByID obtiene un resultado específico por su ID
func (r *resultRepository) GetResultByID(ctx context.Context, resultID uint) (*model.Result, e.ApiError) {
	var result model.Result
	if err := r.db.WithContext(ctx).Preload("Driver").Preload("Session").First(&result, resultID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, e.NewNotFoundApiError("Result not found")
		}
		return nil, e.NewInternalServerApiError("Error finding result", err)
	}
	return &result, nil
}

// UpdateResult actualiza un resultado existente en la base de datos
func (r *resultRepository) UpdateResult(ctx context.Context, result *model.Result) e.ApiError {
	if err := r.db.WithContext(ctx).Save(result).Error; err != nil {
		return e.NewInternalServerApiError("Error updating result", err)
	}
	return nil
}

// DeleteResult elimina un resultado de la base de datos por su ID
func (r *resultRepository) DeleteResult(ctx context.Context, resultID uint) e.ApiError {
	if err := r.db.WithContext(ctx).Delete(&model.Result{}, resultID).Error; err != nil {
		return e.NewInternalServerApiError("Error deleting result", err)
	}
	return nil
}

// GetResultsBySessionID obtiene todos los resultados de una sesión específica
func (r *resultRepository) GetResultsBySessionID(ctx context.Context, sessionID uint) ([]*model.Result, e.ApiError) {
	var results []*model.Result
	if err := r.db.WithContext(ctx).Preload("Driver").Preload("Session").Where("session_id = ?", sessionID).Find(&results).Error; err != nil {
		return nil, e.NewInternalServerApiError("Error finding results by session ID", err)
	}
	return results, nil
}

// GetResultsByDriverID obtiene todos los resultados de un piloto específico
func (r *resultRepository) GetResultsByDriverID(ctx context.Context, driverID uint) ([]*model.Result, e.ApiError) {
	var results []*model.Result
	if err := r.db.WithContext(ctx).Preload("Driver").Preload("Session").Where("driver_id = ?", driverID).Find(&results).Error; err != nil {
		return nil, e.NewInternalServerApiError("Error finding results by driver ID", err)
	}
	return results, nil
}

// GetAllResults obtiene todos los resultados de la base de datos
func (r *resultRepository) GetAllResults(ctx context.Context) ([]*model.Result, e.ApiError) {
	var results []*model.Result
	if err := r.db.WithContext(ctx).Preload("Driver").Preload("Session").Find(&results).Error; err != nil {
		return nil, e.NewInternalServerApiError("Error finding all results", err)
	}
	return results, nil
}

// GetFastestLapInSession obtiene el piloto con el tiempo de vuelta más rápido en una sesión específica
func (r *resultRepository) GetFastestLapInSession(ctx context.Context, sessionID uint) (*model.Result, e.ApiError) {
    var result model.Result
    if err := r.db.WithContext(ctx).Preload("Driver").Preload("Session").Where("session_id = ?", sessionID).Order("fastest_lap_time ASC").First(&result).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, e.NewNotFoundApiError("No fastest lap found for the session")
        }
        return nil, e.NewInternalServerApiError("Error fetching fastest lap for session", err)
    }
    return &result, nil
}

// GetDriverPositionInSession obtiene la posición final de un piloto en una sesión específica
func (r *resultRepository) GetDriverPositionInSession(ctx context.Context, driverID uint, sessionID uint) (int, e.ApiError) {
    var result model.Result
    // Buscar el resultado del piloto en la sesión específica
    if err := r.db.WithContext(ctx).Where("driver_id = ? AND session_id = ?", driverID, sessionID).Select("position").First(&result).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return 0, e.NewNotFoundApiError("No position found for the given driver in the session")
        }
        return 0, e.NewInternalServerApiError("Error fetching driver position in session", err)
    }
    // Devolver únicamente la posición
    return result.Position, nil
}

// GetResultsOrderedByPosition obtiene los resultados de una sesión ordenados por posición
/*
obtiene todos los resultados de una sesión específica (incluyendo detalles de piloto y sesión) y los ordena por posición.
*/
func (r *resultRepository) GetResultsOrderedByPosition(ctx context.Context, sessionID uint) ([]*model.Result, e.ApiError) {
    var results []*model.Result
    if err := r.db.WithContext(ctx).Preload("Driver").Preload("Session").Where("session_id = ?", sessionID).Order("position ASC").Find(&results).Error; err != nil {
        return nil, e.NewInternalServerApiError("Error fetching ordered results for session", err)
    }
    return results, nil
}

// ExistsSessionInResults verifica si existen resultados para un sessionID dado
func (r *resultRepository) ExistsSessionInResults(ctx context.Context, sessionID uint) (bool, e.ApiError) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Result{}).Where("session_id = ?", sessionID).Count(&count).Error; err != nil {
		return false, e.NewInternalServerApiError("Error checking session in results", err)
	}

	return count > 0, nil
}