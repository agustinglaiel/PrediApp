package service

import (
	dto "admin/internal/dto/sessions"
	model "admin/internal/model/sessions"
	repository "admin/internal/repository/sessions"
	e "admin/pkg/utils"
	"context"
	"fmt"
	"time"
)

type sessionService struct{
	sessionsRepo repository.SessionRepository
}

type SessionServiceInterface interface{
	CreateSession(ctx context.Context, request dto.CreateSessionDTO) (dto.ResponseSessionDTO, e.ApiError)
	GetSessionById(ctx context.Context, sessionID uint) (dto.ResponseSessionDTO, e.ApiError)
	UpdateSessionById(ctx context.Context, sessionID uint, request dto.UpdateSessionDTO) (dto.ResponseSessionDTO, e.ApiError)
	DeleteSessionById(ctx context.Context, sessionID uint) e.ApiError
	ListSessionsByYear(ctx context.Context, year int) ([]dto.ResponseSessionDTO, e.ApiError)
	GetSessionNameAndTypeById(ctx context.Context, sessionID uint) (dto.SessionNameAndTypeDTO, e.ApiError)
	ListSessionsByCircuitKey(ctx context.Context, circuitKey int) ([]dto.ResponseSessionDTO, e.ApiError)
	ListSessionsByCountryCode(ctx context.Context, countryCode string) ([]dto.ResponseSessionDTO, e.ApiError)
	ListUpcomingSessions(ctx context.Context) ([]dto.ResponseSessionDTO, e.ApiError)
	ListSessionsBetweenDates(ctx context.Context, startDate time.Time, endDate time.Time) ([]dto.ResponseSessionDTO, e.ApiError)
	FindSessionsByNameAndType(ctx context.Context, sessionName string, sessionType string) ([]dto.ResponseSessionDTO, e.ApiError)
	GetAllSessions(ctx context.Context) ([]dto.ResponseSessionDTO, e.ApiError)
}

func NewSessionService(sessionsRepo repository.SessionRepository) SessionServiceInterface{
	return &sessionService{
		sessionsRepo: sessionsRepo,
	}
}

func (s *sessionService) CreateSession(ctx context.Context, request dto.CreateSessionDTO) (dto.ResponseSessionDTO, e.ApiError) {
	// Validar que session_key sea único
	existingSession, _ := s.sessionsRepo.GetSessionBySessionKey(ctx, request.SessionKey)
	if existingSession != nil {
		return dto.ResponseSessionDTO{}, e.NewBadRequestApiError("session_key ya está en uso")
	}

	 // Validar que la combinación de session_name y session_type sea válida
	 validCombinations := map[string]string{
        "Sprint Qualifying": "Qualifying",
        "Sprint":            "Race",
        "Practice 1":        "Practice",
        "Practice 2":        "Practice",
        "Practice 3":        "Practice",
        "Qualifying":        "Qualifying",
        "Race":              "Race",
    }

	expectedType, ok := validCombinations[request.SessionName]
    if !ok || expectedType != request.SessionType {
        return dto.ResponseSessionDTO{}, e.NewBadRequestApiError("Combinación inválida de session_name y session_type")
    }

    // Validar que no exista ya una combinación idéntica para el mismo fin de semana (mismo circuito y año)
    existingSessions, err := s.sessionsRepo.GetSessionsByCircuitKeyAndYear(ctx, request.CircuitKey, request.Year)
    if err != nil {
        return dto.ResponseSessionDTO{}, e.NewInternalServerApiError("Error validando sesiones existentes", err)
    }

    for _, session := range existingSessions {
        if session.SessionName == request.SessionName && session.SessionType == request.SessionType {
            return dto.ResponseSessionDTO{}, e.NewBadRequestApiError("Ya existe una sesión con el mismo nombre y tipo para este fin de semana")
        }
    }

	// Convert DTO to Model
	newSession := &model.Session{
        CircuitKey:       request.CircuitKey,
        CircuitShortName: request.CircuitShortName,
        CountryCode:      request.CountryCode,
        CountryName:      request.CountryName,
        DateStart:        request.DateStart,
        DateEnd:          request.DateEnd,
        Location:         request.Location,
        SessionKey:       request.SessionKey,
        SessionName:      request.SessionName,
        SessionType:      request.SessionType,
        Year:             request.Year,
        CreatedAt:        time.Now(),
        UpdatedAt:        time.Now(),
    }

	if err := s.sessionsRepo.CreateSession(ctx, newSession); err != nil {
		return dto.ResponseSessionDTO{}, e.NewInternalServerApiError("Error creando la sesión", err)
	}

	// Verificar el ID después de la creación
    fmt.Printf("Session ID after creation: %d\n", newSession.ID)

	// Convert Model to Response DTO
	response := dto.ResponseSessionDTO{
        ID:               newSession.ID,
        CircuitKey:       newSession.CircuitKey,
        CircuitShortName: newSession.CircuitShortName,
        CountryCode:      newSession.CountryCode,
        CountryName:      newSession.CountryName,
        DateStart:        newSession.DateStart,
        DateEnd:          newSession.DateEnd,
        Location:         newSession.Location,
        SessionKey:       newSession.SessionKey,
        SessionName:      newSession.SessionName,
        SessionType:      newSession.SessionType,
        Year:             newSession.Year,
    }

	return response, nil
}

func (s *sessionService) GetSessionById(ctx context.Context, sessionID uint) (dto.ResponseSessionDTO, e.ApiError) {
	session, err := s.sessionsRepo.GetSessionById(ctx, sessionID)
	if err != nil {
		return dto.ResponseSessionDTO{}, err
	}

	// Convert Model to Response DTO
	response := dto.ResponseSessionDTO{
		ID:               session.ID,
		CircuitKey:       session.CircuitKey,
		CircuitShortName: session.CircuitShortName,
		CountryCode:      session.CountryCode,
		CountryName:      session.CountryName,
		DateStart:        session.DateStart,
		DateEnd:          session.DateEnd,
		Location:         session.Location,
		SessionName:      session.SessionName,
		SessionType:      session.SessionType,
		Year:             session.Year,
	}

	return response, nil
}

func (s *sessionService) UpdateSessionById(ctx context.Context, sessionID uint, request dto.UpdateSessionDTO) (dto.ResponseSessionDTO, e.ApiError) {
	// Obtén la sesión existente por su ID
	session, apiErr := s.sessionsRepo.GetSessionById(ctx, sessionID)
	if apiErr != nil {
		return dto.ResponseSessionDTO{}, apiErr
	}

	// Validar que la combinación de session_name y session_type sea válida
	if request.SessionName != "" || request.SessionType != "" {
		validCombinations := map[string]string{
			"Sprint Qualifying": "Qualifying",
			"Sprint":            "Race",
			"Practice 1":        "Practice",
			"Practice 2":        "Practice",
			"Practice 3":        "Practice",
			"Qualifying":        "Qualifying",
			"Race":              "Race",
		}

		if request.SessionName != "" && request.SessionType != "" {
			expectedType, ok := validCombinations[request.SessionName]
			if !ok || expectedType != request.SessionType {
				return dto.ResponseSessionDTO{}, e.NewBadRequestApiError("Combinación inválida de session_name y session_type")
			}
		}
	}

	// Validar que no exista ya una combinación idéntica para el mismo fin de semana (mismo circuito y año)
	if request.SessionName != "" && request.SessionType != "" && request.CircuitKey != 0 && request.Year != 0 {
		existingSessions, err := s.sessionsRepo.GetSessionsByCircuitKeyAndYear(ctx, request.CircuitKey, request.Year)
		if err != nil {
			return dto.ResponseSessionDTO{}, e.NewInternalServerApiError("Error validando sesiones existentes", err)
		}

		for _, existingSession := range existingSessions {
			if existingSession.ID != sessionID && existingSession.SessionName == request.SessionName && existingSession.SessionType == request.SessionType {
				return dto.ResponseSessionDTO{}, e.NewBadRequestApiError("Ya existe una sesión con el mismo nombre y tipo para este fin de semana")
			}
		}
	}

	// Actualiza solo los campos que están presentes en el DTO de actualización
	if request.CircuitKey != 0 {
		session.CircuitKey = request.CircuitKey
	}
	if request.CircuitShortName != "" {
		session.CircuitShortName = request.CircuitShortName
	}
	if request.CountryCode != "" {
		session.CountryCode = request.CountryCode
	}
	if request.CountryKey != 0 {
		session.CountryKey = request.CountryKey
	}
	if request.CountryName != "" {
		session.CountryName = request.CountryName
	}
	if !request.DateStart.IsZero() {
		session.DateStart = request.DateStart
	}
	if !request.DateEnd.IsZero() {
		session.DateEnd = request.DateEnd
	}
	if request.Location != "" {
		session.Location = request.Location
	}
	if request.SessionKey != 0 {
		session.SessionKey = request.SessionKey
	}
	if request.SessionName != "" {
		session.SessionName = request.SessionName
	}
	if request.SessionType != "" {
		session.SessionType = request.SessionType
	}
	if request.Year != 0 {
		session.Year = request.Year
	}

	// Actualiza la sesión en la base de datos
	if apiErr := s.sessionsRepo.UpdateSessionById(ctx, session); apiErr != nil {
		return dto.ResponseSessionDTO{}, apiErr
	}

	// Construye el DTO de respuesta utilizando los valores actualizados del modelo
	response := dto.ResponseSessionDTO{
		ID:               uint(session.ID),
		CircuitKey:       session.CircuitKey,
		CircuitShortName: session.CircuitShortName,
		CountryCode:      session.CountryCode,
		CountryName:      session.CountryName,
		DateStart:        session.DateStart,
		DateEnd:          session.DateEnd,
		Location:         session.Location,
		SessionKey:       session.SessionKey,
		SessionName:      session.SessionName,
		SessionType:      session.SessionType,
		Year:             session.Year,
	}

	return response, nil
}

func (s *sessionService) DeleteSessionById(ctx context.Context, sessionID uint) e.ApiError {
	// Verificar si la sesión existe antes de intentar eliminarla
	session, apiErr := s.sessionsRepo.GetSessionById(ctx, sessionID)
	if apiErr != nil {
		return apiErr
	}

	// Eliminar la sesión utilizando el ID
	if err := s.sessionsRepo.DeleteSessionById(ctx, session.ID); err != nil {
		return e.NewInternalServerApiError("Error eliminando la sesión", err)
	}

	return nil
}

func (s *sessionService) ListSessionsByYear(ctx context.Context, year int) ([]dto.ResponseSessionDTO, e.ApiError) {
	// Llamar a la función del repository para obtener las sesiones por año
	sessions, err := s.sessionsRepo.GetSessionByYear(ctx, year)
	if err != nil {
		return nil, err
	}

	// Convertir el resultado de []*model.Session a []dto.ResponseSessionDTO
	var response []dto.ResponseSessionDTO
	for _, session := range sessions {
		response = append(response, dto.ResponseSessionDTO{
			ID:               session.ID,  // Usamos el ID como identificador principal
			CircuitKey:       session.CircuitKey,
			CircuitShortName: session.CircuitShortName,
			CountryCode:      session.CountryCode,
			CountryName:      session.CountryName,
			DateStart:        session.DateStart,
			DateEnd:          session.DateEnd,
			Location:         session.Location,
			SessionKey:       session.SessionKey,
			SessionName:      session.SessionName,
			SessionType:      session.SessionType,
			Year:             session.Year,
		})
	}

	return response, nil
}

func (s *sessionService) GetSessionNameAndTypeById(ctx context.Context, sessionID uint) (dto.SessionNameAndTypeDTO, e.ApiError) {
	// Llamar a la función del repository para obtener el nombre y tipo de la sesión
	sessionName, sessionType, err := s.sessionsRepo.GetSessionNameAndTypeBySessionID(ctx, sessionID)
	if err != nil {
		return dto.SessionNameAndTypeDTO{}, err
	}

	// Construir el DTO de respuesta
	response := dto.SessionNameAndTypeDTO{
		SessionName: sessionName,
		SessionType: sessionType,
	}

	return response, nil
}

func (s *sessionService) ListSessionsByCircuitKey(ctx context.Context, circuitKey int) ([]dto.ResponseSessionDTO, e.ApiError) {
	// Llamar a la función del repository para obtener las sesiones por CircuitKey
	sessions, err := s.sessionsRepo.GetSessionsByCircuitKey(ctx, circuitKey)
	if err != nil {
		return nil, err
	}

	// Convertir el resultado de []*model.Session a []dto.ResponseSessionDTO
	var response []dto.ResponseSessionDTO
	for _, session := range sessions {
		response = append(response, dto.ResponseSessionDTO{
			ID:               session.ID,  // Usamos el ID como identificador principal
			CircuitKey:       session.CircuitKey,
			CircuitShortName: session.CircuitShortName,
			CountryCode:      session.CountryCode,
			CountryName:      session.CountryName,
			DateStart:        session.DateStart,
			DateEnd:          session.DateEnd,
			Location:         session.Location,
			SessionKey:       session.SessionKey,
			SessionName:      session.SessionName,
			SessionType:      session.SessionType,
			Year:             session.Year,
		})
	}

	return response, nil
}

func (s *sessionService) ListSessionsByCountryCode(ctx context.Context, countryCode string) ([]dto.ResponseSessionDTO, e.ApiError) {
	// Llamar a la función del repository para obtener las sesiones por CountryCode
	sessions, err := s.sessionsRepo.GetSessionsByCountryCode(ctx, countryCode)
	if err != nil {
		return nil, err
	}

	// Convertir el resultado de []*model.Session a []dto.ResponseSessionDTO
	var response []dto.ResponseSessionDTO
	for _, session := range sessions {
		response = append(response, dto.ResponseSessionDTO{
			ID:               session.ID,  // Usamos el ID como identificador principal
			CircuitKey:       session.CircuitKey,
			CircuitShortName: session.CircuitShortName,
			CountryCode:      session.CountryCode,
			CountryName:      session.CountryName,
			DateStart:        session.DateStart,
			DateEnd:          session.DateEnd,
			Location:         session.Location,
			SessionKey:       session.SessionKey,
			SessionName:      session.SessionName,
			SessionType:      session.SessionType,
			Year:             session.Year,
		})
	}

	return response, nil
}

func (s *sessionService) ListUpcomingSessions(ctx context.Context) ([]dto.ResponseSessionDTO, e.ApiError) {
	// Llamar a la función del repository para obtener las próximas sesiones
	sessions, err := s.sessionsRepo.GetUpcomingSessions(ctx)
	if err != nil {
		return nil, err
	}

	// Convertir el resultado de []*model.Session a []dto.ResponseSessionDTO
	var response []dto.ResponseSessionDTO
	for _, session := range sessions {
		response = append(response, dto.ResponseSessionDTO{
			ID:               session.ID,  // Usamos el ID como identificador principal
			CircuitKey:       session.CircuitKey,
			CircuitShortName: session.CircuitShortName,
			CountryCode:      session.CountryCode,
			CountryName:      session.CountryName,
			DateStart:        session.DateStart,
			DateEnd:          session.DateEnd,
			Location:         session.Location,
			SessionKey:       session.SessionKey,
			SessionName:      session.SessionName,
			SessionType:      session.SessionType,
			Year:             session.Year,
		})
	}

	return response, nil
}

func (s *sessionService) ListSessionsBetweenDates(ctx context.Context, startDate time.Time, endDate time.Time) ([]dto.ResponseSessionDTO, e.ApiError) {
	// Llamar a la función del repository para obtener las sesiones entre las fechas especificadas
	sessions, err := s.sessionsRepo.GetSessionsBetweenDates(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Convertir el resultado de []*model.Session a []dto.ResponseSessionDTO
	var response []dto.ResponseSessionDTO
	for _, session := range sessions {
		response = append(response, dto.ResponseSessionDTO{
			ID:               session.ID,  // Usamos el ID como identificador principal
			CircuitKey:       session.CircuitKey,
			CircuitShortName: session.CircuitShortName,
			CountryCode:      session.CountryCode,
			CountryName:      session.CountryName,
			DateStart:        session.DateStart,
			DateEnd:          session.DateEnd,
			Location:         session.Location,
			SessionKey:       session.SessionKey,
			SessionName:      session.SessionName,
			SessionType:      session.SessionType,
			Year:             session.Year,
		})
	}

	return response, nil
}

func (s *sessionService) FindSessionsByNameAndType(ctx context.Context, sessionName string, sessionType string) ([]dto.ResponseSessionDTO, e.ApiError) {
	// Llamar a la función del repository para obtener las sesiones por nombre y tipo
	sessions, err := s.sessionsRepo.GetSessionsByNameAndType(ctx, sessionName, sessionType)
	if err != nil {
		return nil, err
	}

	// Convertir el resultado de []*model.Session a []dto.ResponseSessionDTO
	var response []dto.ResponseSessionDTO
	for _, session := range sessions {
		response = append(response, dto.ResponseSessionDTO{
			ID:               session.ID,
			CircuitKey:       session.CircuitKey,
			CircuitShortName: session.CircuitShortName,
			CountryCode:      session.CountryCode,
			CountryName:      session.CountryName,
			DateStart:        session.DateStart,
			DateEnd:          session.DateEnd,
			Location:         session.Location,
			SessionKey:       session.SessionKey,
			SessionName:      session.SessionName,
			SessionType:      session.SessionType,
			Year:             session.Year,
		})
	}

	return response, nil
}

func (s *sessionService) GetAllSessions(ctx context.Context) ([]dto.ResponseSessionDTO, e.ApiError) {
	// Llamar a la función del repository para obtener todas las sesiones
	sessions, err := s.sessionsRepo.GetAllSessions(ctx)
	if err != nil {
		return nil, err
	}

	// Convertir el resultado de []*model.Session a []dto.ResponseSessionDTO
	var response []dto.ResponseSessionDTO
	for _, session := range sessions {
		response = append(response, dto.ResponseSessionDTO{
			ID:               session.ID,
			CircuitKey:       session.CircuitKey,
			CircuitShortName: session.CircuitShortName,
			CountryCode:      session.CountryCode,
			CountryName:      session.CountryName,
			DateStart:        session.DateStart,
			DateEnd:          session.DateEnd,
			Location:         session.Location,
			SessionKey:       session.SessionKey,
			SessionName:      session.SessionName,
			SessionType:      session.SessionType,
			Year:             session.Year,
		})
	}

	return response, nil
}