package repository

import (
	model "admin/internal/model/sessions"
	e "admin/pkg/utils"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type sessionRepository struct {
	db *gorm.DB
}

type SessionRepository interface{
	CreateSession(ctx context.Context, session *model.Session) e.ApiError
	GetSessionById(ctx context.Context, sessionID uint)(*model.Session, e.ApiError)
	UpdateSessionById(ctx context.Context, session *model.Session) e.ApiError
	DeleteSessionById(ctx context.Context, sessionID uint) e.ApiError
	GetSessionByYear(ctx context.Context, year int) ([]*model.Session, e.ApiError)
	GetSessionNameAndTypeBySessionID(ctx context.Context, sessionID uint) (string, string, e.ApiError)
	GetSessionsByCircuitKey(ctx context.Context, circuitKey int) ([]*model.Session, e.ApiError)
	GetSessionsByCountryCode(ctx context.Context, countryCode string) ([]*model.Session, e.ApiError)
	GetUpcomingSessions(ctx context.Context) ([]*model.Session, e.ApiError)
	GetSessionsBetweenDates(ctx context.Context, startDate, endDate time.Time) ([]*model.Session, e.ApiError)
	GetSessionsByNameAndType(ctx context.Context, sessionName, sessionType string) ([]*model.Session, e.ApiError)
	GetSessionBySessionKey(ctx context.Context, sessionKey int) (*model.Session, e.ApiError)
	GetAllSessions(ctx context.Context) ([]*model.Session, e.ApiError)
	GetSessionsByCircuitKeyAndYear(ctx context.Context, circuitKey, year int) ([]*model.Session, e.ApiError)
}

func NewSessionRepository(db *gorm.DB) SessionRepository{
	return &sessionRepository{db: db}
}

func (s *sessionRepository) CreateSession(ctx context.Context, session *model.Session) e.ApiError{
	if err := s.db.WithContext(ctx).Create(session).Error; err != nil {
        return e.NewInternalServerApiError("Error creando sesión", err)
    }

    // Log para verificar el ID asignado
    fmt.Printf("Session created with ID: %d\n", session.ID)
    
    return nil
}

func (s *sessionRepository) GetSessionById(ctx context.Context, sessionID uint)(*model.Session, e.ApiError){
	var session model.Session
	if err := s.db.WithContext(ctx).First(&session, sessionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound{
			return nil, e.NewNotFoundApiError("Sesion no encontrada")
		}
		return nil, e.NewInternalServerApiError("Error encontrando sesión", err)
	}
	return &session, nil
}

func (s *sessionRepository) UpdateSessionById(ctx context.Context, session *model.Session) e.ApiError {
	if err := s.db.WithContext(ctx).Save(session).Error; err !=nil{
		return e.NewInternalServerApiError("Error actualizando la sesión", err)
	}
	return nil
}

func (s *sessionRepository) DeleteSessionById(ctx context.Context, sessionID uint) e.ApiError {
    // Eliminar físicamente la sesión utilizando el ID
    if err := s.db.WithContext(ctx).Unscoped().Where("id = ?", sessionID).Delete(&model.Session{}).Error; err != nil {
        return e.NewInternalServerApiError("Error eliminando la sesión", err)
    }
    return nil
}

func (s *sessionRepository) GetSessionByYear(ctx context.Context, year int) ([]*model.Session, e.ApiError) {
	var sessions []*model.Session
	if err := s.db.WithContext(ctx).Where("year = ?", year).Find(&sessions).Error; err != nil {
		return nil, e.NewInternalServerApiError("Error encontrando sesiones por año", err)
	}
	return sessions, nil
}

func (s *sessionRepository) GetSessionNameAndTypeBySessionID(ctx context.Context, sessionID uint) (string, string, e.ApiError) {
	var session model.Session
	if err := s.db.WithContext(ctx).Select("session_name", "session_type").First(&session, sessionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", "", e.NewNotFoundApiError("Sesión no encontrada")
		}
		return "", "", e.NewInternalServerApiError("Error obteniendo nombre y tipo de la sesión", err)
	}
	return session.SessionName, session.SessionType, nil
}

func (s *sessionRepository) GetSessionsByCircuitKey(ctx context.Context, circuitKey int) ([]*model.Session, e.ApiError) {
    var sessions []*model.Session
    if err := s.db.WithContext(ctx).Where("circuit_key = ?", circuitKey).Find(&sessions).Error; err != nil {
        return nil, e.NewInternalServerApiError("Error encontrando sesiones por circuito", err)
    }
    return sessions, nil
}

func (s *sessionRepository) GetSessionsByCountryCode(ctx context.Context, countryCode string) ([]*model.Session, e.ApiError) {
    var sessions []*model.Session
    if err := s.db.WithContext(ctx).Where("country_code = ?", countryCode).Find(&sessions).Error; err != nil {
        return nil, e.NewInternalServerApiError("Error encontrando sesiones por país", err)
    }
    return sessions, nil
}

func (s *sessionRepository) GetUpcomingSessions(ctx context.Context) ([]*model.Session, e.ApiError) {
    var sessions []*model.Session
    currentTime := time.Now()
    if err := s.db.WithContext(ctx).Where("date_start > ?", currentTime).Find(&sessions).Error; err != nil {
        return nil, e.NewInternalServerApiError("Error encontrando próximas sesiones", err)
    }
    return sessions, nil
}

func (s *sessionRepository) GetSessionsBetweenDates(ctx context.Context, startDate, endDate time.Time) ([]*model.Session, e.ApiError) {
    var sessions []*model.Session
    if err := s.db.WithContext(ctx).Where("date_start >= ? AND date_end <= ?", startDate, endDate).Find(&sessions).Error; err != nil {
        return nil, e.NewInternalServerApiError("Error encontrando sesiones entre fechas", err)
    }
    return sessions, nil
}

func (s *sessionRepository) GetSessionsByNameAndType(ctx context.Context, sessionName, sessionType string) ([]*model.Session, e.ApiError) {
	var sessions []*model.Session
	if err := s.db.WithContext(ctx).Where("session_name = ? AND session_type = ?", sessionName, sessionType).Find(&sessions).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, e.NewNotFoundApiError("Sesión no encontrada con el nombre y tipo especificado")
		}
		return nil, e.NewInternalServerApiError("Error encontrando sesiones por nombre y tipo", err)
	}
	return sessions, nil
}

func (s *sessionRepository) GetSessionBySessionKey(ctx context.Context, sessionKey int) (*model.Session, e.ApiError) {
	var session model.Session
	if err := s.db.WithContext(ctx).Where("session_key = ?", sessionKey).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No encontrado, pero no es un error
		}
		return nil, e.NewInternalServerApiError("Error buscando sesión por session_key", err)
	}
	return &session, nil
}

func (s *sessionRepository) GetAllSessions(ctx context.Context) ([]*model.Session, e.ApiError) {
	var sessions []*model.Session
	if err := s.db.WithContext(ctx).Find(&sessions).Error; err != nil {
		return nil, e.NewInternalServerApiError("Error obteniendo todas las sesiones", err)
	}
	return sessions, nil
}

func (s *sessionRepository) GetSessionsByCircuitKeyAndYear(ctx context.Context, circuitKey, year int) ([]*model.Session, e.ApiError) {
	var sessions []*model.Session
	if err := s.db.WithContext(ctx).Where("circuit_key = ? AND year = ?", circuitKey, year).Find(&sessions).Error; err != nil {
		return nil, e.NewInternalServerApiError("Error obteniendo sesiones por circuito y año", err)
	}
	return sessions, nil
}