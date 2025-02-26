package repository

import (
	"context"
	"errors"
	"fmt"
	"prodes/internal/model"

	// prodes "prodes/internal/model"
	e "prodes/pkg/utils"

	"gorm.io/gorm"
)

type prodeRepository struct{
	db *gorm.DB
}

type ProdeRepository interface {
	CreateProdeCarrera(ctx context.Context, prode *model.ProdeCarrera) e.ApiError
	CreateProdeSession(ctx context.Context, prode *model.ProdeSession) e.ApiError
	GetProdeCarreraByID(ctx context.Context, prodeID int) (*model.ProdeCarrera, e.ApiError)
	GetProdeSessionByID(ctx context.Context, prodeID int) (*model.ProdeSession, e.ApiError)
	UpdateProdeCarrera(ctx context.Context, prode *model.ProdeCarrera) e.ApiError
	UpdateProdeSession(ctx context.Context, prode *model.ProdeSession) e.ApiError
	DeleteProdeCarreraByID(ctx context.Context, prodeID int, userID int) e.ApiError
    DeleteProdeSessionByID(ctx context.Context, prodeID int, userID int) e.ApiError
	// GetProdeByUserIDAndSessionID(ctx context.Context, userID int, sessionID int) (*model.ProdeCarrera, *model.ProdeSession, e.ApiError)
	GetProdeCarreraBySessionIdAndUserId(ctx context.Context, userID int, sessionID int) (*model.ProdeCarrera, e.ApiError)
	GetProdeSessionBySessionIdAndUserId(ctx context.Context, userID int, sessionID int) (*model.ProdeSession, e.ApiError)
	GetAllProdesBySessionID(ctx context.Context, sessionId int) ([]*model.ProdeCarrera, []*model.ProdeSession, e.ApiError)
	GetProdesByUserID(ctx context.Context, userID int) ([]*model.ProdeCarrera, []*model.ProdeSession, e.ApiError)
	GetProdeCarreraByUserAndSession(ctx context.Context, userID, sessionID int) (*model.ProdeCarrera, e.ApiError)
	GetProdeSessionByUserAndSession(ctx context.Context, userID, sessionID int) (*model.ProdeSession, e.ApiError)
	GetRaceProdesBySession(ctx context.Context, sessionID int) ([]*model.ProdeCarrera, e.ApiError)
	GetSessionProdesBySession(ctx context.Context, sessionID int) ([]*model.ProdeSession, e.ApiError)
}

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
	if err := r.db.WithContext(ctx).Create(prode).Error; err != nil {
		return e.NewInternalServerApiError("error creating prode session", err)
	}
	return nil
}

func (r *prodeRepository) GetProdeCarreraByID(ctx context.Context, prodeID int) (*model.ProdeCarrera, e.ApiError) {
	var prode model.ProdeCarrera

	// Usar Preload para cargar la información de la sesión relacionada
	if err := r.db.WithContext(ctx).Preload("Session").First(&prode, prodeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, e.NewNotFoundApiError("prode carrera not found")
		}
		return nil, e.NewInternalServerApiError("error finding prode carrera", err)
	}

	return &prode, nil
}

func (r *prodeRepository) GetProdeSessionByID(ctx context.Context, prodeID int) (*model.ProdeSession, e.ApiError) {
	var prode model.ProdeSession

	// Usar Preload para cargar la información de la sesión relacionada
	if err := r.db.WithContext(ctx).Preload("Session").First(&prode, prodeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, e.NewNotFoundApiError("prode session not found")
		}
		return nil, e.NewInternalServerApiError("error finding prode session", err)
	}

	return &prode, nil
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

// DeleteProdeByID elimina un pronóstico de carrera por su ID y verifica el userID
func (r *prodeRepository) DeleteProdeCarreraByID(ctx context.Context, prodeID int, userID int) e.ApiError {
	// Validar si el userID es válido (no nulo o mayor a 0)
	if userID <= 0 {
		return e.NewBadRequestApiError("Invalid userID")	
	}	

    if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", prodeID, userID).Delete(&model.ProdeCarrera{}).Error; err != nil {
        return e.NewInternalServerApiError("error deleting prode by ID", err)
    }
    return nil
}

func (r *prodeRepository) DeleteProdeSessionByID(ctx context.Context, prodeID int, userID int) e.ApiError {
	// Validar si el userID es válido (no nulo o mayor a 0)
	if userID <= 0 {
		return e.NewBadRequestApiError("Invalid userID")	
	}	

    if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", prodeID, userID).Delete(&model.ProdeSession{}).Error; err != nil {
        return e.NewInternalServerApiError("error deleting prode session by ID", err)
    }
    return nil
}

// func (r *prodeRepository) GetProdeByUserIDAndSessionID(ctx context.Context, userID int, sessionID int) (*model.ProdeCarrera, *model.ProdeSession, e.ApiError) {
// 	var prodeCarrera model.ProdeCarrera
// 	var prodeSession model.ProdeSession

// 	// Intentar obtener un ProdeCarrera
// 	errCarrera := r.db.WithContext(ctx).Preload("Session").
// 		Where("user_id = ? AND session_id = ?", userID, sessionID).
// 		First(&prodeCarrera).Error

// 	// Si hay un error y no es "record not found", es un error de base de datos
// 	if errCarrera != nil && errCarrera != gorm.ErrRecordNotFound {
// 		return nil, nil, e.NewInternalServerApiError("Error finding prode carrera", errCarrera)
// 	}

// 	// Intentar obtener un ProdeSession
// 	errSession := r.db.WithContext(ctx).Preload("Session").
// 		Where("user_id = ? AND session_id = ?", userID, sessionID).
// 		First(&prodeSession).Error

// 	// Si hay un error y no es "record not found", es un error de base de datos
// 	if errSession != nil && errSession != gorm.ErrRecordNotFound {
// 		return nil, nil, e.NewInternalServerApiError("Error finding prode session", errSession)
// 	}

// 	// Si ambos devuelven record not found, significa que no existe un prode para ese usuario y sesión
// 	if errCarrera == gorm.ErrRecordNotFound && errSession == gorm.ErrRecordNotFound {
// 		return nil, nil, e.NewNotFoundApiError("No prode found for this user and session")
// 	}

// 	// Si existe un ProdeCarrera, retornarlo
// 	if errCarrera == nil {
// 		return &prodeCarrera, nil, nil
// 	}

// 	// Si existe un ProdeSession, retornarlo
// 	return nil, &prodeSession, nil
// }

func (r *prodeRepository) GetProdeCarreraBySessionIdAndUserId(ctx context.Context, userID int, sessionID int) (*model.ProdeCarrera, e.ApiError) {
	var prodeCarrera model.ProdeCarrera

	err := r.db.WithContext(ctx).Preload("Session").
		Where("user_id = ? AND session_id = ?", userID, sessionID).
		First(&prodeCarrera).Error

	// Manejar error de "registro no encontrado"
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundApiError("No race prode found for this user and session")
		}
		return nil, e.NewInternalServerApiError("Error finding prode carrera", err)
	}

	return &prodeCarrera, nil
}

func (r *prodeRepository) GetProdeSessionBySessionIdAndUserId(ctx context.Context, userID int, sessionID int) (*model.ProdeSession, e.ApiError) {
	var prodeSession model.ProdeSession

	err := r.db.WithContext(ctx).Preload("Session").
		Where("user_id = ? AND session_id = ?", userID, sessionID).
		First(&prodeSession).Error

	// Manejar error de "registro no encontrado"
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundApiError("No session prode found for this user and session")
		}
		return nil, e.NewInternalServerApiError("Error finding prode session", err)
	}

	return &prodeSession, nil
}



func (r *prodeRepository) GetAllProdesBySessionID(ctx context.Context, sessionId int) ([]*model.ProdeCarrera, []*model.ProdeSession, e.ApiError) {
	var prodesCarrera []*model.ProdeCarrera
	var prodesSession []*model.ProdeSession

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
func (r *prodeRepository) GetProdesByUserID(ctx context.Context, userID int) ([]*model.ProdeCarrera, []*model.ProdeSession, e.ApiError) {
	var prodesCarrera []*model.ProdeCarrera
	var prodesSession []*model.ProdeSession

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

    result := r.db.WithContext(ctx).Where("user_id = ? AND session_id = ?", userID, sessionID).First(&prode)
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            fmt.Printf("No prode carrera found for userID %d and sessionID %d\n", userID, sessionID)
            return nil, nil // Devolver nil, nil para indicar no encontrado sin error HTTP
        }
        fmt.Printf("Database error for userID %d and sessionID %d: %v\n", userID, sessionID, result.Error)
        fmt.Printf("SQL query for userID %d and sessionID %d: %s\n", userID, sessionID, r.db.ToSQL(func(tx *gorm.DB) *gorm.DB {
            return tx.Where("user_id = ? AND session_id = ?", userID, sessionID).First(&model.ProdeCarrera{})
        }))
        return nil, e.NewInternalServerApiError("error finding prode carrera", result.Error)
    }

    fmt.Printf("Found prode carrera for userID %d and sessionID %d: %+v\n", userID, sessionID, prode)
    return &prode, nil
}

func (r *prodeRepository) GetProdeSessionByUserAndSession(ctx context.Context, userID, sessionID int) (*model.ProdeSession, e.ApiError) {
    var prode model.ProdeSession

    result := r.db.WithContext(ctx).Where("user_id = ? AND session_id = ?", userID, sessionID).First(&prode)
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            fmt.Printf("No prode session found for userID %d and sessionID %d\n", userID, sessionID)
            return nil, nil // Devolver nil, nil para indicar no encontrado sin error HTTP
        }
        fmt.Printf("Database error for userID %d and sessionID %d: %v\n", userID, sessionID, result.Error)
        fmt.Printf("SQL query for userID %d and sessionID %d: %s\n", userID, sessionID, r.db.ToSQL(func(tx *gorm.DB) *gorm.DB {
            return tx.Where("user_id = ? AND session_id = ?", userID, sessionID).First(&model.ProdeSession{})
        }))
        return nil, e.NewInternalServerApiError("error finding prode session", result.Error)
    }

    fmt.Printf("Found prode session for userID %d and sessionID %d: %+v\n", userID, sessionID, prode)
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