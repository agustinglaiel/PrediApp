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
    // cache        *e.Cache
}

type SessionServiceInterface interface{
	CreateSession(ctx context.Context, request dto.CreateSessionDTO) (dto.ResponseSessionDTO, e.ApiError)
	GetSessionById(ctx context.Context, sessionID int) (dto.ResponseSessionDTO, e.ApiError)
	UpdateSessionById(ctx context.Context, sessionID int, request dto.UpdateSessionDTO) (dto.ResponseSessionDTO, e.ApiError)
	DeleteSessionById(ctx context.Context, sessionID int) e.ApiError
	ListSessionsByYear(ctx context.Context, year int) ([]dto.ResponseSessionDTO, e.ApiError)
	GetSessionNameAndTypeById(ctx context.Context, sessionID int) (dto.SessionNameAndTypeDTO, e.ApiError)
	ListSessionsByCircuitKey(ctx context.Context, circuitKey int) ([]dto.ResponseSessionDTO, e.ApiError)
	ListSessionsByCountryCode(ctx context.Context, countryCode string) ([]dto.ResponseSessionDTO, e.ApiError)
	ListUpcomingSessions(ctx context.Context) ([]dto.ResponseSessionDTO, e.ApiError)
    ListPastSessions(ctx context.Context) ([]dto.ResponseSessionDTO, e.ApiError)
	ListSessionsBetweenDates(ctx context.Context, startDate time.Time, endDate time.Time) ([]dto.ResponseSessionDTO, e.ApiError)
	FindSessionsByNameAndType(ctx context.Context, sessionName string, sessionType string) ([]dto.ResponseSessionDTO, e.ApiError)
	GetAllSessions(ctx context.Context) ([]dto.ResponseSessionDTO, e.ApiError)
	GetRaceResultsById(ctx context.Context, sessionID int) (dto.RaceResultsDTO, e.ApiError)
	UpdateResultSCAndVSC(ctx context.Context, sessionID int) e.ApiError 
	UpdateDNFBySessionID(ctx context.Context, sessionID int, dnf int) e.ApiError
	UpdateSessionData(ctx context.Context, sessionID int, location string, sessionName string, sessionType string, year int) e.ApiError
    GetSessionKeyBySessionID(ctx context.Context, sessionID int) (int, e.ApiError)
	UpdateSessionKeyAdmin(ctx context.Context, sessionID int, sessionKey int) e.ApiError
	UpdateDFastLap(ctx context.Context, sessionID int) e.ApiError
}

func NewSessionService(sessionsRepo repository.SessionRepository, client *client.HttpClient, cache *e.Cache) SessionServiceInterface{
	return &sessionService{
		sessionsRepo: sessionsRepo,
		client:       client,  // Pasar el cliente HTTP
        // cache:        cache,
	}
}

func (s *sessionService) CreateSession(ctx context.Context, request dto.CreateSessionDTO) (dto.ResponseSessionDTO, e.ApiError) {

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
    var dnf *int
    var vsc, sf *bool

    if request.SessionType == "Race" && request.SessionName == "Race" {
        // Inicializamos estos valores en nulo, ya que se actualizarán más adelante con datos reales
        dnf = nil
        vsc = nil
        sf = nil
        // dFastLap = nil
    }

    // Validar que no exista ya una combinación idéntica para el mismo fin de semana (mismo location, year, session_name y session_type)
    existingSessions, err := s.sessionsRepo.GetSessionsByLocationAndYear(ctx, request.Location, request.Year)
    if err != nil {
        return dto.ResponseSessionDTO{}, e.NewInternalServerApiError("Error validando sesiones existentes", err)
    }

    // // Imprimir información de las sesiones existentes como JSON para depuración
    // for _, session := range existingSessions {
    //     sessionJSON, err := json.Marshal(session)
    //     if err != nil {
    //         fmt.Printf("Error marshaling session to JSON: %v\n", err)
    //         continue
    //     }
    //     fmt.Printf("Existing Session: %s\n", sessionJSON)
    // }
    
    // Verificar si ya existe una sesión con el mismo SessionName, SessionType y Year para este location
    for _, session := range existingSessions {
        if session.SessionName == request.SessionName &&
           session.SessionType == request.SessionType &&
           session.Year == request.Year { // Añadimos la validación por Year
            return dto.ResponseSessionDTO{}, e.NewBadRequestApiError(
                fmt.Sprintf("Ya existe una sesión con el mismo nombre (%s), tipo (%s) y año (%d) para este fin de semana",
                    request.SessionName, request.SessionType, request.Year),
            )
        }
    }

    // Convertir el DTO en un modelo para guardarlo en la base de datos
    newSession := &model.Session{
        WeekendID:         request.WeekendID,
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
        // DFastLap:          dFastLap,  // ID del piloto que hizo la vuelta más rápida
        CreatedAt:         time.Now(),
        UpdatedAt:         time.Now(),
    }

    if err := s.sessionsRepo.CreateSession(ctx, newSession); err != nil {
        return dto.ResponseSessionDTO{}, e.NewInternalServerApiError("Error creando la sesión", err)
    }

    // Convertir el modelo en un DTO de respuesta
    response := dto.ResponseSessionDTO{
        ID:               newSession.ID,
        WeekendID:        newSession.WeekendID,
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
        // DFastLap:         newSession.DFastLap,  // Actualizar con el nombre correcto
    }

    return response, nil
}

func (s *sessionService) GetSessionById(ctx context.Context, sessionID int) (dto.ResponseSessionDTO, e.ApiError) {
    // Verificas si existe una session con ese sessionID
    session, err := s.sessionsRepo.GetSessionById(ctx, sessionID)
    if err != nil {
        return dto.ResponseSessionDTO{}, err
    }

	// Convert Model to Response DTO
	response := dto.ResponseSessionDTO{
        ID:               session.ID,
        WeekendID:        session.WeekendID,
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
        // DFastLap:         session.DFastLap, // Siempre incluir, incluso si es nil
        VSC:              session.VSC,      // Siempre incluir, incluso si es nil
        SF:               session.SF,       // Siempre incluir, incluso si es nil
        DNF:              session.DNF,      // Siempre incluir, incluso si es nil
    }

	return response, nil
}

func (s *sessionService) UpdateSessionById(ctx context.Context, sessionID int, request dto.UpdateSessionDTO) (dto.ResponseSessionDTO, e.ApiError) {
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
    if request.WeekendID != nil {
        session.WeekendID = *request.WeekendID
    }
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
            // session.DFastLap = request.DFastLap
        }
    }

    // Actualiza la sesión en la base de datos
    if apiErr := s.sessionsRepo.UpdateSessionById(ctx, session); apiErr != nil {
        return dto.ResponseSessionDTO{}, apiErr
    }

    // Construye el DTO de respuesta utilizando los valores actualizados del modelo
    response := dto.ResponseSessionDTO{
        ID:               session.ID,
        WeekendID:        session.WeekendID,
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
        // DFastLap:         session.DFastLap,
    }

    return response, nil
}

func (s *sessionService) DeleteSessionById(ctx context.Context, sessionID int) e.ApiError {
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
    // Llamar al repositorio para obtener las sesiones por año
    sessions, err := s.sessionsRepo.GetSessionByYear(ctx, year)
    if err != nil {
        return nil, err
    }

    // Convertir el resultado de []*model.Session a []dto.ResponseSessionDTO
    var response []dto.ResponseSessionDTO
    for _, session := range sessions {
        response = append(response, dto.ResponseSessionDTO{
            ID:               session.ID,  // Usamos el ID como identificador principal
            WeekendID:        session.WeekendID,
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
            // DFastLap:         session.DFastLap,
            VSC:              session.VSC,
            SF:               session.SF,
            DNF:              session.DNF,
        })
    }

    return response, nil
}

func (s *sessionService) GetSessionNameAndTypeById(ctx context.Context, sessionID int) (dto.SessionNameAndTypeDTO, e.ApiError) {
    // Llamar al repositorio para obtener el nombre y tipo de la sesión
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
    // Llamar al repositorio para obtener las sesiones por CircuitKey
    sessions, err := s.sessionsRepo.GetSessionsByCircuitKey(ctx, circuitKey)
    if err != nil {
        return nil, err
    }

    // Convertir el resultado de []*model.Session a []dto.ResponseSessionDTO
    var response []dto.ResponseSessionDTO
    for _, session := range sessions {
        response = append(response, dto.ResponseSessionDTO{
            ID:               session.ID,
            WeekendID:        session.WeekendID,
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
    // Llamar al repositorio para obtener las sesiones por CountryCode
    sessions, err := s.sessionsRepo.GetSessionsByCountryCode(ctx, countryCode)
    if err != nil {
        return nil, err
    }

    // Convertir el resultado de []*model.Session a []dto.ResponseSessionDTO
    var response []dto.ResponseSessionDTO
    for _, session := range sessions {
        response = append(response, dto.ResponseSessionDTO{
            ID:               session.ID,
            WeekendID:        session.WeekendID,
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
    // Llamar al repositorio para obtener las próximas sesiones
    sessions, err := s.sessionsRepo.GetUpcomingSessions(ctx)
    if err != nil {
        return nil, err
    }

    // Convertir el resultado de []*model.Session a []dto.ResponseSessionDTO
    var response []dto.ResponseSessionDTO
    for _, session := range sessions {
        response = append(response, dto.ResponseSessionDTO{
            ID:               session.ID,
            WeekendID:        session.WeekendID,
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

func (s *sessionService) ListPastSessions(ctx context.Context) ([]dto.ResponseSessionDTO, e.ApiError){
    // Llamar al repositorio para obtener las sesiones pasadas
    sessions, err := s.sessionsRepo.GetPastSessions(ctx)
    if err != nil {
        return nil, err
    }

    // Convertir el resultado de []*model.Session a []dto.ResponseSessionDTO
    var response []dto.ResponseSessionDTO
    for _, session := range sessions {
        response = append(response, dto.ResponseSessionDTO{
            ID:               session.ID,
            WeekendID:        session.WeekendID,
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

    // Si no hay sesiones pasadas, devolver lista vacía
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
            WeekendID:        session.WeekendID,
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

    // Llamar a la función del repositorio para obtener las sesiones por nombre y tipo
    sessions, err := s.sessionsRepo.GetSessionsByNameAndType(ctx, sessionName, sessionType)
    if err != nil {
        return nil, err
    }

    // Convertir el resultado de []*model.Session a []dto.ResponseSessionDTO
    var response []dto.ResponseSessionDTO
    for _, session := range sessions {
        response = append(response, dto.ResponseSessionDTO{
            ID:               session.ID,
            WeekendID:        session.WeekendID,
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
    // Llamar a la función del repositorio para obtener todas las sesiones
    sessions, err := s.sessionsRepo.GetAllSessions(ctx)
    if err != nil {
        return nil, err
    }

    // Convertir el resultado de []*model.Session a []dto.ResponseSessionDTO
    var response []dto.ResponseSessionDTO
    for _, session := range sessions {
        response = append(response, dto.ResponseSessionDTO{
            ID:               session.ID,
            WeekendID:        session.WeekendID,
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

func (s *sessionService) GetRaceResultsById(ctx context.Context, sessionID int) (dto.RaceResultsDTO, e.ApiError) {
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
        // DFastLap:     session.DFastLap,
    }

    return response, nil
}

func (s *sessionService) UpdateResultSCAndVSC(ctx context.Context, sessionID int) e.ApiError {
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

func (s *sessionService) UpdateDNFBySessionID(ctx context.Context, sessionID int, dnf int) e.ApiError {
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

func (s *sessionService) UpdateSessionData(ctx context.Context, sessionID int, location string, sessionName string, sessionType string, year int) e.ApiError {
    // Obtener el session_data desde la API externa usando el cliente HTTP
    sessionData, err := s.client.GetSessionData(location, sessionName, sessionType, year)
    if err != nil {
        return e.NewInternalServerApiError("Error fetching session data", err)
    }

    // Si se encontró un session_key, date_start, y date_end, actualizar la sesión en la base de datos
    if sessionData != nil {
        // Obtener la sesión actual
        session, apiErr := s.sessionsRepo.GetSessionById(ctx, sessionID)
        if apiErr != nil {
            return apiErr
        }

        // Actualizar los campos session_key, date_start y date_end de la sesión
        session.SessionKey = sessionData.SessionKey
        session.DateStart = *sessionData.DateStart
        session.DateEnd = *sessionData.DateEnd
        
        if apiErr := s.sessionsRepo.UpdateSessionById(ctx, session); apiErr != nil {
            return apiErr
        }
    }

    return nil
}

func (s *sessionService) GetSessionKeyBySessionID(ctx context.Context, sessionID int) (int, e.ApiError) {
    // Obtener la sesión por ID
    session, err := s.sessionsRepo.GetSessionById(ctx, sessionID)
    if err != nil {
        return 0, err
    }

    if session.SessionKey == nil {
        return 0, e.NewNotFoundApiError("Session key no encontrado para esta sesión")
    }

    return *session.SessionKey, nil
}

func (s *sessionService) UpdateSessionKeyAdmin(ctx context.Context, sessionID int, sessionKey int) e.ApiError {
    // Obtener la sesión actual por ID
    session, apiErr := s.sessionsRepo.GetSessionById(ctx, sessionID)
    if apiErr != nil {
        return apiErr
    }

    // Actualizar solo el campo session_key con el valor proporcionado manualmente
    session.SessionKey = &sessionKey
    if err := s.sessionsRepo.UpdateSessionKey(ctx, session); err != nil {
        return err
    }

    return nil
}

func (s *sessionService) UpdateDFastLap(ctx context.Context, sessionID int) e.ApiError {
    // Obtener la vuelta más rápida desde el microservicio de results
    fastestLapResult, err := s.client.GetFastestLapBySessionID(sessionID)
    if err != nil {
        fmt.Printf("Error obteniendo el piloto con la vuelta más rápida: %v\n", err)
        return e.NewInternalServerApiError("Error obteniendo el piloto con la vuelta más rápida", err)
    }

    // Verificar que el resultado de la vuelta rápida no esté vacío
    if fastestLapResult == nil {
        fmt.Println("No se encontró una vuelta rápida válida")
        return e.NewNotFoundApiError("No se encontró una vuelta rápida válida para la sesión")
    }

    // Obtener la sesión actual por su ID
    session, apiErr := s.sessionsRepo.GetSessionById(ctx, sessionID)
    if apiErr != nil {
        return apiErr
    }

    // Asignar el ID del piloto que hizo la vuelta más rápida
    // session.DFastLap = &fastestLapResult.Driver.ID

    // Actualizar la sesión utilizando la función UpdateSessionById
    apiErr = s.sessionsRepo.UpdateSessionById(ctx, session)
    if apiErr != nil {
        return apiErr
    }

    fmt.Printf("Vuelta más rápida actualizada correctamente para la sesión %d: Piloto %d\n", sessionID, fastestLapResult.Driver.ID)
    return nil
}