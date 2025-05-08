package repository

import (
	"context"
	"fmt"
	model "sessions/internal/model"
	e "sessions/pkg/utils"
	"time"

	"gorm.io/gorm"
)

type sessionRepository struct {
	db *gorm.DB
}

type SessionRepository interface{
	CreateSession(ctx context.Context, session *model.Session) e.ApiError
	GetSessionById(ctx context.Context, sessionID int)(*model.Session, e.ApiError)
	UpdateSessionById(ctx context.Context, session *model.Session) e.ApiError
	DeleteSessionById(ctx context.Context, sessionID int) e.ApiError
	GetSessionByYear(ctx context.Context, year int) ([]*model.Session, e.ApiError)
	GetSessionNameAndTypeBySessionID(ctx context.Context, sessionID int) (string, string, e.ApiError)
	GetSessionsByCircuitKey(ctx context.Context, circuitKey int) ([]*model.Session, e.ApiError)
	GetSessionsByCountryCode(ctx context.Context, countryCode string) ([]*model.Session, e.ApiError)
	GetUpcomingSessions(ctx context.Context) ([]*model.Session, e.ApiError)
	GetPastSessions(ctx context.Context, year int) ([]*model.Session, e.ApiError)
	GetSessionsBetweenDates(ctx context.Context, startDate time.Time, endDate time.Time) ([]*model.Session, e.ApiError)
	GetSessionsByNameAndType(ctx context.Context, sessionName string, sessionType string) ([]*model.Session, e.ApiError)
	GetSessionBySessionKey(ctx context.Context, sessionKey int) (*model.Session, e.ApiError)
	GetAllSessions(ctx context.Context) ([]*model.Session, e.ApiError)
	GetSessionsByCircuitKeyAndYear(ctx context.Context, circuitKey, year int) ([]*model.Session, e.ApiError)
	UpdateSCAndVSC(ctx context.Context, sessionID int, sc bool, vsc bool) e.ApiError
	UpdateSessionKey(ctx context.Context, session *model.Session) e.ApiError
	UpdateDFastLap(ctx context.Context, sessionID int, driverID int) e.ApiError
	GetSessionsByLocationAndYear(ctx context.Context, location string, year int) ([]*model.Session, e.ApiError)
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

func (s *sessionRepository) GetSessionById(ctx context.Context, sessionID int)(*model.Session, e.ApiError){
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
    // Asegurarse de que el ID esté presente
    if session.ID == 0 {
        return e.NewBadRequestApiError("El ID de la sesión no puede estar vacío")
    }

    // Usar Updates con un mapa para especificar los campos a actualizar
    updates := map[string]interface{}{
        "weekend_id":        session.WeekendID,
        "circuit_key":       session.CircuitKey,
        "circuit_short_name": session.CircuitShortName,
        "country_code":      session.CountryCode,
        "country_key":       session.CountryKey,
        "country_name":      session.CountryName,
        "date_start":        session.DateStart,
        "date_end":          session.DateEnd,
        "location":          session.Location,
        "session_key":       session.SessionKey,
        "session_name":      session.SessionName,
        "session_type":      session.SessionType,
        "year":              session.Year,
        "dnf":               session.DNF,         // Incluir siempre DNF, incluso si es nil
        "vsc":               session.VSC,         // Incluir siempre VSC, incluso si es nil
        "sf":                session.SF,          // Incluir siempre SF, incluso si es nil
        // "d_fast_lap":        session.DFastLap, // Incluir si lo usas
    }

    // Ejecutar la actualización
    result := s.db.WithContext(ctx).Model(&model.Session{}).Where("id = ?", session.ID).Updates(updates)
    if result.Error != nil {
        return e.NewInternalServerApiError("Error actualizando la sesión", result.Error)
    }

    // Verificar si se actualizó alguna fila
    if result.RowsAffected == 0 {
        return e.NewNotFoundApiError("No se encontró la sesión para actualizar")
    }

    // Verificar los datos actualizados en la base de datos
    updatedSession := &model.Session{}
    if err := s.db.WithContext(ctx).First(updatedSession, session.ID).Error; err != nil {
        return e.NewInternalServerApiError("Error verificando la sesión actualizada", err)
    }

    return nil
}

func (s *sessionRepository) DeleteSessionById(ctx context.Context, sessionID int) e.ApiError {
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

func (s *sessionRepository) GetSessionNameAndTypeBySessionID(ctx context.Context, sessionID int) (string, string, e.ApiError) {
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
    currentYear := currentTime.Year() // Obtiene el año actual (por ejemplo, 2025)

    if err := s.db.WithContext(ctx).
        Where("date_start > ?", currentTime).
        Where("YEAR(date_start) = ?", currentYear).
        Find(&sessions).Error; err != nil {
        return nil, e.NewInternalServerApiError("Error encontrando próximas sesiones", err)
    }
    return sessions, nil
}

func (s *sessionRepository) GetPastSessions(ctx context.Context, year int) ([]*model.Session, e.ApiError) {
    var sessions []*model.Session
    currentTime := time.Now()

    // Filtrar sesiones del año especificado cuya date_start sea anterior a la fecha actual
    if err := s.db.WithContext(ctx).
        Where("date_start < ?", currentTime).
        Where("year = ?", year).
        Find(&sessions).Error; err != nil {
        return nil, e.NewInternalServerApiError("Error encontrando sesiones pasadas", err)
    }
    return sessions, nil
}

func (s *sessionRepository) GetSessionsBetweenDates(ctx context.Context, startDate time.Time, endDate time.Time) ([]*model.Session, e.ApiError) {
    var sessions []*model.Session
    if err := s.db.WithContext(ctx).Where("date_start >= ? AND date_end <= ?", startDate, endDate).Find(&sessions).Error; err != nil {
        return nil, e.NewInternalServerApiError("Error encontrando sesiones entre fechas", err)
    }
    return sessions, nil
}

func (s *sessionRepository) GetSessionsByNameAndType(ctx context.Context, sessionName string, sessionType string) ([]*model.Session, e.ApiError) {
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

func (s *sessionRepository) UpdateSCAndVSC(ctx context.Context, sessionID int, sc bool, vsc bool) e.ApiError {
    // Actualizar solo los campos SC y VSC en la sesión
    if err := s.db.WithContext(ctx).Model(&model.Session{}).Where("id = ?", sessionID).Updates(map[string]interface{}{
        "sf":  sc,
        "vsc": vsc,
    }).Error; err != nil {
        return e.NewInternalServerApiError("Error actualizando SC y VSC en la sesión", err)
    }
    return nil
}

func (s *sessionRepository) UpdateSessionKey(ctx context.Context, session *model.Session) e.ApiError {
	// Actualizar el campo session_key de la sesión
	if err := s.db.WithContext(ctx).Model(&model.Session{}).Where("id = ?", session.ID).Update("session_key", session.SessionKey).Error; err != nil {
		return e.NewInternalServerApiError("Error actualizando el session_key en la base de datos", err)
	}
	return nil
}

// UpdateDFastLap actualiza el valor del campo DFastLap en una sesión específica
func (s *sessionRepository) UpdateDFastLap(ctx context.Context, sessionID int, driverID int) e.ApiError {
    // Actualizar el campo DFastLap de la sesión
    if err := s.db.WithContext(ctx).Model(&model.Session{}).Where("id = ?", sessionID).Update("d_fast_lap", driverID).Error; err != nil {
        return e.NewInternalServerApiError("Error actualizando el DFastLap en la sesión", err)
    }
    return nil
}

// GetSessionsByLocationAndYear obtiene todas las sesiones de un fin de semana y año específicos
func (s *sessionRepository) GetSessionsByLocationAndYear(ctx context.Context, location string, year int) ([]*model.Session, e.ApiError) {
	var sessions []*model.Session
	if err := s.db.WithContext(ctx).Where("location = ? AND year = ?", location, year).Find(&sessions).Error; err != nil {
		return nil, e.NewInternalServerApiError("Error obteniendo sesiones por localización y año", err)
	}
	return sessions, nil
}