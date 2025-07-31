package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	model "prediapp.local/db/model"
	client "prediapp.local/prodes/internal/client"
	prodes "prediapp.local/prodes/internal/dto"
	repository "prediapp.local/prodes/internal/repository"
	e "prediapp.local/prodes/pkg/utils"

	"gorm.io/gorm"
)

type prodeService struct {
	prodeRepo     repository.ProdeRepository
	sessionClient *client.HttpClient
	userClient    *client.HttpClient
	driverClient  *client.HttpClient
	resultsClient *client.HttpClient
	// cache         *e.Cache
}

type ProdeServiceInterface interface {
	CreateProdeCarrera(ctx context.Context, request prodes.CreateProdeCarreraDTO) (prodes.ResponseProdeCarreraDTO, e.ApiError)
	CreateProdeSession(ctx context.Context, request prodes.CreateProdeSessionDTO) (prodes.ResponseProdeSessionDTO, e.ApiError)
	UpdateProdeCarrera(ctx context.Context, request prodes.UpdateProdeCarreraDTO) (prodes.ResponseProdeCarreraDTO, e.ApiError)
	UpdateProdeSession(ctx context.Context, request prodes.UpdateProdeSessionDTO) (prodes.ResponseProdeSessionDTO, e.ApiError)
	DeleteProdeById(ctx context.Context, prodeID int) e.ApiError
	GetProdesByUserId(ctx context.Context, userID int) ([]prodes.ResponseProdeCarreraDTO, []prodes.ResponseProdeSessionDTO, e.ApiError)
	GetRaceProdesBySession(ctx context.Context, sessionID int) ([]prodes.ResponseProdeCarreraDTO, e.ApiError)
	UpdateRaceProdeForUserBySessionId(ctx context.Context, userID int, sessionID int, updatedProde prodes.UpdateProdeCarreraDTO) (prodes.ResponseProdeCarreraDTO, e.ApiError)
	GetSessionProdeBySession(ctx context.Context, sessionID int) ([]prodes.ResponseProdeSessionDTO, e.ApiError)
	GetUserProdes(ctx context.Context, userID int) ([]prodes.ResponseProdeCarreraDTO, []prodes.ResponseProdeSessionDTO, e.ApiError)
	GetDriverDetails(ctx context.Context, driverID int) (prodes.DriverDTO, e.ApiError)
	GetAllDrivers(ctx context.Context) ([]prodes.DriverDTO, e.ApiError)
	GetTopDriversBySessionId(ctx context.Context, sessionID, n int) ([]prodes.TopDriverDTO, e.ApiError)
	GetProdeByUserAndSession(ctx context.Context, userID int, sessionID int) (*prodes.ResponseProdeCarreraDTO, *prodes.ResponseProdeSessionDTO, e.ApiError)
	UpdateScoresForRaceProdes(ctx context.Context, sessionID int) e.ApiError
	UpdateScoresForSessionProdes(ctx context.Context, sessionID int) e.ApiError
	// UpdateUserScores(ctx context.Context) e.ApiError
}

// NewProdeService crea una nueva instancia de ProdeService con inyección de dependencias
func NewPrediService(prodeRepo repository.ProdeRepository, sessionClient *client.HttpClient, userClient *client.HttpClient, driverClient *client.HttpClient, resultsClient *client.HttpClient) ProdeServiceInterface {
	return &prodeService{
		prodeRepo:     prodeRepo,
		sessionClient: sessionClient,
		userClient:    userClient,
		driverClient:  driverClient,
		resultsClient: resultsClient,
		// cache:         cache,
	}
}

func (s *prodeService) CreateProdeCarrera(ctx context.Context, request prodes.CreateProdeCarreraDTO) (prodes.ResponseProdeCarreraDTO, e.ApiError) {
	existingProde, err := s.prodeRepo.GetProdeCarreraBySessionIdAndUserId(ctx, request.UserID, request.SessionID)

	if err == nil {
		// Si ya existe un ProdeCarrera, actualizarlo en lugar de crear uno nuevo
		updateRequest := prodes.UpdateProdeCarreraDTO{
			ProdeID:   existingProde.ID,
			UserID:    existingProde.UserID,
			SessionID: existingProde.SessionID,
			P1:        request.P1,
			P2:        request.P2,
			P3:        request.P3,
			P4:        request.P4,
			P5:        request.P5,
			// FastestLap: request.FastestLap,
			VSC: request.VSC,
			SC:  request.SC,
			DNF: request.DNF,
		}
		return s.UpdateProdeCarrera(ctx, updateRequest)
	}

	// Hacer la llamada al cliente HTTP para obtener la información de la sesión
	sessionInfo, err := s.sessionClient.GetSessionNameAndType(request.SessionID)
	if err != nil {
		return prodes.ResponseProdeCarreraDTO{}, e.NewInternalServerApiError("Error fetching session name and type from sessions service", err)
	}

	// Validar tanto el session_name como el session_type
	if !isRaceSession(sessionInfo.SessionName, sessionInfo.SessionType) {
		return prodes.ResponseProdeCarreraDTO{}, e.NewBadRequestApiError("La sesión asociada no es una carrera válida (Race), no se puede crear un ProdeCarrera")
	}

	// Convertir DTO a modelo
	prode := model.ProdeCarrera{
		UserID:    request.UserID,
		SessionID: request.SessionID,
		P1:        request.P1,
		P2:        request.P2,
		P3:        request.P3,
		P4:        request.P4,
		P5:        request.P5,
		// FastestLap: request.FastestLap,
		VSC:   request.VSC,
		SC:    request.SC,
		DNF:   request.DNF,
		Score: 0,
	}

	// Crear el pronóstico de carrera en la base de datos
	err = s.prodeRepo.CreateProdeCarrera(ctx, &prode)
	if err != nil {
		return prodes.ResponseProdeCarreraDTO{}, e.NewInternalServerApiError("Error creando el pronóstico de carrera", err)
	}

	// // Invalidar caché relevante
	// cacheKeys := []string{
	// 	fmt.Sprintf("prode:user:%d", request.UserID),
	// 	fmt.Sprintf("race_prodes:session:%d", request.SessionID),
	// 	fmt.Sprintf("prode:user:%d:session:%d", request.UserID, request.SessionID),
	// }
	// for _, key := range cacheKeys {
	// 	s.cache.Delete(key)
	// 	fmt.Printf("Invalidated cache for key=%s\n", key)
	// }

	// Convertir el modelo a DTO de respuesta
	response := prodes.ResponseProdeCarreraDTO{
		ID:        prode.ID,
		UserID:    prode.UserID,
		SessionID: prode.SessionID,
		P1:        prode.P1,
		P2:        prode.P2,
		P3:        prode.P3,
		P4:        prode.P4,
		P5:        prode.P5,
		// FastestLap: prode.FastestLap,
		VSC:   prode.VSC,
		SC:    prode.SC,
		DNF:   prode.DNF,
		Score: prode.Score,
	}

	return response, nil
}

func (s *prodeService) CreateProdeSession(ctx context.Context, request prodes.CreateProdeSessionDTO) (prodes.ResponseProdeSessionDTO, e.ApiError) {
	existingProde, err := s.prodeRepo.GetProdeSessionBySessionIdAndUserId(ctx, request.UserID, request.SessionID)

	if err == nil {
		// Si ya existe un ProdeSession, actualizarlo en lugar de crear uno nuevo
		updateRequest := prodes.UpdateProdeSessionDTO{
			ProdeID:   existingProde.ID,
			UserID:    existingProde.UserID,
			SessionID: existingProde.SessionID,
			P1:        request.P1,
			P2:        request.P2,
			P3:        request.P3,
		}
		return s.UpdateProdeSession(ctx, updateRequest)
	}

	// Obtener la información de la sesión desde el microservicio de sesiones
	sessionInfo, err := s.sessionClient.GetSessionNameAndType(request.SessionID)
	if err != nil {
		return prodes.ResponseProdeSessionDTO{}, e.NewInternalServerApiError("Error fetching session name and type from sessions service", err)
	}

	// Verificar si la sesión es de tipo "Race"
	if isRaceSession(sessionInfo.SessionName, sessionInfo.SessionType) {
		return prodes.ResponseProdeSessionDTO{}, e.NewBadRequestApiError("La sesión asociada no es una carrera válida (Race), no se puede crear un ProdeCarrera")
	}

	// Convertir DTO a modelo
	prode := model.ProdeSession{
		UserID:    request.UserID,
		SessionID: request.SessionID,
		P1:        request.P1,
		P2:        request.P2,
		P3:        request.P3,
		Score:     0,
	}

	// Crear el pronóstico de sesión en la base de datos
	err = s.prodeRepo.CreateProdeSession(ctx, &prode)
	if err != nil {
		return prodes.ResponseProdeSessionDTO{}, e.NewInternalServerApiError("Error creando el pronóstico de sesión", err)
	}

	// // Invalidar caché relevante
	// cacheKeys := []string{
	// 	fmt.Sprintf("prode:user:%d", request.UserID),
	// 	fmt.Sprintf("session_prodes:session:%d", request.SessionID),
	// 	fmt.Sprintf("prode:user:%d:session:%d", request.UserID, request.SessionID),
	// }
	// for _, key := range cacheKeys {
	// 	s.cache.Delete(key)
	// 	fmt.Printf("Invalidated cache for key=%s\n", key)
	// }

	// Convertir el modelo a DTO de respuesta
	response := prodes.ResponseProdeSessionDTO{
		ID:        prode.ID,
		UserID:    prode.UserID,
		SessionID: prode.SessionID,
		P1:        prode.P1,
		P2:        prode.P2,
		P3:        prode.P3,
		Score:     prode.Score,
	}

	return response, nil
}

func (s *prodeService) UpdateProdeCarrera(ctx context.Context, request prodes.UpdateProdeCarreraDTO) (prodes.ResponseProdeCarreraDTO, e.ApiError) {
	// Buscar el prode existente para obtener los valores originales de SessionID y UserID
	existingProde, err := s.prodeRepo.GetProdeCarreraByID(ctx, request.ProdeID)
	if err != nil {
		return prodes.ResponseProdeCarreraDTO{}, e.NewNotFoundApiError("El pronóstico de carrera no fue encontrado")
	}

	// Obtener los detalles de la sesión directamente del microservicio de sesiones
	sessionDetails, httpErr := s.sessionClient.GetSessionByID(existingProde.SessionID)
	if httpErr != nil {
		return prodes.ResponseProdeCarreraDTO{}, e.NewInternalServerApiError("Error fetching session details", httpErr)
	}

	// Validar si la sesión ya ha comenzado
	if time.Now().After(sessionDetails.DateStart) {
		return prodes.ResponseProdeCarreraDTO{}, e.NewForbiddenApiError("No se puede actualizar el pronóstico, la carrera ya ha comenzado.")
	}

	// Proceder con la actualización del ProdeCarrera
	// Aquí usamos los valores originales de SessionID y UserID para evitar cambios no permitidos
	prode := model.ProdeCarrera{
		ID:        existingProde.ID,
		UserID:    existingProde.UserID,    // Mantener el UserID original
		SessionID: existingProde.SessionID, // Mantener el SessionID original
		P1:        request.P1,
		P2:        request.P2,
		P3:        request.P3,
		P4:        request.P4,
		P5:        request.P5,
		// FastestLap: request.FastestLap,
		VSC:       request.VSC,
		SC:        request.SC,
		DNF:       request.DNF,
		CreatedAt: existingProde.CreatedAt,
		UpdatedAt: time.Now(),
	}

	err = s.prodeRepo.UpdateProdeCarrera(ctx, &prode)
	if err != nil {
		return prodes.ResponseProdeCarreraDTO{}, e.NewInternalServerApiError("Error actualizando el pronóstico de carrera", err)
	}

	// // Invalidar caché relevante
	// cacheKeys := []string{
	// 	fmt.Sprintf("prode:user:%d", prode.UserID),
	// 	fmt.Sprintf("race_prodes:session:%d", prode.SessionID),
	// 	fmt.Sprintf("prode:user:%d:session:%d", prode.UserID, prode.SessionID),
	// }
	// for _, key := range cacheKeys {
	// 	s.cache.Delete(key)
	// 	fmt.Printf("Invalidated cache for key=%s\n", key)
	// }

	response := prodes.ResponseProdeCarreraDTO{
		ID:        prode.ID,
		UserID:    prode.UserID,
		SessionID: prode.SessionID,
		P1:        prode.P1,
		P2:        prode.P2,
		P3:        prode.P3,
		P4:        prode.P4,
		P5:        prode.P5,
		// FastestLap: prode.FastestLap,
		VSC:   prode.VSC,
		SC:    prode.SC,
		DNF:   prode.DNF,
		Score: prode.Score,
	}

	return response, nil
}

func (s *prodeService) UpdateProdeSession(ctx context.Context, request prodes.UpdateProdeSessionDTO) (prodes.ResponseProdeSessionDTO, e.ApiError) {
	// Buscar el prode existente para obtener los valores originales de SessionID y UserID
	existingProde, err := s.prodeRepo.GetProdeSessionByID(ctx, request.ProdeID)
	if err != nil {
		return prodes.ResponseProdeSessionDTO{}, e.NewNotFoundApiError("El pronóstico de sesión no fue encontrado")
	}

	// Obtener los detalles de la sesión directamente desde el microservicio de sesiones
	sessionDetails, httpErr := s.sessionClient.GetSessionByID(existingProde.SessionID)
	if httpErr != nil {
		return prodes.ResponseProdeSessionDTO{}, e.NewInternalServerApiError("Error fetching session details", httpErr)
	}

	// Validar si la sesión ya ha comenzado
	if time.Now().After(sessionDetails.DateStart) {
		return prodes.ResponseProdeSessionDTO{}, e.NewForbiddenApiError("No se puede actualizar el pronóstico, la sesión ya ha comenzado.")
	}

	// Proceder con la actualización del ProdeSession
	// Usar los valores originales de SessionID y UserID
	prode := model.ProdeSession{
		ID:        existingProde.ID,
		UserID:    existingProde.UserID,    // Mantener el UserID original
		SessionID: existingProde.SessionID, // Mantener el SessionID original
		P1:        request.P1,
		P2:        request.P2,
		P3:        request.P3,
		CreatedAt: existingProde.CreatedAt,
		UpdatedAt: time.Now(),
	}

	err = s.prodeRepo.UpdateProdeSession(ctx, &prode)
	if err != nil {
		return prodes.ResponseProdeSessionDTO{}, e.NewInternalServerApiError("Error actualizando el pronóstico de sesión", err)
	}

	// // Invalidar caché relevante
	// cacheKeys := []string{
	// 	fmt.Sprintf("prode:user:%d", prode.UserID),
	// 	fmt.Sprintf("session_prodes:session:%d", prode.SessionID),
	// 	fmt.Sprintf("prode:user:%d:session:%d", prode.UserID, prode.SessionID),
	// }
	// for _, key := range cacheKeys {
	// 	s.cache.Delete(key)
	// 	fmt.Printf("Invalidated cache for key=%s\n", key)
	// }

	response := prodes.ResponseProdeSessionDTO{
		ID:        prode.ID,
		UserID:    prode.UserID,
		SessionID: prode.SessionID,
		P1:        prode.P1,
		P2:        prode.P2,
		P3:        prode.P3,
		Score:     prode.Score,
	}

	return response, nil
}

func (s *prodeService) DeleteProdeCarrera(ctx context.Context, prodeID int) e.ApiError {
	prode, err := s.prodeRepo.GetProdeCarreraByID(ctx, prodeID)
	if err != nil {
		return e.NewNotFoundApiError("El pronóstico de carrera no fue encontrado")
	}

	err = s.prodeRepo.DeleteProdeCarreraByID(ctx, prode.ID, prode.UserID)
	if err != nil {
		return e.NewInternalServerApiError("Error eliminando el pronóstico de carrera", err)
	}

	// Invalidar caché relevante
	// cacheKeys := []string{
	// 	fmt.Sprintf("prode:user:%d", prode.UserID),
	// 	fmt.Sprintf("race_prodes:session:%d", prode.SessionID),
	// 	fmt.Sprintf("prode:user:%d:session:%d", prode.UserID, prode.SessionID),
	// }
	// for _, key := range cacheKeys {
	// 	s.cache.Delete(key)
	// 	fmt.Printf("Invalidated cache for key=%s\n", key)
	// }

	return nil
}

func (s *prodeService) DeleteProdeSession(ctx context.Context, prodeID int) e.ApiError {
	prode, err := s.prodeRepo.GetProdeSessionByID(ctx, prodeID)
	if err != nil {
		return e.NewNotFoundApiError("El pronóstico de sesión no fue encontrado")
	}

	err = s.prodeRepo.DeleteProdeSessionByID(ctx, prode.ID, prode.UserID)
	if err != nil {
		return e.NewInternalServerApiError("Error eliminando el pronóstico de sesión", err)
	}

	// Invalidar caché relevante
	// cacheKeys := []string{
	// 	fmt.Sprintf("prode:user:%d", prode.UserID),
	// 	fmt.Sprintf("session_prodes:session:%d", prode.SessionID),
	// 	fmt.Sprintf("prode:user:%d:session:%d", prode.UserID, prode.SessionID),
	// }
	// for _, key := range cacheKeys {
	// 	s.cache.Delete(key)
	// 	fmt.Printf("Invalidated cache for key=%s\n", key)
	// }

	return nil
}

func (s *prodeService) DeleteProdeById(ctx context.Context, prodeID int) e.ApiError {
	sessionInfo, err := s.sessionClient.GetSessionNameAndType(prodeID)
	if err != nil {
		return e.NewInternalServerApiError("Error fetching session name and type from sessions service", err)
	}

	if isRaceSession(sessionInfo.SessionName, sessionInfo.SessionType) {
		if err := s.DeleteProdeCarrera(ctx, prodeID); err != nil {
			return err
		}
	} else {
		if err := s.DeleteProdeSession(ctx, prodeID); err != nil {
			return err
		}
	}

	return nil
}

func (s *prodeService) GetProdesByUserId(ctx context.Context, userID int) ([]prodes.ResponseProdeCarreraDTO, []prodes.ResponseProdeSessionDTO, e.ApiError) {
	// cacheKey := fmt.Sprintf("prode:user:%d", userID)
	// if cached, exists := s.cache.Get(cacheKey); exists {
	// 	if result, ok := cached.(struct {
	// 		Carrera []prodes.ResponseProdeCarreraDTO
	// 		Session []prodes.ResponseProdeSessionDTO
	// 	}); ok {
	// 		fmt.Printf("Cache hit for prodes:user:%d\n", userID)
	// 		return result.Carrera, result.Session, nil
	// 	}
	// }

	carreraProdes, sessionProdes, err := s.prodeRepo.GetProdesByUserID(ctx, userID)
	if err != nil {
		return nil, nil, e.NewInternalServerApiError("Error fetching prodes by user ID", err)
	}

	var carreraResponses []prodes.ResponseProdeCarreraDTO
	for _, prode := range carreraProdes {
		carreraResponses = append(carreraResponses, prodes.ResponseProdeCarreraDTO{
			ID:        prode.ID,
			UserID:    prode.UserID,
			SessionID: prode.SessionID,
			P1:        prode.P1,
			P2:        prode.P2,
			P3:        prode.P3,
			P4:        prode.P4,
			P5:        prode.P5,
			VSC:       prode.VSC,
			SC:        prode.SC,
			DNF:       prode.DNF,
			Score:     prode.Score,
		})
	}

	var sessionResponses []prodes.ResponseProdeSessionDTO
	for _, prode := range sessionProdes {
		sessionResponses = append(sessionResponses, prodes.ResponseProdeSessionDTO{
			ID:        prode.ID,
			UserID:    prode.UserID,
			SessionID: prode.SessionID,
			P1:        prode.P1,
			P2:        prode.P2,
			P3:        prode.P3,
			Score:     prode.Score,
		})
	}

	// Cachear el resultado
	// s.cache.Set(cacheKey, struct {
	// 	Carrera []prodes.ResponseProdeCarreraDTO
	// 	Session []prodes.ResponseProdeSessionDTO
	// }{Carrera: carreraResponses, Session: sessionResponses}, 5*time.Minute)
	// fmt.Printf("Cached prodes for user:%d\n", userID)

	return carreraResponses, sessionResponses, nil
}

func (s *prodeService) GetProdeByUserAndSession(ctx context.Context, userID, sessionID int) (*prodes.ResponseProdeCarreraDTO, *prodes.ResponseProdeSessionDTO, e.ApiError) {
	// cacheKey := fmt.Sprintf("prode:user:%d:session:%d", userID, sessionID)
	// if cached, exists := s.cache.Get(cacheKey); exists {
	// 	if result, ok := cached.(struct {
	// 		Carrera *prodes.ResponseProdeCarreraDTO
	// 		Session *prodes.ResponseProdeSessionDTO
	// 	}); ok {
	// 		fmt.Printf("Cache hit for prode:user:%d:session:%d\n", userID, sessionID)
	// 		return result.Carrera, result.Session, nil
	// 	}
	// }

	sessionInfo, err := s.sessionClient.GetSessionNameAndType(sessionID)
	if err != nil {
		fmt.Printf("Error fetching session info: %v\n", err)
		return nil, nil, e.NewInternalServerApiError("Error fetching session name and type from sessions service", err)
	}

	var carreraResponse *prodes.ResponseProdeCarreraDTO
	var sessionResponse *prodes.ResponseProdeSessionDTO

	if isRaceSession(sessionInfo.SessionName, sessionInfo.SessionType) {
		prode, err := s.prodeRepo.GetProdeCarreraByUserAndSession(ctx, userID, sessionID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Printf("No prode carrera found for userID %d and sessionID %d\n", userID, sessionID)
				return nil, nil, nil
			}
			fmt.Printf("Database error for userID %d and sessionID %d: %v\n", userID, sessionID, err)
			return nil, nil, e.NewInternalServerApiError("Error fetching prode carrera", err)
		}

		if prode != nil {
			carreraResponse = &prodes.ResponseProdeCarreraDTO{
				ID:        prode.ID,
				UserID:    prode.UserID,
				SessionID: prode.SessionID,
				P1:        prode.P1,
				P2:        prode.P2,
				P3:        prode.P3,
				P4:        prode.P4,
				P5:        prode.P5,
				VSC:       prode.VSC,
				SC:        prode.SC,
				DNF:       prode.DNF,
				Score:     prode.Score,
			}
		}
	} else {
		prode, err := s.prodeRepo.GetProdeSessionByUserAndSession(ctx, userID, sessionID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Printf("No prode session found for userID %d and sessionID %d\n", userID, sessionID)
				return nil, nil, nil
			}
			fmt.Printf("Database error for userID %d and sessionID %d: %v\n", userID, sessionID, err)
			return nil, nil, e.NewInternalServerApiError("Error fetching prode session", err)
		}

		if prode != nil {
			sessionResponse = &prodes.ResponseProdeSessionDTO{
				ID:        prode.ID,
				UserID:    prode.UserID,
				SessionID: prode.SessionID,
				P1:        prode.P1,
				P2:        prode.P2,
				P3:        prode.P3,
				Score:     prode.Score,
			}
		}
	}

	// Cachear el resultado
	// s.cache.Set(cacheKey, struct {
	// 	Carrera *prodes.ResponseProdeCarreraDTO
	// 	Session *prodes.ResponseProdeSessionDTO
	// }{Carrera: carreraResponse, Session: sessionResponse}, 5*time.Minute)
	// fmt.Printf("Cached prode for user:%d:session:%d\n", userID, sessionID)

	return carreraResponse, sessionResponse, nil
}

func (s *prodeService) GetRaceProdesBySession(ctx context.Context, sessionID int) ([]prodes.ResponseProdeCarreraDTO, e.ApiError) {
	// cacheKey := fmt.Sprintf("race_prodes:session:%d", sessionID)
	// if cached, exists := s.cache.Get(cacheKey); exists {
	// 	if raceProdes, ok := cached.([]prodes.ResponseProdeCarreraDTO); ok {
	// 		fmt.Printf("Cache hit for race_prodes:session:%d\n", sessionID)
	// 		return raceProdes, nil
	// 	}
	// }

	sessionInfo, err := s.sessionClient.GetSessionNameAndType(sessionID)
	if err != nil {
		return nil, e.NewInternalServerApiError("Error fetching session name and type from sessions service", err)
	}

	if !isRaceSession(sessionInfo.SessionName, sessionInfo.SessionType) {
		return nil, e.NewBadRequestApiError("La sesión no es una carrera válida (Race), no se pueden buscar los ProdesCarrera")
	}

	raceProdes, err := s.prodeRepo.GetRaceProdesBySession(ctx, sessionID)
	if err != nil {
		return nil, e.NewInternalServerApiError("Error fetching race prodes for the session", err)
	}

	var raceProdeResponses []prodes.ResponseProdeCarreraDTO
	for _, prode := range raceProdes {
		raceProdeResponses = append(raceProdeResponses, prodes.ResponseProdeCarreraDTO{
			ID:        prode.ID,
			UserID:    prode.UserID,
			SessionID: prode.SessionID,
			P1:        prode.P1,
			P2:        prode.P2,
			P3:        prode.P3,
			P4:        prode.P4,
			P5:        prode.P5,
			VSC:       prode.VSC,
			SC:        prode.SC,
			DNF:       prode.DNF,
			Score:     prode.Score,
		})
	}

	// // Cachear el resultado
	// ttl := 5 * time.Minute
	// if sessionInfo.DateEnd.Before(time.Now()) {
	// 	ttl = 24 * time.Hour // Sesiones finalizadas son inmutables
	// }
	// s.cache.Set(cacheKey, raceProdeResponses, ttl)
	// fmt.Printf("Cached race_prodes for session:%d\n", sessionID)

	return raceProdeResponses, nil
}

func (s *prodeService) UpdateRaceProdeForUserBySessionId(ctx context.Context, userID int, sessionID int, updatedProde prodes.UpdateProdeCarreraDTO) (prodes.ResponseProdeCarreraDTO, e.ApiError) {
	sessionDetails, err := s.sessionClient.GetSessionByID(sessionID)
	if err != nil {
		return prodes.ResponseProdeCarreraDTO{}, e.NewInternalServerApiError("Error fetching session details", err)
	}

	if !isRaceSession(sessionDetails.SessionName, sessionDetails.SessionType) {
		return prodes.ResponseProdeCarreraDTO{}, e.NewBadRequestApiError("La sesión no es de tipo 'Race'. No se puede actualizar un ProdeCarrera")
	}

	if time.Now().After(sessionDetails.DateStart) {
		return prodes.ResponseProdeCarreraDTO{}, e.NewForbiddenApiError("No se puede actualizar el pronóstico, la carrera ya ha comenzado")
	}

	prode := model.ProdeCarrera{
		ID:        updatedProde.ProdeID,
		UserID:    userID,
		SessionID: sessionID,
		P1:        updatedProde.P1,
		P2:        updatedProde.P2,
		P3:        updatedProde.P3,
		P4:        updatedProde.P4,
		P5:        updatedProde.P5,
		VSC:       updatedProde.VSC,
		SC:        updatedProde.SC,
		DNF:       updatedProde.DNF,
	}

	err = s.prodeRepo.UpdateProdeCarrera(ctx, &prode)
	if err != nil {
		return prodes.ResponseProdeCarreraDTO{}, e.NewInternalServerApiError("Error actualizando el pronóstico de carrera", err)
	}

	// Invalidar caché relevante
	// cacheKeys := []string{
	// 	fmt.Sprintf("prode:user:%d", userID),
	// 	fmt.Sprintf("race_prodes:session:%d", sessionID),
	// 	fmt.Sprintf("prode:user:%d:session:%d", userID, sessionID),
	// }
	// for _, key := range cacheKeys {
	// 	s.cache.Delete(key)
	// 	fmt.Printf("Invalidated cache for key=%s\n", key)
	// }

	response := prodes.ResponseProdeCarreraDTO{
		ID:        prode.ID,
		UserID:    prode.UserID,
		SessionID: prode.SessionID,
		P1:        prode.P1,
		P2:        prode.P2,
		P3:        prode.P3,
		P4:        prode.P4,
		P5:        prode.P5,
		VSC:       prode.VSC,
		SC:        prode.SC,
		DNF:       prode.DNF,
		Score:     prode.Score,
	}

	return response, nil
}

func (s *prodeService) GetSessionProdeBySession(ctx context.Context, sessionID int) ([]prodes.ResponseProdeSessionDTO, e.ApiError) {
	// cacheKey := fmt.Sprintf("session_prodes:session:%d", sessionID)
	// if cached, exists := s.cache.Get(cacheKey); exists {
	// 	if sessionProdes, ok := cached.([]prodes.ResponseProdeSessionDTO); ok {
	// 		fmt.Printf("Cache hit for session_prodes:session:%d\n", sessionID)
	// 		return sessionProdes, nil
	// 	}
	// }

	sessionInfo, err := s.sessionClient.GetSessionNameAndType(sessionID)
	if err != nil {
		return nil, e.NewInternalServerApiError("Error fetching session name and type from sessions service", err)
	}

	if isRaceSession(sessionInfo.SessionName, sessionInfo.SessionType) {
		return nil, e.NewBadRequestApiError("La sesión es una carrera (Race), no se pueden buscar los ProdesSession")
	}

	sessionProdes, err := s.prodeRepo.GetSessionProdesBySession(ctx, sessionID)
	if err != nil {
		return nil, e.NewInternalServerApiError("Error fetching session prodes for the session", err)
	}

	var sessionProdeResponses []prodes.ResponseProdeSessionDTO
	for _, prode := range sessionProdes {
		sessionProdeResponses = append(sessionProdeResponses, prodes.ResponseProdeSessionDTO{
			ID:        prode.ID,
			UserID:    prode.UserID,
			SessionID: prode.SessionID,
			P1:        prode.P1,
			P2:        prode.P2,
			P3:        prode.P3,
			Score:     prode.Score,
		})
	}

	// Cachear el resultado
	// ttl := 5 * time.Minute
	// if sessionInfo.DateEnd.Before(time.Now()) {
	// 	ttl = 24 * time.Hour // Sesiones finalizadas son inmutables
	// }
	// s.cache.Set(cacheKey, sessionProdeResponses, ttl)
	// fmt.Printf("Cached session_prodes for session:%d\n", sessionID)

	return sessionProdeResponses, nil
}

func (s *prodeService) GetUserProdes(ctx context.Context, userID int) ([]prodes.ResponseProdeCarreraDTO, []prodes.ResponseProdeSessionDTO, e.ApiError) {
	// cacheKey := fmt.Sprintf("prode:user:%d", userID)
	// if cached, exists := s.cache.Get(cacheKey); exists {
	// 	if result, ok := cached.(struct {
	// 		Carrera []prodes.ResponseProdeCarreraDTO
	// 		Session []prodes.ResponseProdeSessionDTO
	// 	}); ok {
	// 		fmt.Printf("Cache hit for prodes:user:%d\n", userID)
	// 		return result.Carrera, result.Session, nil
	// 	}
	// }

	userExists, err := s.userClient.GetUserByID(userID)
	if err != nil || !userExists {
		return nil, nil, e.NewNotFoundApiError("User not found")
	}

	carreraProdes, sessionProdes, err := s.prodeRepo.GetProdesByUserID(ctx, userID)
	if err != nil {
		return nil, nil, e.NewInternalServerApiError("Error fetching user prodes", err)
	}

	var carreraResponses []prodes.ResponseProdeCarreraDTO
	for _, prode := range carreraProdes {
		carreraResponses = append(carreraResponses, prodes.ResponseProdeCarreraDTO{
			ID:        prode.ID,
			UserID:    prode.UserID,
			SessionID: prode.SessionID,
			P1:        prode.P1,
			P2:        prode.P2,
			P3:        prode.P3,
			P4:        prode.P4,
			P5:        prode.P5,
			VSC:       prode.VSC,
			SC:        prode.SC,
			DNF:       prode.DNF,
			Score:     prode.Score,
		})
	}

	var sessionResponses []prodes.ResponseProdeSessionDTO
	for _, prode := range sessionProdes {
		sessionResponses = append(sessionResponses, prodes.ResponseProdeSessionDTO{
			ID:        prode.ID,
			UserID:    prode.UserID,
			SessionID: prode.SessionID,
			P1:        prode.P1,
			P2:        prode.P2,
			P3:        prode.P3,
			Score:     prode.Score,
		})
	}

	// Cachear el resultado
	// s.cache.Set(cacheKey, struct {
	// 	Carrera []prodes.ResponseProdeCarreraDTO
	// 	Session []prodes.ResponseProdeSessionDTO
	// }{Carrera: carreraResponses, Session: sessionResponses}, 5*time.Minute)
	// fmt.Printf("Cached prodes for user:%d\n", userID)

	return carreraResponses, sessionResponses, nil
}

func (s *prodeService) GetDriverDetails(ctx context.Context, driverID int) (prodes.DriverDTO, e.ApiError) {
	// cacheKey := fmt.Sprintf("driver:%d", driverID)
	// if cached, exists := s.cache.Get(cacheKey); exists {
	// 	if driver, ok := cached.(prodes.DriverDTO); ok {
	// 		fmt.Printf("Cache hit for driver:%d\n", driverID)
	// 		return driver, nil
	// 	}
	// }

	driverDetails, err := s.driverClient.GetDriverByID(driverID)
	if err != nil {
		return prodes.DriverDTO{}, e.NewInternalServerApiError("Error fetching driver details from drivers service", err)
	}

	response := prodes.DriverDTO{
		ID:          driverDetails.ID,
		FirstName:   driverDetails.FirstName,
		LastName:    driverDetails.LastName,
		FullName:    driverDetails.FullName,
		NameAcronym: driverDetails.NameAcronym,
		TeamName:    driverDetails.TeamName,
	}

	// Cachear el resultado
	// s.cache.Set(cacheKey, response, 24*time.Hour)
	// fmt.Printf("Cached driver:%d\n", driverID)

	return response, nil
}

func (s *prodeService) GetAllDrivers(ctx context.Context) ([]prodes.DriverDTO, e.ApiError) {
	// cacheKey := "all_drivers"
	// if cached, exists := s.cache.Get(cacheKey); exists {
	// 	if drivers, ok := cached.([]prodes.DriverDTO); ok {
	// 		fmt.Printf("Cache hit for all_drivers\n")
	// 		return drivers, nil
	// 	}
	// }

	drivers, err := s.driverClient.GetAllDrivers()
	if err != nil {
		return nil, e.NewInternalServerApiError("Error fetching all drivers from drivers service", err)
	}

	var driverResponses []prodes.DriverDTO
	for _, driver := range drivers {
		driverResponses = append(driverResponses, prodes.DriverDTO{
			ID:          driver.ID,
			FirstName:   driver.FirstName,
			LastName:    driver.LastName,
			FullName:    driver.FullName,
			NameAcronym: driver.NameAcronym,
			TeamName:    driver.TeamName,
		})
	}

	// Cachear el resultado
	// s.cache.Set(cacheKey, driverResponses, 24*time.Hour)
	// fmt.Printf("Cached all_drivers\n")

	return driverResponses, nil
}

func (s *prodeService) GetTopDriversBySessionId(ctx context.Context, sessionID, n int) ([]prodes.TopDriverDTO, e.ApiError) {
	topDrivers, err := s.resultsClient.GetTopDriversBySession(sessionID, n)
	if err != nil {
		return nil, e.NewInternalServerApiError("Error fetching top drivers from results service", err)
	}

	// Retornar los pilotos obtenidos
	return topDrivers, nil
}

func (s *prodeService) UpdateScoresForRaceProdes(ctx context.Context, sessionID int) e.ApiError {
	sessionDetails, err := s.sessionClient.GetSessionByID(sessionID)
	if err != nil {
		return e.NewInternalServerApiError("Error fetching session details", err)
	}

	if !isRaceSession(sessionDetails.SessionName, sessionDetails.SessionType) {
		return e.NewBadRequestApiError("La sesión no es de tipo 'Race'; no se pueden recalcular prodes carrera")
	}

	// Valores reales con defaults
	realVSC := false
	if sessionDetails.VSC != nil {
		realVSC = *sessionDetails.VSC
	}
	realSC := false
	if sessionDetails.SC != nil {
		realSC = *sessionDetails.SC
	}
	realDNF := 0
	if sessionDetails.DNF != nil {
		realDNF = *sessionDetails.DNF
	}

	realTopDrivers, err := s.resultsClient.GetTopDriversBySession(sessionID, 5)
	if err != nil {
		return e.NewInternalServerApiError("Error fetching top 5 drivers for race session", err)
	}

	raceProdes, err := s.prodeRepo.GetRaceProdesBySession(ctx, sessionID)
	if err != nil {
		return e.NewInternalServerApiError("Error fetching race prodes for session", err)
	}

	// Acumular deltas por usuario
	deltaPorUsuario := make(map[int]int)

	// Calcular nuevos scores y acumular deltas
	for _, prode := range raceProdes {
		newScore := calculateRaceScore(prode, realTopDrivers, realVSC, realSC, realDNF)
		delta := newScore - prode.Score
		if delta == 0 {
			continue
		}
		prode.Score = newScore
		deltaPorUsuario[prode.UserID] += delta
	}

	// Persistir prodes actualizados
	for _, prode := range raceProdes {
		if err := s.prodeRepo.UpdateProdeCarrera(ctx, prode); err != nil {
			return e.NewInternalServerApiError("Error updating race prode score", err)
		}
	}

	// Aplicar deltas acumulados a cada user
	for userID, delta := range deltaPorUsuario {
		if apiErr := s.prodeRepo.IncrementUserScore(ctx, userID, delta); apiErr != nil {
			return e.NewInternalServerApiError("Error updating user total score", apiErr)
		}
	}

	return nil
}

func (s *prodeService) UpdateScoresForSessionProdes(ctx context.Context, sessionID int) e.ApiError {
	realTopDrivers, err := s.resultsClient.GetTopDriversBySession(sessionID, 3)
	if err != nil {
		return e.NewInternalServerApiError("Error fetching real top drivers for session", err)
	}

	prodesSession, err := s.prodeRepo.GetSessionProdesBySession(ctx, sessionID)
	if err != nil {
		return e.NewInternalServerApiError("Error fetching prodes session for scoring", err)
	}

	deltaPorUsuario := make(map[int]int)

	for _, prode := range prodesSession {
		newScore := calculateSessionScore(prode, realTopDrivers)
		delta := newScore - prode.Score
		if delta == 0 {
			continue
		}
		prode.Score = newScore
		deltaPorUsuario[prode.UserID] += delta
	}

	for _, prode := range prodesSession {
		if err := s.prodeRepo.UpdateProdeSession(ctx, prode); err != nil {
			return e.NewInternalServerApiError("Error updating prode session score", err)
		}
	}

	for userID, delta := range deltaPorUsuario {
		if apiErr := s.prodeRepo.IncrementUserScore(ctx, userID, delta); apiErr != nil {
			return e.NewInternalServerApiError("Error updating user total score", apiErr)
		}
	}

	return nil
}

func calculateRaceScore(prode *model.ProdeCarrera, realTop []prodes.TopDriverDTO, realVSC bool, realSC bool, realDNF int) int {
	score := 0

	// P1
	if len(realTop) > 0 {
		if prode.P1 == realTop[0].DriverID {
			score += 3
		} else if driverInList(prode.P1, realTop) {
			score += 1
		}
	}

	// P2
	if len(realTop) > 1 {
		if prode.P2 == realTop[1].DriverID {
			score += 3
		} else if driverInList(prode.P2, realTop) {
			score += 1
		}
	}

	// P3
	if len(realTop) > 2 {
		if prode.P3 == realTop[2].DriverID {
			score += 3
		} else if driverInList(prode.P3, realTop) {
			score += 1
		}
	}

	// P4
	if len(realTop) > 3 {
		if prode.P4 == realTop[3].DriverID {
			score += 3
		} else if driverInList(prode.P4, realTop) {
			score += 1
		}
	}

	// P5
	if len(realTop) > 4 {
		if prode.P5 == realTop[4].DriverID {
			score += 3
		} else if driverInList(prode.P5, realTop) {
			score += 1
		}
	}

	// 2. Comparar VSC
	if prode.VSC == realVSC {
		score += 2
	}

	// 3. Comparar SC
	if prode.SC == realSC {
		score += 2
	}

	// 4. Comparar DNF
	if prode.DNF == realDNF {
		score += 5
	}

	return score
}

func calculateSessionScore(prode *model.ProdeSession, realTop []prodes.TopDriverDTO) int {
	score := 0

	// P1
	if len(realTop) > 0 {
		if prode.P1 == realTop[0].DriverID {
			score += 3
		} else if driverInList(prode.P1, realTop) {
			score += 1
		}
	}

	// P2
	if len(realTop) > 1 {
		if prode.P2 == realTop[1].DriverID {
			score += 3
		} else if driverInList(prode.P2, realTop) {
			score += 1
		}
	}

	// P3
	if len(realTop) > 2 {
		if prode.P3 == realTop[2].DriverID {
			score += 3
		} else if driverInList(prode.P3, realTop) {
			score += 1
		}
	}

	return score
}

func driverInList(driverID int, realTop []prodes.TopDriverDTO) bool {
	for _, driver := range realTop {
		if driver.DriverID == driverID {
			return true
		}
	}
	return false
}

func isRaceSession(sessionName string, sessionType string) bool {
	return sessionName == "Race" && sessionType == "Race"
}
