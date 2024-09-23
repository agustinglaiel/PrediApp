package service

import (
	"context"
	"fmt"
	"sessions/internal/client"
	dto "sessions/internal/dto"
	model "sessions/internal/model"
	repository "sessions/internal/repository"
	e "sessions/pkg/utils"
	"time"
)

type sessionService struct {
    sessionsRepo repository.SessionRepository
    client       *client.HttpClient  // Agregar el cliente HTTP
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
	GetRaceResultsById(ctx context.Context, sessionID uint) (dto.RaceResultsDTO, e.ApiError)
	UpdateResultSCAndVSC(ctx context.Context, sessionID uint) e.ApiError 
	//CalculateDNF(ctx context.Context, sessionID uint) e.ApiError
	UpdateDNFBySessionID(ctx context.Context, sessionID uint, dnf int) e.ApiError
	UpdateSessionKey(ctx context.Context, sessionID uint, location, sessionName, sessionType string, year int) e.ApiError
    GetSessionKeyBySessionID(ctx context.Context, sessionID uint) (int, e.ApiError)
}

func NewSessionService(sessionsRepo repository.SessionRepository, client *client.HttpClient) SessionServiceInterface{
	return &sessionService{
		sessionsRepo: sessionsRepo,
		client:       client,  // Pasar el cliente HTTP
	}
}

func (s *sessionService) CreateSession(ctx context.Context, request dto.CreateSessionDTO) (dto.ResponseSessionDTO, e.ApiError) {
    /* ESTO LO ELIMINO PORQUE DEJAMOS DE PASAR EL SESSION_KEY SINO QUE AHORA LO OBTENEMOS A PARTIR DE LA API
	// Validar que session_key sea único
    existingSession, _ := s.sessionsRepo.GetSessionBySessionKey(ctx, request.SessionKey)
    if existingSession != nil {
        return dto.ResponseSessionDTO{}, e.NewBadRequestApiError("session_key ya está en uso")
    }
	*/

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

    // Para todas las sesiones que no son "Race", estos campos deben ser nulos
    var dnf, dFastLap *int
    var vsc, sf *bool

    if request.SessionType == "Race" && request.SessionName == "Race" {
        // Inicializamos estos valores en nulo, ya que se actualizarán más adelante con datos reales
        dnf = nil
        vsc = nil
        sf = nil
        dFastLap = nil
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

    // Convertir el DTO en un modelo para guardarlo en la base de datos
    newSession := &model.Session{
        CircuitKey:        request.CircuitKey,
        CircuitShortName:  request.CircuitShortName,
        CountryCode:       request.CountryCode,
        CountryName:       request.CountryName,
        DateStart:         request.DateStart,
        DateEnd:           request.DateEnd,
        Location:          request.Location,
        SessionKey:        nil,
        SessionName:       request.SessionName,
        SessionType:       request.SessionType,
        Year:              request.Year,
        DNF:               dnf,       // Número de pilotos que no terminaron
        VSC:               vsc,       // Indica si hubo VSC
        SF:                sf,        // Indica si hubo SC
        DFastLap:          dFastLap,  // ID del piloto que hizo la vuelta más rápida
        CreatedAt:         time.Now(),
        UpdatedAt:         time.Now(),
    }

    if err := s.sessionsRepo.CreateSession(ctx, newSession); err != nil {
        return dto.ResponseSessionDTO{}, e.NewInternalServerApiError("Error creando la sesión", err)
    }

    // Verificar el ID después de la creación
    fmt.Printf("Session ID after creation: %d\n", newSession.ID)

    // Convertir el modelo en un DTO de respuesta
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
        DNF:              newSession.DNF,
        VSC:              newSession.VSC,
        SF:               newSession.SF,
        DFastLap:         newSession.DFastLap,  // Actualizar con el nombre correcto
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
		SessionKey:       session.SessionKey,
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

    // Si se intenta actualizar session_name o session_type, validar la combinación
    if request.SessionName != nil || request.SessionType != nil {
        validCombinations := map[string]string{
            "Sprint Qualifying": "Qualifying",
            "Sprint":            "Race",
            "Practice 1":        "Practice",
            "Practice 2":        "Practice",
            "Practice 3":        "Practice",
            "Qualifying":        "Qualifying",
            "Race":              "Race",
        }

        // Validar que la combinación sea correcta si ambos campos están presentes
        if request.SessionName != nil && request.SessionType != nil {
            expectedType, ok := validCombinations[*request.SessionName]
            if !ok || expectedType != *request.SessionType {
                return dto.ResponseSessionDTO{}, e.NewBadRequestApiError("Combinación inválida de session_name y session_type")
            }
        } else if request.SessionName != nil {
            // Si solo se está actualizando session_name, mantener el session_type existente para validar
            expectedType, ok := validCombinations[*request.SessionName]
            if !ok || expectedType != session.SessionType {
                return dto.ResponseSessionDTO{}, e.NewBadRequestApiError("Combinación inválida de session_name con el session_type existente")
            }
        } else if request.SessionType != nil {
            // Si solo se está actualizando session_type, mantener el session_name existente para validar
            expectedType, ok := validCombinations[session.SessionName]
            if !ok || expectedType != *request.SessionType {
                return dto.ResponseSessionDTO{}, e.NewBadRequestApiError("Combinación inválida de session_type con el session_name existente")
            }
        }
    }

    // Validar que no exista ya una combinación idéntica para el mismo fin de semana (mismo circuito y año)
    if request.SessionName != nil && request.SessionType != nil && request.CircuitKey != nil && request.Year != nil {
        existingSessions, err := s.sessionsRepo.GetSessionsByCircuitKeyAndYear(ctx, *request.CircuitKey, *request.Year)
        if err != nil {
            return dto.ResponseSessionDTO{}, e.NewInternalServerApiError("Error validando sesiones existentes", err)
        }

        for _, existingSession := range existingSessions {
            if existingSession.ID != sessionID && existingSession.SessionName == *request.SessionName && existingSession.SessionType == *request.SessionType {
                return dto.ResponseSessionDTO{}, e.NewBadRequestApiError("Ya existe una sesión con el mismo nombre y tipo para este fin de semana")
            }
        }
    }

    // Actualiza solo los campos que están presentes en el DTO de actualización
    if request.CircuitKey != nil {
        session.CircuitKey = *request.CircuitKey
    }
    if request.CircuitShortName != nil {
        session.CircuitShortName = *request.CircuitShortName
    }
    if request.CountryCode != nil {
        session.CountryCode = *request.CountryCode
    }
    if request.CountryKey != nil {
        session.CountryKey = *request.CountryKey
    }
    if request.CountryName != nil {
        session.CountryName = *request.CountryName
    }
    if request.DateStart != nil {
        session.DateStart = *request.DateStart
    }
    if request.DateEnd != nil {
        session.DateEnd = *request.DateEnd
    }
    if request.Location != nil {
        session.Location = *request.Location
    }
    if request.SessionKey != nil {
        session.SessionKey = request.SessionKey
    }
    if request.SessionName != nil {
        session.SessionName = *request.SessionName
    }
    if request.SessionType != nil {
        session.SessionType = *request.SessionType
    }
    if request.Year != nil {
        session.Year = *request.Year
    }

    // Si la sesión es de tipo "Race", actualizamos los campos relacionados con resultados de la carrera
    if session.SessionType == "Race" && session.SessionName == "Race" {
        if request.DNF != nil {
            session.DNF = request.DNF
        }
        if request.VSC != nil {
            session.VSC = request.VSC
        }
        if request.SF != nil {
            session.SF = request.SF
        }
        if request.DFastLap != nil {
            session.DFastLap = request.DFastLap
        }
    }

    // Actualiza la sesión en la base de datos
    if apiErr := s.sessionsRepo.UpdateSessionById(ctx, session); apiErr != nil {
        return dto.ResponseSessionDTO{}, apiErr
    }

    // Construye el DTO de respuesta utilizando los valores actualizados del modelo
    response := dto.ResponseSessionDTO{
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
        DNF:              session.DNF,
        VSC:              session.VSC,
        SF:               session.SF,
        DFastLap:         session.DFastLap,
    }

    return response, nil
}

func (s *sessionService) DeleteSessionById(ctx context.Context, sessionID uint) e.ApiError {
	// Verificar si la sesión existe antes de intentar eliminarla
	session, apiErr := s.sessionsRepo.GetSessionById(ctx, sessionID)
	if apiErr != nil {
		return apiErr
	}

	// (Opcional) Log de eliminación
	fmt.Printf("Eliminando la sesión con ID: %d\n", sessionID)

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
			// Puedes incluir aquí los campos adicionales si es relevante
			DFastLap:         session.DFastLap,
			VSC:              session.VSC,
			SF:               session.SF,
			DNF:              session.DNF,
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

	// Si no hay próximas sesiones, devolver lista vacía
	if len(response) == 0 {
		return []dto.ResponseSessionDTO{}, nil
	}

	return response, nil
}

func (s *sessionService) ListSessionsBetweenDates(ctx context.Context, startDate time.Time, endDate time.Time) ([]dto.ResponseSessionDTO, e.ApiError) {
	// Validar que endDate no sea anterior a startDate
	if endDate.Before(startDate) {
		return nil, e.NewBadRequestApiError("La fecha de finalización no puede ser anterior a la fecha de inicio")
	}

	// Llamar a la función del repository para obtener las sesiones entre las fechas especificadas
	sessions, err := s.sessionsRepo.GetSessionsBetweenDates(ctx, startDate, endDate)
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

	// Si no hay sesiones entre las fechas, devolver una lista vacía
	if len(response) == 0 {
		return []dto.ResponseSessionDTO{}, nil
	}

	return response, nil
}

func (s *sessionService) FindSessionsByNameAndType(ctx context.Context, sessionName string, sessionType string) ([]dto.ResponseSessionDTO, e.ApiError) {
	// Validar que sessionName y sessionType no estén vacíos
	if sessionName == "" || sessionType == "" {
		return nil, e.NewBadRequestApiError("El nombre y tipo de sesión son obligatorios")
	}

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

	// Si no se encuentran sesiones con el nombre y tipo especificado, devolver lista vacía
	if len(response) == 0 {
		return []dto.ResponseSessionDTO{}, nil
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

	// Si no se encuentran sesiones, devolver lista vacía
	if len(response) == 0 {
		return []dto.ResponseSessionDTO{}, nil
	}

	return response, nil
}

func (s *sessionService) GetRaceResultsById(ctx context.Context, sessionID uint) (dto.RaceResultsDTO, e.ApiError) {
    // Obtener la sesión por ID
    session, apiErr := s.sessionsRepo.GetSessionById(ctx, sessionID)
    if apiErr != nil {
        return dto.RaceResultsDTO{}, apiErr
    }

    // Validar que la sesión sea de tipo "Race"
    if session.SessionType != "Race" || session.SessionName != "Race" {
        return dto.RaceResultsDTO{}, e.NewBadRequestApiError("Solo las sesiones de tipo 'Race' tienen resultados de carrera")
    }

    // Construir y devolver el DTO de resultados de carrera
    response := dto.RaceResultsDTO{
        DNF:          session.DNF,
        VSC:          session.VSC,
        SF:           session.SF,
        DFastLap:     session.DFastLap,
    }

    return response, nil
}

func (s *sessionService) UpdateResultSCAndVSC(ctx context.Context, sessionID uint) e.ApiError {
    // Obtener la sesión por su ID para tener el SessionKey
    session, apiErr := s.sessionsRepo.GetSessionById(ctx, sessionID)
    if apiErr != nil {
        return apiErr
    }

	fmt.Printf("Session Key: %v\n", session.SessionKey)

    // Usar el SessionKey para hacer la llamada a la API externa
    raceControlData, err := s.client.GetRaceControlData(session.SessionKey)
    if err != nil {
        return e.NewInternalServerApiError("Error fetching race control data", err)
    }

	// Depurar los datos de control de carrera
	//fmt.Printf("Race Control Data: %+v\n", raceControlData)

    // Procesar los datos de control de carrera para actualizar VSC y SC
    var vsc, sc bool
    for _, event := range raceControlData {
        if event.Category == "SafetyCar" {
            if event.Message == "VIRTUAL SAFETY CAR DEPLOYED" {
                vsc = true
            } else if event.Message == "SAFETY CAR DEPLOYED" {
                sc = true
            }
        }
    }

	fmt.Printf("VSC: %t, SC: %t\n", vsc, sc)

    // Llamar al repository para actualizar solo los campos SC y VSC
	//fmt.Printf("Actualizando SC: %t, VSC: %t para la session ID: %d\n", sc, vsc, sessionID)
    if err := s.sessionsRepo.UpdateSCAndVSC(ctx, sessionID, sc, vsc); err != nil {
        return e.NewInternalServerApiError("Error updating SC and VSC in session", err)
    }

    return nil
}

func (s *sessionService) UpdateDNFBySessionID(ctx context.Context, sessionID uint, dnf int) e.ApiError {
    // Obtener la sesión por su ID para asegurarnos que existe y que sea una carrera
    session, apiErr := s.sessionsRepo.GetSessionById(ctx, sessionID)
    if apiErr != nil {
        return apiErr
    }

    // Validar que sea una sesión de tipo "Race"
    if session.SessionName != "Race" || session.SessionType != "Race" {
        return e.NewBadRequestApiError("La sesión no es una carrera (Race)")
    }

    // Actualizar el valor de DNF en la sesión
    session.DNF = &dnf

    // Guardar la actualización en el repositorio
    if err := s.sessionsRepo.UpdateSessionById(ctx, session); err != nil {
        return e.NewInternalServerApiError("Error actualizando la cantidad de DNF", err)
    }

    return nil
}

func (s *sessionService) UpdateSessionKey(ctx context.Context, sessionID uint, location, sessionName, sessionType string, year int) e.ApiError {
	// Obtener el session_key desde la API externa usando el cliente HTTP
	sessionKey, err := s.client.GetSessionKey(location, sessionName, sessionType, year)
	if err != nil {
		return e.NewInternalServerApiError("Error fetching session key", err)
	}

	// Si se encontró un session_key, actualizar la sesión en la base de datos
	if sessionKey != nil {
		// Obtener la sesión actual
		session, apiErr := s.sessionsRepo.GetSessionById(ctx, sessionID)
		if apiErr != nil {
			return apiErr
		}

		// Actualizar el session_key de la sesión
		session.SessionKey = sessionKey
		if apiErr := s.sessionsRepo.UpdateSessionKey(ctx, session); apiErr != nil {
			return apiErr
		}
	}

	return nil
}

/*
func (s *sessionService) CalculateDNF(ctx context.Context, sessionID uint) e.ApiError {
    // Obtener la sesión de la base de datos
    session, err := s.sessionsRepo.GetSessionById(ctx, sessionID)
    if err != nil {
        return err
    }

    // Validar que sea una sesión de tipo "Race"
    if session.SessionType != "Race" || session.SessionName != "Race" {
        return e.NewBadRequestApiError("Esta sesión no es de tipo 'Race'")
    }

    // Obtener el session_key (clave única) para la API externa
    sessionKey := session.SessionKey

    // Obtener las vueltas de todos los pilotos para la sesión a través del cliente HTTP
    lapsData, err := s.client.GetLapsData(sessionKey)
    if err != nil {
        return e.NewInternalServerApiError("Error obteniendo los datos de las vueltas", err)  // Adaptamos el error
    }

    // Lógica para calcular cuántos pilotos no completaron la carrera (DNF)
    totalDNF := 0
    for _, lap := range lapsData {
        // Si el tiempo de vuelta es nulo o cero, se considera que el piloto no terminó
        if lap.LapDuration == nil || *lap.LapDuration == 0 {  // Corregimos la comparación
            totalDNF++
        }
    }

    // Actualizar el campo DNF de la sesión con el valor calculado
    session.DNF = &totalDNF
    if err := s.sessionsRepo.UpdateSessionById(ctx, session); err != nil {
        return e.NewInternalServerApiError("Error actualizando la sesión con el valor de DNF", err)
    }

    return nil
}
*/

func (s *sessionService) GetSessionKeyBySessionID(ctx context.Context, sessionID uint) (int, e.ApiError) {
    session, err := s.sessionsRepo.GetSessionById(ctx, sessionID)
    if err != nil {
        return 0, err
    }

    if session.SessionKey == nil {
        return 0, e.NewNotFoundApiError("Session key no encontrado para esta sesión")
    }

    return *session.SessionKey, nil
}
