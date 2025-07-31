package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	e "prediapp.local/results/pkg/utils"

	"prediapp.local/db/model"
	"prediapp.local/results/internal/client"
	"prediapp.local/results/internal/dto"
	"prediapp.local/results/internal/repository"

	"gorm.io/gorm"
)

type resultService struct {
	resultRepo     repository.ResultRepository
	driversClient  *client.HttpClient
	sessionsClient *client.HttpClient
	usersClient    *client.HttpClient
	externalClient *client.HttpClient
	cache          *e.Cache
}

type ResultService interface {
	FetchResultsFromExternalAPI(ctx context.Context, sessionId int) ([]dto.ResponseResultDTO, e.ApiError)
	FetchNonRaceSessionResults(ctx context.Context, sessionId int) ([]dto.ResponseResultDTO, e.ApiError)
	UpdateResult(ctx context.Context, resultID int, request dto.UpdateResultDTO) (dto.ResponseResultDTO, e.ApiError)
	GetResultsOrderedByPosition(ctx context.Context, sessionID int) ([]dto.ResponseResultDTO, e.ApiError)
	GetFastestLapInSession(ctx context.Context, sessionID int) (dto.ResponseResultDTO, e.ApiError)
	CreateResult(ctx context.Context, request dto.CreateResultDTO) (dto.ResponseResultDTO, e.ApiError)
	DeleteResult(ctx context.Context, resultID int) e.ApiError
	GetAllResults(ctx context.Context) ([]dto.ResponseResultDTO, e.ApiError)
	GetTopNDriversInSession(ctx context.Context, sessionID int, n int) ([]dto.TopDriverDTO, e.ApiError)
	DeleteAllResultsForSession(ctx context.Context, sessionID int) e.ApiError
	CreateSessionResultsAdmin(ctx context.Context, bulkRequest dto.CreateBulkResultsDTO) ([]dto.ResponseResultDTO, e.ApiError)
}

func NewResultService(
	resultRepo repository.ResultRepository,
	driversClient *client.HttpClient,
	sessionsClient *client.HttpClient,
	usersClient *client.HttpClient,
	externalClient *client.HttpClient,
	cache *e.Cache,
) ResultService {
	return &resultService{
		resultRepo:     resultRepo,
		driversClient:  driversClient,
		sessionsClient: sessionsClient,
		usersClient:    usersClient,
		externalClient: externalClient,
		cache:          cache,
	}
}

// FetchResultsFromExternalAPI obtiene los resultados de una API externa y los inserta o actualiza en la base de datos
func (s *resultService) FetchResultsFromExternalAPI(ctx context.Context, sessionID int) ([]dto.ResponseResultDTO, e.ApiError) {
	// 1. Obtener sessionKey llamando al otro microservicio
	sessionKey, err := s.sessionsClient.GetSessionKeyBySessionID(sessionID)
	if err != nil {
		return nil, e.NewInternalServerApiError("Error obteniendo session key", err)
	}
	fmt.Println("Session Key obtenida:", sessionKey)

	// 2. Obtener las "positions" desde la API externa
	positions, err := s.externalClient.GetPositions(sessionKey)
	if err != nil {
		return nil, e.NewInternalServerApiError("Error fetching positions from external API", err)
	}

	// Mapa para quedarnos con la última posición reportada (driverNumber -> *int)
	finalPositions := make(map[int]*int)

	// Agrupar posiciones por driverNumber y ordenar por date para obtener la última posición
	positionsByDriver := make(map[int][]dto.Position)
	for _, pos := range positions {
		positionsByDriver[pos.DriverNumber] = append(positionsByDriver[pos.DriverNumber], pos)
	}

	// Ordenar posiciones por date y tomar la última para cada piloto
	for driverNumber, driverPositions := range positionsByDriver {
		sort.Slice(driverPositions, func(i, j int) bool {
			return driverPositions[i].Date < driverPositions[j].Date
		})
		if len(driverPositions) > 0 {
			finalPositions[driverNumber] = driverPositions[len(driverPositions)-1].Position
		}
	}

	// 3. Encontrar al piloto en posición 1 para determinar el número total de vueltas
	var position1DriverNumber int
	for driverNumber, position := range finalPositions {
		if position != nil && *position == 1 {
			position1DriverNumber = driverNumber
			break
		}
	}

	// Si no encontramos un piloto en posición 1, no podemos determinar las vueltas totales
	if position1DriverNumber == 0 {
		return nil, e.NewInternalServerApiError("No se encontró un piloto en posición 1 para determinar las vueltas totales", nil)
	}

	// Obtener las vueltas del piloto en posición 1
	position1Laps, err := s.externalClient.GetLaps(sessionKey, position1DriverNumber)
	if err != nil {
		return nil, e.NewInternalServerApiError(fmt.Sprintf("Error obteniendo vueltas del piloto en posición 1 (driver %d): %v", position1DriverNumber, err), err)
	}

	totalSessionLaps := len(position1Laps)

	var responseResults []dto.ResponseResultDTO

	// 4. Para cada driverNumber, determinamos la vuelta más rápida y actualizamos/insertamos en DB
	for driverNumber, positionPtr := range finalPositions {
		laps, err := s.externalClient.GetLaps(sessionKey, driverNumber)
		if err != nil {
			fmt.Printf("Error obteniendo vueltas para driver %d: %v\n", driverNumber, err)
			continue
		}

		// 5. Encontrar la vuelta más rápida o asignar 0 si no hay vueltas
		var fastestLap float64
		if len(laps) == 0 {
			fastestLap = 0
		} else {
			for _, lap := range laps {
				if fastestLap == 0 || lap.LapDuration < fastestLap {
					fastestLap = lap.LapDuration
				}
			}
		}

		// 6. Determinar el status basado en las vueltas
		var status string
		if len(laps) == 0 {
			status = "DNS" // Did Not Start
		} else if len(laps) < totalSessionLaps {
			status = "DNF" // Did Not Finish
		} else {
			status = "FINISHED" // Completó todas las vueltas
		}

		// 7. Obtener info completa del driver desde microservicio de drivers
		driverInfo, err := s.driversClient.GetDriverByNumber(driverNumber)
		if err != nil {
			fmt.Printf("Error obteniendo info del driver_number %d: %v\n", driverNumber, err)
			continue
		}

		// 8. Ver si ya existe un resultado para (driver, session)
		existingResult, _ := s.resultRepo.GetResultByDriverAndSession(ctx, driverInfo.ID, sessionID)

		if existingResult == nil {
			// Crear nuevo result
			newResult := &model.Result{
				SessionID:      sessionID,
				DriverID:       driverInfo.ID,
				Position:       positionPtr,
				Status:         status,
				FastestLapTime: fastestLap,
			}
			if err := s.resultRepo.CreateResult(ctx, newResult); err != nil {
				return nil, e.NewInternalServerApiError("Error inserting result", err)
			}
			existingResult = newResult
		} else {
			// Actualizar
			existingResult.Position = positionPtr
			existingResult.Status = status
			existingResult.FastestLapTime = fastestLap
			if err := s.resultRepo.UpdateResult(ctx, existingResult); err != nil {
				return nil, e.NewInternalServerApiError("Error updating existing result", err)
			}
		}

		// 9. Obtener la info de la sesión para armar el DTO de respuesta
		sessionData, err := s.sessionsClient.GetSessionByID(sessionID)
		if err != nil {
			return nil, e.NewInternalServerApiError("Error fetching session data", err)
		}

		// 10. Construir el DTO
		responseResult := dto.ResponseResultDTO{
			ID:             existingResult.ID,
			Position:       existingResult.Position,
			Status:         existingResult.Status,
			FastestLapTime: existingResult.FastestLapTime,
			Driver: dto.ResponseDriverDTO{
				ID:          driverInfo.ID,
				FirstName:   driverInfo.FirstName,
				LastName:    driverInfo.LastName,
				FullName:    driverInfo.FullName,
				NameAcronym: driverInfo.NameAcronym,
				TeamName:    driverInfo.TeamName,
			},
			Session: dto.ResponseSessionDTO{
				ID:               sessionData.ID,
				CircuitShortName: sessionData.CircuitShortName,
				CountryName:      sessionData.CountryName,
				Location:         sessionData.Location,
				SessionName:      sessionData.SessionName,
				SessionType:      sessionData.SessionType,
				DateStart:        sessionData.DateStart,
			},
			CreatedAt: existingResult.CreatedAt,
			UpdatedAt: existingResult.UpdatedAt,
		}
		responseResults = append(responseResults, responseResult)
	}

	return responseResults, nil
}

func (s *resultService) FetchNonRaceSessionResults(ctx context.Context, sessionId int) ([]dto.ResponseResultDTO, e.ApiError) {
	// 1. Obtener información de la sesión para verificar que no es Race
	sessionData, err := s.sessionsClient.GetSessionByID(sessionId)
	if err != nil {
		return nil, e.NewInternalServerApiError("Error fetching session data", err)
	}

	// Verificar que session_name y session_type no sean "Race"
	if strings.ToLower(sessionData.SessionName) == "race" || strings.ToLower(sessionData.SessionType) == "race" {
		return nil, e.NewBadRequestApiError("Este endpoint solo soporta sesiones no Race")
	}

	// 2. Obtener sessionKey llamando al microservicio de sessions
	sessionKey, err := s.sessionsClient.GetSessionKeyBySessionID(sessionId)
	if err != nil {
		return nil, e.NewInternalServerApiError("Error obteniendo session key", err)
	}

	// 3. Obtener las posiciones desde la API externa
	positions, err := s.externalClient.GetPositions(sessionKey)
	if err != nil {
		return nil, e.NewInternalServerApiError("Error fetching positions from external API", err)
	}

	// Mapa para quedarnos con la última posición reportada (driverNumber -> *int)
	finalPositions := make(map[int]*int)

	// Agrupar posiciones por driverNumber y ordenar por date para obtener la última posición
	positionsByDriver := make(map[int][]dto.Position)
	for _, pos := range positions {
		positionsByDriver[pos.DriverNumber] = append(positionsByDriver[pos.DriverNumber], pos)
	}

	// Ordenar posiciones por date y tomar la última para cada piloto
	for driverNumber, driverPositions := range positionsByDriver {
		sort.Slice(driverPositions, func(i, j int) bool {
			return driverPositions[i].Date < driverPositions[j].Date
		})
		if len(driverPositions) > 0 {
			finalPositions[driverNumber] = driverPositions[len(driverPositions)-1].Position
		}
	}

	var responseResults []dto.ResponseResultDTO

	// 4. Procesar cada piloto y sus posiciones
	for driverNumber, positionPtr := range finalPositions {
		// Obtener info completa del driver desde microservicio de drivers
		driverInfo, err := s.driversClient.GetDriverByNumber(driverNumber)
		if err != nil {
			fmt.Printf("Error obteniendo info del driver_number %d: %v\n", driverNumber, err)
			continue
		}

		// Verificar si ya existe un resultado para (driver, session)
		existingResult, _ := s.resultRepo.GetResultByDriverAndSession(ctx, driverInfo.ID, sessionId)

		// Status por defecto para sesiones no Race
		status := "FINISHED"

		if existingResult == nil {
			newResult := &model.Result{
				SessionID:      sessionId,
				DriverID:       driverInfo.ID,
				Position:       positionPtr,
				Status:         status,
				FastestLapTime: 0, // No se consideran vueltas
			}
			if err := s.resultRepo.CreateResult(ctx, newResult); err != nil {
				return nil, e.NewInternalServerApiError("Error inserting result", err)
			}
			existingResult = newResult
		} else {
			// Actualizar
			existingResult.Position = positionPtr
			existingResult.Status = status
			existingResult.FastestLapTime = 0 // No se consideran vueltas
			if err := s.resultRepo.UpdateResult(ctx, existingResult); err != nil {
				return nil, e.NewInternalServerApiError("Error updating existing result", err)
			}
		}

		// Construir el DTO
		responseResult := dto.ResponseResultDTO{
			ID:             existingResult.ID,
			Position:       existingResult.Position,
			Status:         existingResult.Status,
			FastestLapTime: existingResult.FastestLapTime,
			Driver: dto.ResponseDriverDTO{
				ID:          driverInfo.ID,
				FirstName:   driverInfo.FirstName,
				LastName:    driverInfo.LastName,
				FullName:    driverInfo.FullName,
				NameAcronym: driverInfo.NameAcronym,
				TeamName:    driverInfo.TeamName,
			},
			Session: dto.ResponseSessionDTO{
				ID:               sessionData.ID,
				CircuitShortName: sessionData.CircuitShortName,
				CountryName:      sessionData.CountryName,
				Location:         sessionData.Location,
				SessionName:      sessionData.SessionName,
				SessionType:      sessionData.SessionType,
				DateStart:        sessionData.DateStart,
			},
			CreatedAt: existingResult.CreatedAt,
			UpdatedAt: existingResult.UpdatedAt,
		}
		responseResults = append(responseResults, responseResult)
	}

	return responseResults, nil
}

// ESTO SOLO SIRVE PARA CREAR UN RESULTADO A LA VEZ
// CreateResult crea un nuevo resultado
func (s *resultService) CreateResult(ctx context.Context, request dto.CreateResultDTO) (dto.ResponseResultDTO, e.ApiError) {
	// 1. Validar status
	//    Ejemplo de validación de status vs. position
	validStatuses := map[string]bool{"FINISHED": true, "DNF": true, "DNS": true, "DSQ": true}
	if request.Status == "" {
		// Si no viene, por defecto interpretamos FINISHED si hay position, DNF si no
		if request.Position != nil {
			request.Status = "FINISHED"
		} else {
			request.Status = "DNF"
		}
	} else {
		// Verificar si el status es uno de los válidos
		if !validStatuses[request.Status] {
			return dto.ResponseResultDTO{}, e.NewBadRequestApiError(fmt.Sprintf("Status inválido: %s", request.Status))
		}
	}

	// 2. Reglas:
	//    - Si status == "FINISHED", position != nil y entre 1..20
	//    - Si status != "FINISHED", position == nil
	if request.Status == "FINISHED" {
		if request.Position == nil {
			return dto.ResponseResultDTO{}, e.NewBadRequestApiError("Debe proporcionar una posición cuando el estado es FINISHED")
		}
		if *request.Position < 1 || *request.Position > 20 {
			return dto.ResponseResultDTO{}, e.NewBadRequestApiError("La posición debe estar entre 1 y 20 si es FINISHED")
		}
	} else {
		// DNF, DNS, DSQ => position debe ser nil
		if request.Position != nil {
			return dto.ResponseResultDTO{}, e.NewBadRequestApiError(fmt.Sprintf("No puede proporcionar Position cuando Status es %s", request.Status))
		}
	}

	// 3. Validar fastestLapTime si lo deseas
	if request.FastestLapTime != 0 && request.FastestLapTime < 30 {
		return dto.ResponseResultDTO{}, e.NewBadRequestApiError("Fastest lap time debe ser mayor a 30 (o 0 si se omite)")
	}

	// 4. Revisar si ya existe un resultado para (driver, session)
	existingResult, _ := s.resultRepo.GetResultByDriverAndSession(ctx, request.DriverID, request.SessionID)
	if existingResult != nil {
		return dto.ResponseResultDTO{}, e.NewBadRequestApiError("Ya existe un resultado para este driver en esta sesión")
	}

	// 5. Crear modelo
	newResult := &model.Result{
		SessionID:      request.SessionID,
		DriverID:       request.DriverID,
		Position:       request.Position,
		Status:         request.Status,
		FastestLapTime: request.FastestLapTime,
	}

	// 6. Insertar en DB
	if err := s.resultRepo.CreateResult(ctx, newResult); err != nil {
		return dto.ResponseResultDTO{}, e.NewInternalServerApiError("Error creando resultado", err)
	}

	// 8. Construir DTO de respuesta
	response := dto.ResponseResultDTO{
		ID:             newResult.ID,
		Position:       newResult.Position,
		Status:         newResult.Status,
		FastestLapTime: newResult.FastestLapTime,
		Driver:         dto.ResponseDriverDTO{ID: newResult.DriverID},
		Session:        dto.ResponseSessionDTO{ID: newResult.SessionID},
		CreatedAt:      newResult.CreatedAt,
		UpdatedAt:      newResult.UpdatedAt,
	}

	return response, nil
}

// UpdateResult actualiza un resultado existente
func (s *resultService) UpdateResult(ctx context.Context, resultID int, request dto.UpdateResultDTO) (dto.ResponseResultDTO, e.ApiError) {
	// 1. Buscar el resultado en DB
	result, err := s.resultRepo.GetResultByID(ctx, resultID)
	if err != nil {
		return dto.ResponseResultDTO{}, e.NewBadRequestApiError("Error obteniendo el resultado por su ID")
	}

	// 2. Actualizar STATUS
	validStatuses := map[string]bool{"FINISHED": true, "DNF": true, "DNS": true, "DSQ": true}
	if request.Status != "" {
		// Si viene un nuevo Status, validarlo
		if !validStatuses[request.Status] {
			return dto.ResponseResultDTO{}, e.NewBadRequestApiError(fmt.Sprintf("Status inválido: %s", request.Status))
		}
		result.Status = request.Status
	}

	// 3. Actualizar POSITION si viene
	if request.Position != nil {
		// Si la nueva position no es nil, forzamos status = FINISHED
		if result.Status != "" && result.Status != "FINISHED" {
			return dto.ResponseResultDTO{}, e.NewBadRequestApiError(
				fmt.Sprintf("No se puede asignar Position si el Status es %s", result.Status),
			)
		}
		if *request.Position < 1 || *request.Position > 20 {
			return dto.ResponseResultDTO{}, e.NewBadRequestApiError("La posición debe estar entre 1 y 20")
		}
		// Marcamos status "FINISHED" si no se había puesto
		if result.Status == "" || result.Status == "DNF" || result.Status == "DNS" || result.Status == "DSQ" {
			result.Status = "FINISHED"
		}
		result.Position = request.Position
	}

	// 4. Actualizar fastestLapTime si != 0
	if request.FastestLapTime != 0 {
		if request.FastestLapTime < 30 {
			return dto.ResponseResultDTO{}, e.NewBadRequestApiError("Invalid fastest lap time, must be > 30")
		}
		result.FastestLapTime = request.FastestLapTime
	}

	// 5. Persistir cambios
	if err := s.resultRepo.UpdateResult(ctx, result); err != nil {
		return dto.ResponseResultDTO{}, e.NewInternalServerApiError("Error updating result", err)
	}

	// 7. Construir respuesta
	response := dto.ResponseResultDTO{
		ID:             result.ID,
		Position:       result.Position,
		Status:         result.Status,
		FastestLapTime: result.FastestLapTime,
		Driver: dto.ResponseDriverDTO{
			ID:          result.Driver.ID,
			FirstName:   result.Driver.FirstName,
			LastName:    result.Driver.LastName,
			FullName:    result.Driver.FullName,
			NameAcronym: result.Driver.NameAcronym,
			TeamName:    result.Driver.TeamName,
		},
		Session: dto.ResponseSessionDTO{
			ID:               result.Session.ID,
			CircuitShortName: result.Session.CircuitShortName,
			CountryName:      result.Session.CountryName,
			Location:         result.Session.Location,
			SessionName:      result.Session.SessionName,
			SessionType:      result.Session.SessionType,
			DateStart:        result.Session.DateStart,
		},
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}

	return response, nil
}

// GetResultsOrderedByPosition obtiene los resultados de una sesión específica ordenados por posición
func (s *resultService) GetResultsOrderedByPosition(ctx context.Context, sessionID int) ([]dto.ResponseResultDTO, e.ApiError) {
	// Verificar caché
	// cacheKey := fmt.Sprintf("results:session:%d", sessionID)
	// if cached, exists := s.cache.Get(cacheKey); exists {
	// 	if results, ok := cached.([]dto.ResponseResultDTO); ok {
	// 		fmt.Printf("Cache hit for results:session:%d\n", sessionID)
	// 		return results, nil
	// 	}
	// }

	// Verificar si existe el sessionID en la tabla de resultados
	exists, err := s.resultRepo.ExistsSessionInResults(ctx, sessionID)
	if err != nil {
		return nil, e.NewInternalServerApiError("Error verifying session existence in results", err)
	}
	if !exists {
		return nil, e.NewNotFoundApiError("No results found for the given session ID")
	}

	// Obtener los resultados ordenados por posición
	results, err := s.resultRepo.GetResultsOrderedByPosition(ctx, sessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundApiError("No results found for this session")
		}
		return nil, e.NewInternalServerApiError("Error retrieving results", err)
	}

	var responseResults []dto.ResponseResultDTO
	for _, result := range results {
		response := dto.ResponseResultDTO{
			ID:             result.ID,
			Position:       result.Position,
			Status:         result.Status,
			FastestLapTime: result.FastestLapTime,
			Driver: dto.ResponseDriverDTO{
				ID:            result.Driver.ID,
				BroadcastName: result.Driver.BroadcastName,
				CountryCode:   result.Driver.CountryCode,
				DriverNumber:  result.Driver.DriverNumber,
				FirstName:     result.Driver.FirstName,
				LastName:      result.Driver.LastName,
				FullName:      result.Driver.FullName,
				NameAcronym:   result.Driver.NameAcronym,
				TeamName:      result.Driver.TeamName,
			},
			Session: dto.ResponseSessionDTO{
				ID:               result.Session.ID,
				CircuitShortName: result.Session.CircuitShortName,
				CountryName:      result.Session.CountryName,
				Location:         result.Session.Location,
				SessionName:      result.Session.SessionName,
				SessionType:      result.Session.SessionType,
				DateStart:        result.Session.DateStart,
			},
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt,
		}
		responseResults = append(responseResults, response)
	}

	// // Cachear el resultado
	// ttl := 5 * time.Minute
	// if len(results) > 0 && results[0].Session.DateEnd.Before(time.Now()) {
	// 	ttl = 24 * time.Hour // Sesiones finalizadas son inmutables
	// }
	// s.cache.Set(cacheKey, responseResults, ttl)
	// fmt.Printf("Cached results for session:%d\n", sessionID)

	return responseResults, nil
}

// GetFastestLapInSession obtiene el piloto con la vuelta más rápida en una sesión específica
func (s *resultService) GetFastestLapInSession(ctx context.Context, sessionID int) (dto.ResponseResultDTO, e.ApiError) {
	// Verificar caché
	// cacheKey := fmt.Sprintf("fastest_lap:session:%d", sessionID)
	// if cached, exists := s.cache.Get(cacheKey); exists {
	// 	if result, ok := cached.(dto.ResponseResultDTO); ok {
	// 		fmt.Printf("Cache hit for fastest_lap:session:%d\n", sessionID)
	// 		return result, nil
	// 	}
	// }

	// Verificar si existe el sessionID en la tabla de resultados
	exists, err := s.resultRepo.ExistsSessionInResults(ctx, sessionID)
	if err != nil {
		return dto.ResponseResultDTO{}, e.NewInternalServerApiError("Error verifying session existence in results", err)
	}
	if !exists {
		return dto.ResponseResultDTO{}, e.NewNotFoundApiError("No results found for the given session ID")
	}

	// Obtener la vuelta más rápida de la sesión
	results, err := s.resultRepo.GetResultsBySessionID(ctx, sessionID)
	if err != nil {
		return dto.ResponseResultDTO{}, e.NewInternalServerApiError("Error fetching session results", err)
	}

	var fastestResult *model.Result
	for _, result := range results {
		if result.FastestLapTime > 0 {
			if fastestResult == nil || result.FastestLapTime < fastestResult.FastestLapTime {
				fastestResult = result
			}
		}
	}

	if fastestResult == nil {
		return dto.ResponseResultDTO{}, e.NewNotFoundApiError("No valid lap times found for the given session")
	}

	// Convertir el resultado más rápido a DTO
	response := dto.ResponseResultDTO{
		ID:             fastestResult.ID,
		Position:       fastestResult.Position,
		FastestLapTime: fastestResult.FastestLapTime,
		Driver: dto.ResponseDriverDTO{
			ID:          fastestResult.Driver.ID,
			FirstName:   fastestResult.Driver.FirstName,
			LastName:    fastestResult.Driver.LastName,
			FullName:    fastestResult.Driver.FullName,
			NameAcronym: fastestResult.Driver.NameAcronym,
			TeamName:    fastestResult.Driver.TeamName,
		},
		Session: dto.ResponseSessionDTO{
			ID:               fastestResult.Session.ID,
			CircuitShortName: fastestResult.Session.CircuitShortName,
			CountryName:      fastestResult.Session.CountryName,
			Location:         fastestResult.Session.Location,
			SessionName:      fastestResult.Session.SessionName,
			SessionType:      fastestResult.Session.SessionType,
			DateStart:        fastestResult.Session.DateStart,
		},
		CreatedAt: fastestResult.CreatedAt,
		UpdatedAt: fastestResult.UpdatedAt,
	}

	// Cachear el resultado
	// ttl := 5 * time.Minute
	// if fastestResult.Session.DateEnd.Before(time.Now()) {
	// 	ttl = 24 * time.Hour
	// }
	// s.cache.Set(cacheKey, response, ttl)
	// fmt.Printf("Cached fastest_lap for session:%d\n", sessionID)

	return response, nil
}

// DeleteResult elimina un resultado específico
func (s *resultService) DeleteResult(ctx context.Context, resultID int) e.ApiError {
	if resultID == 0 {
		return e.NewBadRequestApiError("El ID del resultado no puede ser 0")
	}

	result, err := s.resultRepo.GetResultByID(ctx, resultID)
	if err != nil {
		if err == e.NewNotFoundApiError("Result not found") {
			return e.NewNotFoundApiError("El resultado con el ID proporcionado no existe, no se puede eliminar")
		}
		return e.NewInternalServerApiError("Error al verificar la existencia del resultado", err)
	}

	fmt.Printf("Eliminando resultado: ID=%d, DriverID=%d, SessionID=%d\n", result.ID, result.DriverID, result.SessionID)

	if err := s.resultRepo.DeleteResult(ctx, resultID); err != nil {
		return e.NewInternalServerApiError("Error al eliminar el resultado", err)
	}

	// // Invalidar caché relevante
	// cacheKeys := []string{
	// 	fmt.Sprintf("results:session:%d", result.SessionID),
	// 	fmt.Sprintf("fastest_lap:session:%d", result.SessionID),
	// 	fmt.Sprintf("top_drivers:session:%d:n:%d", result.SessionID, 20),
	// 	"all_results",
	// }
	// for _, key := range cacheKeys {
	// 	s.cache.Delete(key)
	// 	fmt.Printf("Invalidated cache for key=%s\n", key)
	// }

	return nil
}

// DeleteAllResultsForSession elimina todos los resultados asociados a una sesión específica
func (s *resultService) DeleteAllResultsForSession(ctx context.Context, sessionID int) e.ApiError {
	if sessionID == 0 {
		return e.NewBadRequestApiError("El ID de la sesión no puede ser 0")
	}

	results, err := s.resultRepo.GetResultsBySessionID(ctx, sessionID)
	if err != nil {
		return e.NewInternalServerApiError("Error al obtener los resultados de la sesión", err)
	}

	if len(results) == 0 {
		return e.NewNotFoundApiError("No se encontraron resultados para la sesión especificada")
	}

	for _, result := range results {
		if err := s.resultRepo.DeleteResult(ctx, result.ID); err != nil {
			return e.NewInternalServerApiError(fmt.Sprintf("Error al eliminar el resultado con ID %d", result.ID), err)
		}
	}

	// // Invalidar caché relevante
	// cacheKeys := []string{
	// 	fmt.Sprintf("results:session:%d", sessionID),
	// 	fmt.Sprintf("fastest_lap:session:%d", sessionID),
	// 	fmt.Sprintf("top_drivers:session:%d:n:%d", sessionID, 20),
	// 	"all_results",
	// }
	// for _, key := range cacheKeys {
	// 	s.cache.Delete(key)
	// 	fmt.Printf("Invalidated cache for key=%s\n", key)
	// }

	return nil
}

// GetAllResults obtiene todos los resultados de la base de datos
func (s *resultService) GetAllResults(ctx context.Context) ([]dto.ResponseResultDTO, e.ApiError) {
	// // Verificar caché
	// cacheKey := "all_results"
	// if cached, exists := s.cache.Get(cacheKey); exists {
	// 	if results, ok := cached.([]dto.ResponseResultDTO); ok {
	// 		fmt.Printf("Cache hit for all_results\n")
	// 		return results, nil
	// 	}
	// }

	results, err := s.resultRepo.GetAllResults(ctx)
	if err != nil {
		return nil, e.NewInternalServerApiError("Error al obtener todos los resultados", err)
	}

	if len(results) == 0 {
		return nil, e.NewNotFoundApiError("No se encontraron resultados en la base de datos")
	}

	var responseResults []dto.ResponseResultDTO
	for _, result := range results {
		response := dto.ResponseResultDTO{
			ID:             result.ID,
			Position:       result.Position,
			FastestLapTime: result.FastestLapTime,
			Driver: dto.ResponseDriverDTO{
				ID:          result.Driver.ID,
				FirstName:   result.Driver.FirstName,
				LastName:    result.Driver.LastName,
				FullName:    result.Driver.FullName,
				NameAcronym: result.Driver.NameAcronym,
				TeamName:    result.Driver.TeamName,
			},
			Session: dto.ResponseSessionDTO{
				ID:               result.Session.ID,
				CircuitShortName: result.Session.CircuitShortName,
				CountryName:      result.Session.CountryName,
				Location:         result.Session.Location,
				SessionName:      result.Session.SessionName,
				SessionType:      result.Session.SessionType,
				DateStart:        result.Session.DateStart,
			},
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt,
		}
		responseResults = append(responseResults, response)
	}

	// // Cachear el resultado
	// s.cache.Set(cacheKey, responseResults, 5*time.Minute)
	// fmt.Printf("Cached all_results\n")

	return responseResults, nil
}

// GetTopNDriversInSession obtiene los mejores N pilotos de una sesión específica.
func (s *resultService) GetTopNDriversInSession(ctx context.Context, sessionID int, n int) ([]dto.TopDriverDTO, e.ApiError) {
	if sessionID == 0 {
		return nil, e.NewBadRequestApiError("El ID de la sesión no puede ser 0")
	}
	if n < 1 {
		return nil, e.NewBadRequestApiError("El número de pilotos a obtener debe ser mayor que 0")
	}

	// // Verificar caché
	// cacheKey := fmt.Sprintf("top_drivers:session:%d:n:%d", sessionID, n)
	// if cached, exists := s.cache.Get(cacheKey); exists {
	// 	if topDrivers, ok := cached.([]dto.TopDriverDTO); ok {
	// 		fmt.Printf("Cache hit for top_drivers:session:%d:n:%d\n", sessionID, n)
	// 		return topDrivers, nil
	// 	}
	// }

	results, err := s.resultRepo.GetResultsOrderedByPosition(ctx, sessionID)
	if err != nil {
		return nil, e.NewInternalServerApiError("Error obteniendo resultados de la sesión", err)
	}

	if len(results) == 0 {
		return nil, e.NewNotFoundApiError("No se encontraron resultados para la sesión")
	}

	var finishedResults []*model.Result
	for _, r := range results {
		if r.Position != nil {
			finishedResults = append(finishedResults, r)
		}
	}
	if len(finishedResults) == 0 {
		return nil, e.NewNotFoundApiError("Ningún piloto terminó la sesión")
	}

	if n > len(finishedResults) {
		n = len(finishedResults)
	}
	if n > 20 {
		n = 20
	}

	var topDrivers []dto.TopDriverDTO
	for i := 0; i < n; i++ {
		topDrivers = append(topDrivers, dto.TopDriverDTO{
			Position: *finishedResults[i].Position,
			DriverID: finishedResults[i].DriverID,
		})
	}

	// Cachear el resultado
	// ttl := 5 * time.Minute
	// if len(results) > 0 && results[0].Session.DateEnd.Before(time.Now()) {
	// 	ttl = 24 * time.Hour
	// }
	// s.cache.Set(cacheKey, topDrivers, ttl)
	// fmt.Printf("Cached top_drivers for session:%d:n:%d\n", sessionID, n)

	return topDrivers, nil
}

func (s *resultService) CreateSessionResultsAdmin(ctx context.Context, bulkRequest dto.CreateBulkResultsDTO) ([]dto.ResponseResultDTO, e.ApiError) {
	if bulkRequest.SessionID == 0 {
		return nil, e.NewBadRequestApiError("El session_id no puede ser 0")
	}

	var resultsToCreate []*model.Result
	var resultsToUpdate []*model.Result

	validStatuses := map[string]bool{"FINISHED": true, "DNF": true, "DNS": true, "DSQ": true}

	existingResults, err := s.resultRepo.GetResultsBySessionID(ctx, bulkRequest.SessionID)
	if err != nil {
		return nil, e.NewInternalServerApiError("Error obteniendo resultados existentes", err)
	}

	existingResultsMap := make(map[int]*model.Result)
	for _, r := range existingResults {
		existingResultsMap[r.DriverID] = r
	}

	for _, item := range bulkRequest.Results {
		if item.Status == "" {
			if item.Position != nil {
				item.Status = "FINISHED"
			} else {
				item.Status = "DNF"
			}
		} else {
			if !validStatuses[item.Status] {
				return nil, e.NewBadRequestApiError(fmt.Sprintf("Status inválido: %s", item.Status))
			}
		}

		if item.Status == "FINISHED" {
			if item.Position == nil {
				return nil, e.NewBadRequestApiError("Debe proporcionar una posición si el status es FINISHED")
			}
			if *item.Position < 1 || *item.Position > 20 {
				return nil, e.NewBadRequestApiError(
					fmt.Sprintf("Posición inválida para driver_id %d. Debe estar entre 1 y 20", item.DriverID),
				)
			}
		} else {
			if item.Position != nil {
				return nil, e.NewBadRequestApiError(
					fmt.Sprintf("No puede dar Position si el status es %s (driver_id %d)", item.Status, item.DriverID),
				)
			}
		}

		if item.FastestLapTime != 0 && item.FastestLapTime < 30 {
			return nil, e.NewBadRequestApiError(
				fmt.Sprintf("FastestLapTime inválido para driver_id %d. Debe ser >30 o 0", item.DriverID),
			)
		}

		if existingResult, exists := existingResultsMap[item.DriverID]; exists {
			existingResult.Position = item.Position
			existingResult.Status = item.Status
			existingResult.FastestLapTime = item.FastestLapTime
			resultsToUpdate = append(resultsToUpdate, existingResult)
		} else {
			newResult := &model.Result{
				SessionID:      bulkRequest.SessionID,
				DriverID:       item.DriverID,
				Position:       item.Position,
				Status:         item.Status,
				FastestLapTime: item.FastestLapTime,
			}
			resultsToCreate = append(resultsToCreate, newResult)
		}
	}

	txErr := s.resultRepo.SessionCreateOrUpdateResultsAdmin(ctx, resultsToCreate, resultsToUpdate)
	if txErr != nil {
		return nil, e.NewInternalServerApiError("Error creando o actualizando resultados masivamente", txErr)
	}

	updatedResults, err := s.resultRepo.GetResultsBySessionID(ctx, bulkRequest.SessionID)
	if err != nil {
		return nil, e.NewInternalServerApiError("Error obteniendo resultados actualizados", err)
	}

	// Invalidar caché relevante
	// cacheKeys := []string{
	// 	fmt.Sprintf("results:session:%d", bulkRequest.SessionID),
	// 	fmt.Sprintf("fastest_lap:session:%d", bulkRequest.SessionID),
	// 	fmt.Sprintf("top_drivers:session:%d:n:%d", bulkRequest.SessionID, 20),
	// 	"all_results",
	// }
	// for _, key := range cacheKeys {
	// 	s.cache.Delete(key)
	// 	fmt.Printf("Invalidated cache for key=%s\n", key)
	// }

	var responseResults []dto.ResponseResultDTO
	for _, r := range updatedResults {
		responseResults = append(responseResults, dto.ResponseResultDTO{
			ID:             r.ID,
			Position:       r.Position,
			Status:         r.Status,
			FastestLapTime: r.FastestLapTime,
			Driver:         dto.ResponseDriverDTO{ID: r.DriverID},
			Session:        dto.ResponseSessionDTO{ID: r.SessionID},
			CreatedAt:      r.CreatedAt,
			UpdatedAt:      r.UpdatedAt,
		})
	}

	return responseResults, nil
}
