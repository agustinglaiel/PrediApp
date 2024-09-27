package repository

import (
	"context"
	"prodes/internal/model"
	prodes "prodes/internal/model"
	e "prodes/pkg/utils"

	"gorm.io/gorm"
)

type prodeRepository struct{
	db *gorm.DB
}

type ProdeRepository interface {
	CreateProdeCarrera(ctx context.Context, prode *prodes.ProdeCarrera) e.ApiError
	CreateProdeSession(ctx context.Context, prode *prodes.ProdeSession) e.ApiError
	GetProdeCarreraByID(ctx context.Context, prodeID int) (*prodes.ProdeCarrera, e.ApiError)
	GetProdeSessionByID(ctx context.Context, prodeID int) (*prodes.ProdeSession, e.ApiError)
	UpdateProdeCarrera(ctx context.Context, prode *prodes.ProdeCarrera) e.ApiError
	UpdateProdeSession(ctx context.Context, prode *prodes.ProdeSession) e.ApiError
	DeleteProdeCarreraByID(ctx context.Context, prodeID int, userID int) e.ApiError
    DeleteProdeSessionByID(ctx context.Context, prodeID int, userID int) e.ApiError
	GetProdesByUserIDAndSessionID(ctx context.Context, userID int, sessionId int) ([]*prodes.ProdeCarrera, []*prodes.ProdeSession, e.ApiError)
	GetAllProdesBySessionID(ctx context.Context, sessionId int) ([]*prodes.ProdeCarrera, []*prodes.ProdeSession, e.ApiError)
	GetProdesByUserID(ctx context.Context, userID int) ([]*prodes.ProdeCarrera, []*prodes.ProdeSession, e.ApiError)
	GetProdeCarreraByUserAndSession(ctx context.Context, userID, sessionID int) (*model.ProdeCarrera, e.ApiError)
	GetProdeSessionByUserAndSession(ctx context.Context, userID, sessionID int) (*model.ProdeSession, e.ApiError)
	GetRaceProdesBySession(ctx context.Context, sessionID int) ([]*model.ProdeCarrera, e.ApiError)
	GetSessionProdesBySession(ctx context.Context, sessionID int) ([]*model.ProdeSession, e.ApiError)
}

func NewProdeRepository(db *gorm.DB) ProdeRepository {
	return &prodeRepository{db: db}
}

func (r *prodeRepository) CreateProdeCarrera(ctx context.Context, prode *prodes.ProdeCarrera) e.ApiError {
	if err := r.db.WithContext(ctx).Create(prode).Error; err != nil {
		return e.NewInternalServerApiError("error creating prode carrera", err)
	}
	return nil
}

func (r *prodeRepository) CreateProdeSession(ctx context.Context, prode *prodes.ProdeSession) e.ApiError {
	if err := r.db.WithContext(ctx).Create(prode).Error; err != nil {
		return e.NewInternalServerApiError("error creating prode session", err)
	}
	return nil
}

func (r *prodeRepository) GetProdeCarreraByID(ctx context.Context, prodeID int) (*prodes.ProdeCarrera, e.ApiError) {
	var prode prodes.ProdeCarrera

	// Usar Preload para cargar la información de la sesión relacionada
	if err := r.db.WithContext(ctx).Preload("Session").First(&prode, prodeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, e.NewNotFoundApiError("prode carrera not found")
		}
		return nil, e.NewInternalServerApiError("error finding prode carrera", err)
	}

	return &prode, nil
}

func (r *prodeRepository) GetProdeSessionByID(ctx context.Context, prodeID int) (*prodes.ProdeSession, e.ApiError) {
	var prode prodes.ProdeSession

	// Usar Preload para cargar la información de la sesión relacionada
	if err := r.db.WithContext(ctx).Preload("Session").First(&prode, prodeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, e.NewNotFoundApiError("prode session not found")
		}
		return nil, e.NewInternalServerApiError("error finding prode session", err)
	}

	return &prode, nil
}

func (r *prodeRepository) UpdateProdeCarrera(ctx context.Context, prode *prodes.ProdeCarrera) e.ApiError {
	if err := r.db.WithContext(ctx).Save(prode).Error; err != nil {
		return e.NewInternalServerApiError("error updating prode carrera", err)
	}
	return nil
}

func (r *prodeRepository) UpdateProdeSession(ctx context.Context, prode *prodes.ProdeSession) e.ApiError {
	if err := r.db.WithContext(ctx).Save(prode).Error; err != nil {
		return e.NewInternalServerApiError("error updating prode session", err)
	}
	return nil
}

// DeleteProdeByID elimina un pronóstico de carrera por su ID y verifica el userID
func (r *prodeRepository) DeleteProdeCarreraByID(ctx context.Context, prodeID int, userID int) e.ApiError {
	// Validar si el userID es válido (no nulo o mayor a 0)
	if userID <= 0 {
		return e.NewBadRequestApiError("Invalid userID")	
	}	

    if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", prodeID, userID).Delete(&prodes.ProdeCarrera{}).Error; err != nil {
        return e.NewInternalServerApiError("error deleting prode by ID", err)
    }
    return nil
}

func (r *prodeRepository) DeleteProdeSessionByID(ctx context.Context, prodeID int, userID int) e.ApiError {
	// Validar si el userID es válido (no nulo o mayor a 0)
	if userID <= 0 {
		return e.NewBadRequestApiError("Invalid userID")	
	}	

    if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", prodeID, userID).Delete(&prodes.ProdeSession{}).Error; err != nil {
        return e.NewInternalServerApiError("error deleting prode session by ID", err)
    }
    return nil
}

func (r *prodeRepository) GetProdesByUserIDAndSessionID(ctx context.Context, userID int, sessionId int) ([]*prodes.ProdeCarrera, []*prodes.ProdeSession, e.ApiError) {
	var prodesCarrera []*prodes.ProdeCarrera
	var prodesSession []*prodes.ProdeSession

	// Usar Preload para cargar la información de la sesión relacionada
	if err := r.db.WithContext(ctx).Preload("Session").Where("user_id = ? AND session_id = ?", userID, sessionId).Find(&prodesCarrera).Error; err != nil {
		return nil, nil, e.NewInternalServerApiError("error finding prodes carrera", err)
	}

	if err := r.db.WithContext(ctx).Preload("Session").Where("user_id = ? AND session_id = ?", userID, sessionId).Find(&prodesSession).Error; err != nil {
		return nil, nil, e.NewInternalServerApiError("error finding prodes session", err)
	}

	return prodesCarrera, prodesSession, nil
}

func (r *prodeRepository) GetAllProdesBySessionID(ctx context.Context, sessionId int) ([]*prodes.ProdeCarrera, []*prodes.ProdeSession, e.ApiError) {
	var prodesCarrera []*prodes.ProdeCarrera
	var prodesSession []*prodes.ProdeSession

	// Usar Preload para cargar la información de la sesión relacionada
	if err := r.db.WithContext(ctx).Preload("Session").Where("session_id = ?", sessionId).Find(&prodesCarrera).Error; err != nil {
		return nil, nil, e.NewInternalServerApiError("error finding prodes carrera", err)
	}

	if err := r.db.WithContext(ctx).Preload("Session").Where("session_id = ?", sessionId).Find(&prodesSession).Error; err != nil {
		return nil, nil, e.NewInternalServerApiError("error finding prodes session", err)
	}

	return prodesCarrera, prodesSession, nil
}

// GetProdesByUserID obtiene todos los prodes (carrera y sesión) realizados por un usuario
func (r *prodeRepository) GetProdesByUserID(ctx context.Context, userID int) ([]*prodes.ProdeCarrera, []*prodes.ProdeSession, e.ApiError) {
	var prodesCarrera []*prodes.ProdeCarrera
	var prodesSession []*prodes.ProdeSession

	// Usar Preload para cargar la información de la sesión relacionada
	if err := r.db.WithContext(ctx).Preload("Session").Where("user_id = ?", userID).Find(&prodesCarrera).Error; err != nil {
		return nil, nil, e.NewInternalServerApiError("error finding race predictions", err)
	}

	if err := r.db.WithContext(ctx).Preload("Session").Where("user_id = ?", userID).Find(&prodesSession).Error; err != nil {
		return nil, nil, e.NewInternalServerApiError("error finding session predictions", err)
	}

	return prodesCarrera, prodesSession, nil
}

func (r *prodeRepository) GetProdeCarreraByUserAndSession(ctx context.Context, userID, sessionID int) (*model.ProdeCarrera, e.ApiError) {
    var prode model.ProdeCarrera

    // Buscar un único prode para el usuario y la sesión dados
    if err := r.db.WithContext(ctx).Where("user_id = ? AND session_id = ?", userID, sessionID).First(&prode).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, e.NewNotFoundApiError("prode carrera not found")
        }
        return nil, e.NewInternalServerApiError("error finding prode carrera", err)
    }

    return &prode, nil
}

func (r *prodeRepository) GetProdeSessionByUserAndSession(ctx context.Context, userID, sessionID int) (*model.ProdeSession, e.ApiError) {
    var prode model.ProdeSession

    // Buscar un único prode de sesión para el usuario y la sesión dados
    if err := r.db.WithContext(ctx).Where("user_id = ? AND session_id = ?", userID, sessionID).First(&prode).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, e.NewNotFoundApiError("prode session not found")
        }
        return nil, e.NewInternalServerApiError("error finding prode session", err)
    }

    return &prode, nil
}

func (r *prodeRepository) GetRaceProdesBySession(ctx context.Context, sessionID int) ([]*model.ProdeCarrera, e.ApiError) {
    var raceProdes []*model.ProdeCarrera

    // Usar Preload para cargar la información de la sesión relacionada
    if err := r.db.WithContext(ctx).Preload("Session").Where("session_id = ?", sessionID).Find(&raceProdes).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, e.NewNotFoundApiError("No se encontraron pronósticos de carrera para la sesión")
        }
        return nil, e.NewInternalServerApiError("Error fetching race prodes for session", err)
    }

    return raceProdes, nil
}

func (r *prodeRepository) GetSessionProdesBySession(ctx context.Context, sessionID int) ([]*model.ProdeSession, e.ApiError) {
    var prodesSession []*model.ProdeSession

    // Usar Preload para cargar la información de la sesión relacionada
    if err := r.db.WithContext(ctx).Preload("Session").Where("session_id = ?", sessionID).Find(&prodesSession).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, e.NewNotFoundApiError("No se encontraron pronósticos de sesión para la sesión")
        }
        return nil, e.NewInternalServerApiError("Error fetching session prodes for session", err)
    }

    return prodesSession, nil
}