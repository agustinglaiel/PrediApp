package service

import (
	"context"
	"errors"
	"fmt"
	client "prodes/internal/client"
	prodes "prodes/internal/dto"
	model "prodes/internal/model"
	repository "prodes/internal/repository"
	e "prodes/pkg/utils"
	"time"

	"gorm.io/gorm"
)

type prodeService struct {
	prodeRepo  repository.ProdeRepository
	httpClient *client.HttpClient
	cache      *e.Cache  // Agregar la caché
}

type ProdeServiceInterface interface {
	CreateProdeCarrera(ctx context.Context, request prodes.CreateProdeCarreraDTO) (prodes.ResponseProdeCarreraDTO, e.ApiError)
	CreateProdeSession(ctx context.Context, request prodes.CreateProdeSessionDTO) (prodes.ResponseProdeSessionDTO, e.ApiError)
	UpdateProdeCarrera(ctx context.Context, request prodes.UpdateProdeCarreraDTO) (prodes.ResponseProdeCarreraDTO, e.ApiError)
	UpdateProdeSession(ctx context.Context, request prodes.UpdateProdeSessionDTO) (prodes.ResponseProdeSessionDTO, e.ApiError)
	DeleteProdeById(ctx context.Context, prodeID int) e.ApiError 
	GetProdesByUserId(ctx context.Context, userID int) ([]prodes.ResponseProdeCarreraDTO, []prodes.ResponseProdeSessionDTO, e.ApiError)
	// GetRaceProdeByUserAndSession(ctx context.Context, userID, sessionID int) (prodes.ResponseProdeCarreraDTO, e.ApiError)
	// GetSessionProdeByUserAndSession(ctx context.Context, userID, sessionID int) (prodes.ResponseProdeSessionDTO, e.ApiError)
	GetRaceProdesBySession(ctx context.Context, sessionID int) ([]prodes.ResponseProdeCarreraDTO, e.ApiError)
	UpdateRaceProdeForUserBySessionId(ctx context.Context, userID int, sessionID int, updatedProde prodes.UpdateProdeCarreraDTO) (prodes.ResponseProdeCarreraDTO, e.ApiError)
	GetSessionProdeBySession(ctx context.Context, sessionID int) ([]prodes.ResponseProdeSessionDTO, e.ApiError)
	GetUserProdes(ctx context.Context, userID int) ([]prodes.ResponseProdeCarreraDTO, []prodes.ResponseProdeSessionDTO, e.ApiError)
	GetDriverDetails(ctx context.Context, driverID int) (prodes.DriverDTO, e.ApiError)
	GetAllDrivers(ctx context.Context) ([]prodes.DriverDTO, e.ApiError)
	GetTopDriversBySessionId(ctx context.Context, sessionID, n int) ([]prodes.TopDriverDTO, e.ApiError)
    GetProdeByUserAndSession (ctx context.Context, userID int, sessionID int)(*prodes.ResponseProdeCarreraDTO, *prodes.ResponseProdeSessionDTO, e.ApiError)
    UpdateScoresForSessionProdes(ctx context.Context, sessionID int) e.ApiError
}

// NewProdeService crea una nueva instancia de ProdeService con inyección de dependencias
func NewPrediService(prodeRepo repository.ProdeRepository, httpClient *client.HttpClient, cache *e.Cache) ProdeServiceInterface {
	return &prodeService{
		prodeRepo:  prodeRepo,
		httpClient: httpClient,
		cache:      cache,
	}
}

func (s *prodeService) CreateProdeCarrera(ctx context.Context, request prodes.CreateProdeCarreraDTO) (prodes.ResponseProdeCarreraDTO, e.ApiError) {
    existingProde, err := s.prodeRepo.GetProdeCarreraBySessionIdAndUserId(ctx, request.UserID, request.SessionID)

    if err == nil {
        // Si ya existe un ProdeCarrera, actualizarlo en lugar de crear uno nuevo
        updateRequest := prodes.UpdateProdeCarreraDTO{
            ProdeID:    existingProde.ID,
            UserID:     existingProde.UserID,
            SessionID:  existingProde.SessionID,
            P1:         request.P1,
            P2:         request.P2,
            P3:         request.P3,
            P4:         request.P4,
            P5:         request.P5,
            // FastestLap: request.FastestLap,
            VSC:        request.VSC,
            SC:         request.SC,
            DNF:        request.DNF,
        }
        return s.UpdateProdeCarrera(ctx, updateRequest)
    }

    // if err != e.NewNotFoundApiError("No prode found for this user and session") {
    //     // Si ocurrió un error diferente a "registro no encontrado", devolver el error
    //     return prodes.ResponseProdeCarreraDTO{}, e.NewInternalServerApiError("Error checking existing prode", err)
    // }

    // Hacer la llamada al cliente HTTP para obtener la información de la sesión
    sessionInfo, err := s.httpClient.GetSessionNameAndType(request.SessionID)
    if err != nil {
        return prodes.ResponseProdeCarreraDTO{}, e.NewInternalServerApiError("Error fetching session name and type from sessions service", err)
    }

    // Validar tanto el session_name como el session_type
    if !isRaceSession(sessionInfo.SessionName, sessionInfo.SessionType) {
        return prodes.ResponseProdeCarreraDTO{}, e.NewBadRequestApiError("La sesión asociada no es una carrera válida (Race), no se puede crear un ProdeCarrera")
    }

    // Convertir DTO a modelo
    prode := model.ProdeCarrera{
        UserID:     request.UserID,
        SessionID:  request.SessionID,
        P1:         request.P1,
        P2:         request.P2,
        P3:         request.P3,
        P4:         request.P4,
        P5:         request.P5,
        // FastestLap: request.FastestLap,
        VSC:        request.VSC,
        SC:         request.SC,
        DNF:        request.DNF,
        Score:      0,
    }

    // Crear el pronóstico de carrera en la base de datos
    err = s.prodeRepo.CreateProdeCarrera(ctx, &prode)
    if err != nil {
        return prodes.ResponseProdeCarreraDTO{}, e.NewInternalServerApiError("Error creando el pronóstico de carrera", err)
    }

    // Convertir el modelo a DTO de respuesta
    response := prodes.ResponseProdeCarreraDTO{
        ID:         prode.ID,
        UserID:     prode.UserID,
        SessionID:  prode.SessionID,
        P1:         prode.P1,
        P2:         prode.P2,
        P3:         prode.P3,
        P4:         prode.P4,
        P5:         prode.P5,
        // FastestLap: prode.FastestLap,
        VSC:        prode.VSC,
        SC:         prode.SC,
        DNF:        prode.DNF,
        Score:      prode.Score,
    }

    return response, nil
}

func (s *prodeService) CreateProdeSession(ctx context.Context, request prodes.CreateProdeSessionDTO) (prodes.ResponseProdeSessionDTO, e.ApiError) {
    existingProde, err := s.prodeRepo.GetProdeSessionBySessionIdAndUserId(ctx, request.UserID, request.SessionID)

    if err == nil {
        // Si ya existe un ProdeSession, actualizarlo en lugar de crear uno nuevo
        updateRequest := prodes.UpdateProdeSessionDTO{
            ProdeID:    existingProde.ID,
            UserID:     existingProde.UserID,
            SessionID:  existingProde.SessionID,
            P1:         request.P1,
            P2:         request.P2,
            P3:         request.P3,
        }
        return s.UpdateProdeSession(ctx, updateRequest)
    }

    // if err != e.NewNotFoundApiError("No prode found for this user and session") {
    //     // Si ocurrió un error diferente a "registro no encontrado", devolver el error
    //     return prodes.ResponseProdeSessionDTO{}, e.NewInternalServerApiError("Error checking existing prode", err)
    // }
    
    // Obtener la información de la sesión desde el microservicio de sesiones
    sessionInfo, err := s.httpClient.GetSessionNameAndType(request.SessionID)
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
    sessionDetails, httpErr := s.httpClient.GetSessionByID(existingProde.SessionID)
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
        ID:         existingProde.ID,
        UserID:     existingProde.UserID,  // Mantener el UserID original
        SessionID:  existingProde.SessionID,  // Mantener el SessionID original
        P1:         request.P1,
        P2:         request.P2,
        P3:         request.P3,
        P4:         request.P4,
        P5:         request.P5,
        // FastestLap: request.FastestLap,
        VSC:        request.VSC,
        SC:         request.SC,
        DNF:        request.DNF,
        CreatedAt:  existingProde.CreatedAt,
        UpdatedAt:  time.Now(),
    }

    err = s.prodeRepo.UpdateProdeCarrera(ctx, &prode)
    if err != nil {
        return prodes.ResponseProdeCarreraDTO{}, e.NewInternalServerApiError("Error actualizando el pronóstico de carrera", err)
    }

	response := prodes.ResponseProdeCarreraDTO{
        ID:         prode.ID,
        UserID:     prode.UserID,
        SessionID:  prode.SessionID,
        P1:         prode.P1,
        P2:         prode.P2,
        P3:         prode.P3,
        P4:         prode.P4,
        P5:         prode.P5,
        // FastestLap: prode.FastestLap,
        VSC:        prode.VSC,
        SC:         prode.SC,
        DNF:        prode.DNF,
        Score:      prode.Score,
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
    sessionDetails, httpErr := s.httpClient.GetSessionByID(existingProde.SessionID)
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
        UserID:    existingProde.UserID,  // Mantener el UserID original
        SessionID: existingProde.SessionID,  // Mantener el SessionID original
        P1:        request.P1,
        P2:        request.P2,
        P3:        request.P3,
        CreatedAt:  existingProde.CreatedAt,
        UpdatedAt:  time.Now(),
    }

    err = s.prodeRepo.UpdateProdeSession(ctx, &prode)
    if err != nil {
        return prodes.ResponseProdeSessionDTO{}, e.NewInternalServerApiError("Error actualizando el pronóstico de sesión", err)
    }

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
	// Buscar el prode de carrera por ID
	prode, err := s.prodeRepo.GetProdeCarreraByID(ctx, prodeID)
	if err != nil {
		return err
	}

	// Eliminar el prode de carrera
	if err := s.prodeRepo.DeleteProdeCarreraByID(ctx, prode.ID, prode.UserID); err != nil {
		return err
	}

	return nil
}

func (s *prodeService) DeleteProdeSession(ctx context.Context, prodeID int) e.ApiError {
	// Buscar el prode de sesión por ID
	prode, err := s.prodeRepo.GetProdeSessionByID(ctx, prodeID)
	if err != nil {
		return err
	}

	// Eliminar el prode de sesión
	if err := s.prodeRepo.DeleteProdeSessionByID(ctx, prode.ID, prode.UserID); err != nil {
		return err
	}

	return nil
}

func (s *prodeService) DeleteProdeById(ctx context.Context, prodeID int) e.ApiError {
    // Obtener los datos de la sesión directamente desde el microservicio de sessions
    sessionInfo, err := s.httpClient.GetSessionNameAndType(prodeID)
    if err != nil {
        return e.NewInternalServerApiError("Error fetching session name and type from sessions service", err)
    }

    // Verificar si la sesión es de tipo "Race" tanto en session_name como en session_type
    if isRaceSession(sessionInfo.SessionName, sessionInfo.SessionType) {
        // Es carrera, entonces elimina el prode en race_prode
        if err := s.DeleteProdeCarrera(ctx, prodeID); err != nil {
            return err
        }
    } else {
        // No es carrera, elimina en session_prode
        if err := s.DeleteProdeSession(ctx, prodeID); err != nil {
            return err
        }
    }

    return nil
}

func (s *prodeService) GetProdesByUserId(ctx context.Context, userID int) ([]prodes.ResponseProdeCarreraDTO, []prodes.ResponseProdeSessionDTO, e.ApiError) {
	carreraProdes, sessionProdes, err := s.prodeRepo.GetProdesByUserID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	
	var carreraResponses []prodes.ResponseProdeCarreraDTO
	for _, prode := range carreraProdes {
		carreraResponses = append(carreraResponses, prodes.ResponseProdeCarreraDTO{
			ID:         prode.ID,
			UserID:     prode.UserID,
			SessionID:  prode.SessionID,
			P1:         prode.P1,
			P2:         prode.P2,
			P3:         prode.P3,
			P4:         prode.P4,
			P5:         prode.P5,
			// FastestLap: prode.FastestLap,
			VSC:        prode.VSC,
			SC:         prode.SC,
			DNF:        prode.DNF,
            Score:      prode.Score,
		})
	}

	var sessionResponses []prodes.ResponseProdeSessionDTO
	for _, prode := range sessionProdes {
		sessionResponses = append(sessionResponses, prodes.ResponseProdeSessionDTO{
			ID:      prode.ID,
			UserID:  prode.UserID,
			SessionID: prode.SessionID,
			P1:      prode.P1,
			P2:      prode.P2,
			P3:      prode.P3,
            Score:   prode.Score,
		})
	}

	return carreraResponses, sessionResponses, nil
}

// func (s *prodeService) GetRaceProdeByUserAndSession(ctx context.Context, userID, sessionID int) (prodes.ResponseProdeCarreraDTO, e.ApiError) {
//     fmt.Printf("Fetching race prode for userID: %d, sessionID: %d\n", userID, sessionID)

//     sessionInfo, err := s.httpClient.GetSessionNameAndType(sessionID)
//     if err != nil {
//         fmt.Printf("Error fetching session info: %v\n", err)
//         return prodes.ResponseProdeCarreraDTO{}, e.NewInternalServerApiError("Error fetching session name and type from sessions service", err)
//     }

//     fmt.Printf("Session info for sessionID %d: %+v\n", sessionID, sessionInfo)

//     if !isRaceSession(sessionInfo.SessionName, sessionInfo.SessionType) {
//         fmt.Printf("Session %d is not a Race session: Name=%s, Type=%s\n", sessionID, sessionInfo.SessionName, sessionInfo.SessionType)
//         return prodes.ResponseProdeCarreraDTO{}, nil // 200 OK con null (sin pronóstico)
//     }

//     prode, err := s.prodeRepo.GetProdeCarreraByUserAndSession(ctx, userID, sessionID)
//     if err != nil {
//         if errors.Is(err, gorm.ErrRecordNotFound) {
//             fmt.Printf("No prode carrera found for userID %d and sessionID %d\n", userID, sessionID)
//             return prodes.ResponseProdeCarreraDTO{}, nil // 200 OK con null (sin pronóstico)
//         }
//         fmt.Printf("Database error for userID %d and sessionID %d: %v\n", userID, sessionID, err)
//         return prodes.ResponseProdeCarreraDTO{}, err
//     }

//     response := prodes.ResponseProdeCarreraDTO{
//         ID:         prode.ID,
//         UserID:     prode.UserID,
//         SessionID:  prode.SessionID,
//         P1:         prode.P1,
//         P2:         prode.P2,
//         P3:         prode.P3,
//         P4:         prode.P4,
//         P5:         prode.P5,
//         FastestLap: prode.FastestLap,
//         VSC:        prode.VSC,
//         SC:         prode.SC,
//         DNF:        prode.DNF,
//     }

//     fmt.Printf("Found prode carrera for userID %d and sessionID %d: %+v\n", userID, sessionID, response)
//     return response, nil
// }

// func (s *prodeService) GetSessionProdeByUserAndSession(ctx context.Context, userID, sessionID int) (prodes.ResponseProdeSessionDTO, e.ApiError) {
//     fmt.Printf("Fetching session prode for userID: %d, sessionID: %d\n", userID, sessionID)

//     sessionInfo, err := s.httpClient.GetSessionNameAndType(sessionID)
//     if err != nil {
//         fmt.Printf("Error fetching session info: %v\n", err)
//         return prodes.ResponseProdeSessionDTO{}, e.NewInternalServerApiError("Error fetching session name and type from sessions service", err)
//     }

//     fmt.Printf("Session info for sessionID %d: %+v\n", sessionID, sessionInfo)

//     if isRaceSession(sessionInfo.SessionName, sessionInfo.SessionType) {
//         fmt.Printf("Session %d is a Race session: Name=%s, Type=%s\n", sessionID, sessionInfo.SessionName, sessionInfo.SessionType)
//         return prodes.ResponseProdeSessionDTO{}, nil // 200 OK con null (sin pronóstico)
//     }

//     prode, err := s.prodeRepo.GetProdeSessionByUserAndSession(ctx, userID, sessionID)
//     if err != nil {
//         if errors.Is(err, gorm.ErrRecordNotFound) {
//             fmt.Printf("No prode session found for userID %d and sessionID %d\n", userID, sessionID)
//             return prodes.ResponseProdeSessionDTO{}, nil // 200 OK con null (sin pronóstico)
//         }
//         fmt.Printf("Database error for userID %d and sessionID %d: %v\n", userID, sessionID, err)
//         return prodes.ResponseProdeSessionDTO{}, err
//     }

//     response := prodes.ResponseProdeSessionDTO{
//         ID:        prode.ID,
//         UserID:    prode.UserID,
//         SessionID: prode.SessionID,
//         P1:        prode.P1,
//         P2:        prode.P2,
//         P3:        prode.P3,
//     }

//     fmt.Printf("Found prode session for userID %d and sessionID %d: %+v\n", userID, sessionID, response)
//     return response, nil
// }

func (s *prodeService) GetProdeByUserAndSession(ctx context.Context, userID, sessionID int) (*prodes.ResponseProdeCarreraDTO, *prodes.ResponseProdeSessionDTO, e.ApiError) {
    fmt.Printf("Fetching prode for userID: %d, sessionID: %d\n", userID, sessionID)

    // Obtener los datos de la sesión directamente desde el microservicio de sessions
    sessionInfo, err := s.httpClient.GetSessionNameAndType(sessionID)
    if err != nil {
        fmt.Printf("Error fetching session info: %v\n", err)
        return nil, nil, e.NewInternalServerApiError("Error fetching session name and type from sessions service", err)
    }

    fmt.Printf("Session info for sessionID %d: %+v\n", sessionID, sessionInfo)

    // Determinar si es una sesión de carrera (Race - Race)
    if isRaceSession(sessionInfo.SessionName, sessionInfo.SessionType) {
        // Buscar en prode_carreras
        prode, err := s.prodeRepo.GetProdeCarreraByUserAndSession(ctx, userID, sessionID)
        if err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                fmt.Printf("No prode carrera found for userID %d and sessionID %d\n", userID, sessionID)
                return nil, nil, nil // Devolver nil para ambos DTOs (200 OK implícito con [] en el controlador)
            }
            fmt.Printf("Database error for userID %d and sessionID %d: %v\n", userID, sessionID, err)
            return nil, nil, err // Devolver error si no es 404
        }

        // Verificar si prode es nil antes de acceder a sus campos
        if prode == nil {
            fmt.Printf("Prode carrera is nil for userID %d and sessionID %d\n", userID, sessionID)
            return nil, nil, nil // Devolver nil para ambos DTOs (200 OK implícito con [] en el controlador)
        }

        // Convertir prode_carrera a ResponseProdeCarreraDTO
        response := &prodes.ResponseProdeCarreraDTO{
            ID:         prode.ID,
            UserID:     prode.UserID,
            SessionID:  prode.SessionID,
            P1:         prode.P1,
            P2:         prode.P2,
            P3:         prode.P3,
            P4:         prode.P4,
            P5:         prode.P5,
            // FastestLap: prode.FastestLap,
            VSC:        prode.VSC,
            SC:         prode.SC,
            DNF:        prode.DNF,
            Score:      prode.Score,
        }
        fmt.Printf("Found prode carrera for userID %d and sessionID %d: %+v\n", userID, sessionID, response)
        return response, nil, nil // Devolver solo el DTO de carrera, nil para sesión, y nil para error
    }

    // Buscar en prode_sessions (cualquier otra sesión)
    prode, err := s.prodeRepo.GetProdeSessionByUserAndSession(ctx, userID, sessionID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            fmt.Printf("No prode session found for userID %d and sessionID %d\n", userID, sessionID)
            return nil, nil, nil // Devolver nil para ambos DTOs (200 OK implícito con [] en el controlador)
        }
        fmt.Printf("Database error for userID %d and sessionID %d: %v\n", userID, sessionID, err)
        return nil, nil, err // Devolver error si no es 404
    }

    // Verificar si prode es nil antes de acceder a sus campos
    if prode == nil {
        fmt.Printf("Prode session is nil for userID %d and sessionID %d\n", userID, sessionID)
        return nil, nil, nil // Devolver nil para ambos DTOs (200 OK implícito con [] en el controlador)
    }

    // Convertir prode_session a ResponseProdeSessionDTO
    response := &prodes.ResponseProdeSessionDTO{
        ID:        prode.ID,
        UserID:    prode.UserID,
        SessionID: prode.SessionID,
        P1:        prode.P1,
        P2:        prode.P2,
        P3:        prode.P3,
        Score:     prode.Score,
    }
    fmt.Printf("Found prode session for userID %d and sessionID %d: %+v\n", userID, sessionID, response)
    return nil, response, nil // Devolver nil para carrera, solo el DTO de sesión, y nil para error
}


func (s *prodeService) GetRaceProdesBySession(ctx context.Context, sessionID int) ([]prodes.ResponseProdeCarreraDTO, e.ApiError) {
    // Obtener los datos de la sesión directamente desde el microservicio de sessions
    sessionInfo, err := s.httpClient.GetSessionNameAndType(sessionID)
    if err != nil {
        return nil, e.NewInternalServerApiError("Error fetching session name and type from sessions service", err)
    }

    // Verificar que la sesión sea de tipo "Race"
    if !isRaceSession(sessionInfo.SessionName, sessionInfo.SessionType) {
        return nil, e.NewBadRequestApiError("La sesión no es una carrera válida (Race), no se pueden buscar los ProdesCarrera")
    }

    // Obtener todos los pronósticos de carrera para la sesión específica
    raceProdes, err := s.prodeRepo.GetRaceProdesBySession(ctx, sessionID)
    if err != nil {
        return nil, e.NewInternalServerApiError("Error fetching race prodes for the session", err)
    }

    // Convertir los modelos a DTOs de respuesta
    var raceProdeResponses []prodes.ResponseProdeCarreraDTO
    for _, prode := range raceProdes {
        raceProdeResponses = append(raceProdeResponses, prodes.ResponseProdeCarreraDTO{
            ID:         prode.ID,
            UserID:     prode.UserID,
            SessionID:  prode.SessionID,
            P1:         prode.P1,
            P2:         prode.P2,
            P3:         prode.P3,
            P4:         prode.P4,
            P5:         prode.P5,
            // FastestLap: prode.FastestLap,
            VSC:        prode.VSC,
            SC:         prode.SC,
            DNF:        prode.DNF,
            Score:      prode.Score,
        })
    }

    // Retornar la lista de pronósticos de carrera
    return raceProdeResponses, nil
}

func (s *prodeService) UpdateRaceProdeForUserBySessionId(ctx context.Context, userID int, sessionID int, updatedProde prodes.UpdateProdeCarreraDTO) (prodes.ResponseProdeCarreraDTO, e.ApiError) {
    // Obtener detalles de la sesión directamente desde el microservicio de sessions
    sessionDetails, err := s.httpClient.GetSessionByID(sessionID)
    if err != nil {
        return prodes.ResponseProdeCarreraDTO{}, e.NewInternalServerApiError("Error fetching session details", err)
    }

    // Validar si la sesión es una carrera (Race)
    if !isRaceSession(sessionDetails.SessionName, sessionDetails.SessionType) {
        return prodes.ResponseProdeCarreraDTO{}, e.NewBadRequestApiError("La sesión no es de tipo 'Race'. No se puede actualizar un ProdeCarrera.")
    }

    // Validar si la carrera ya ha comenzado comparando la fecha de inicio
    if time.Now().After(sessionDetails.DateStart) {
        return prodes.ResponseProdeCarreraDTO{}, e.NewForbiddenApiError("No se puede actualizar el pronóstico, la carrera ya ha comenzado.")
    }

    // Convertir el DTO de actualización en un modelo de ProdeCarrera
    prode := model.ProdeCarrera{
        ID:         updatedProde.ProdeID,
        UserID:     userID,
        SessionID:  sessionID,
        P1:         updatedProde.P1,
        P2:         updatedProde.P2,
        P3:         updatedProde.P3,
        P4:         updatedProde.P4,
        P5:         updatedProde.P5,
        // FastestLap: updatedProde.FastestLap,
        VSC:        updatedProde.VSC,
        SC:         updatedProde.SC,
        DNF:        updatedProde.DNF,
    }

    // Llamar al repositorio para actualizar el ProdeCarrera
    err = s.prodeRepo.UpdateProdeCarrera(ctx, &prode)
    if err != nil {
        return prodes.ResponseProdeCarreraDTO{}, e.NewInternalServerApiError("Error actualizando el pronóstico de carrera", err)
    }

    // Convertir el modelo actualizado en un DTO de respuesta
    response := prodes.ResponseProdeCarreraDTO{
        ID:         prode.ID,
        UserID:     prode.UserID,
        SessionID:  prode.SessionID,
        P1:         prode.P1,
        P2:         prode.P2,
        P3:         prode.P3,
        P4:         prode.P4,
        P5:         prode.P5,
        // FastestLap: prode.FastestLap,
        VSC:        prode.VSC,
        SC:         prode.SC,
        DNF:        prode.DNF,
    }

    return response, nil
}

func (s *prodeService) GetSessionProdeBySession(ctx context.Context, sessionID int) ([]prodes.ResponseProdeSessionDTO, e.ApiError) {
    // Obtener detalles de la sesión directamente desde el microservicio de sessions
    sessionInfo, err := s.httpClient.GetSessionNameAndType(sessionID)
    if err != nil {
        return nil, e.NewInternalServerApiError("Error fetching session name and type from sessions service", err)
    }

    // Verificar que la sesión NO sea de tipo "Race"
    if isRaceSession(sessionInfo.SessionName, sessionInfo.SessionType) {
        return nil, e.NewBadRequestApiError("La sesión es una carrera (Race), no se pueden buscar los ProdesSession")
    }

    // Obtener todos los pronósticos de sesión para la sesión específica
    sessionProdes, err := s.prodeRepo.GetSessionProdesBySession(ctx, sessionID)
    if err != nil {
        return nil, e.NewInternalServerApiError("Error fetching session prodes for the session", err)
    }

    // Convertir los modelos a DTOs de respuesta
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

    // Retornar la lista de pronósticos de sesión
    return sessionProdeResponses, nil
}

func (s *prodeService) GetUserProdes(ctx context.Context, userID int) ([]prodes.ResponseProdeCarreraDTO, []prodes.ResponseProdeSessionDTO, e.ApiError) {
    // Llamar al cliente HTTP para verificar si el usuario existe en el microservicio de users
    userExists, err := s.httpClient.GetUserByID(userID)
    if err != nil || !userExists {
        return nil, nil, e.NewNotFoundApiError("User not found")
    }

    // Obtener todos los prodes (carrera y sesión) para el usuario
    carreraProdes, sessionProdes, err := s.prodeRepo.GetProdesByUserID(ctx, userID)
    if err != nil {
        return nil, nil, e.NewInternalServerApiError("Error fetching user prodes", err)
    }

    // Convertir los prodes de carrera a DTOs de respuesta
    var carreraResponses []prodes.ResponseProdeCarreraDTO
    for _, prode := range carreraProdes {
        carreraResponses = append(carreraResponses, prodes.ResponseProdeCarreraDTO{
            ID:         prode.ID,
            UserID:     prode.UserID,
            SessionID:  prode.SessionID,
            P1:         prode.P1,
            P2:         prode.P2,
            P3:         prode.P3,
            P4:         prode.P4,
            P5:         prode.P5,
            // FastestLap: prode.FastestLap,
            VSC:        prode.VSC,
            SC:         prode.SC,
            DNF:        prode.DNF,
            Score:      prode.Score,
        })
    }

    // Convertir los prodes de sesión a DTOs de respuesta
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

    return carreraResponses, sessionResponses, nil
}

func (s *prodeService) GetDriverDetails(ctx context.Context, driverID int) (prodes.DriverDTO, e.ApiError) {
    // Llamar al microservicio de drivers para obtener los detalles del piloto
    driverDetails, err := s.httpClient.GetDriverByID(driverID)
    if err != nil {
        return prodes.DriverDTO{}, e.NewInternalServerApiError("Error fetching driver details from drivers service", err)
    }

    // Convertir los detalles del piloto a DTO de respuesta
    response := prodes.DriverDTO{
        ID:          driverDetails.ID,
        FirstName:   driverDetails.FirstName,
        LastName:    driverDetails.LastName,
        FullName:    driverDetails.FullName,
        NameAcronym: driverDetails.NameAcronym,
        TeamName:    driverDetails.TeamName,
    }

    return response, nil
}

func (s *prodeService) GetAllDrivers(ctx context.Context) ([]prodes.DriverDTO, e.ApiError) {
    // Llamar al microservicio de drivers para obtener todos los pilotos
    drivers, err := s.httpClient.GetAllDrivers()
    if err != nil {
        return nil, e.NewInternalServerApiError("Error fetching all drivers from drivers service", err)
    }

    // Convertir los detalles de los pilotos a DTOs de respuesta
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

    return driverResponses, nil
}

func (s *prodeService) GetTopDriversBySessionId(ctx context.Context, sessionID, n int) ([]prodes.TopDriverDTO, e.ApiError) {
    topDrivers, err := s.httpClient.GetTopDriversBySession(sessionID, n)
    if err != nil {
        return nil, e.NewInternalServerApiError("Error fetching top drivers from results service", err)
    }

    // Retornar los pilotos obtenidos
    return topDrivers, nil
}

func (s *prodeService) UpdateScoresForSessionProdes(ctx context.Context, sessionID int) e.ApiError {
    // 1. Obtener el top 3 real de la sesión
    fmt.Println("Fetching top drivers for session:", sessionID) // Debugging
    realTopDrivers, err := s.httpClient.GetTopDriversBySession(sessionID, 3)
    if err != nil {
        return e.NewInternalServerApiError("Error fetching real top drivers for session", err)
    }
    
    fmt.Println("Real top drivers:", realTopDrivers) // Debugging

    // 2. Obtener todos los ProdeSession de la BD para esta sesión
    prodesSession, apiErr := s.prodeRepo.GetSessionProdesBySession(ctx, sessionID)
    if apiErr != nil {
        return e.NewInternalServerApiError("Error fetching prodes session for scoring", apiErr)
    }

    // 3. Recorrer cada prode y calcular
    for _, prode := range prodesSession {
        // Calcular el score
        score := calculateSessionScore(prode, realTopDrivers)
        prode.Score = score

        // Actualizar en la BD
        if err := s.prodeRepo.UpdateProdeSession(ctx, prode); err != nil {
            return e.NewInternalServerApiError("Error updating prode session score", err)
        }
    }

    return nil
}

// calculateSessionScore calcula el puntaje basándose en P1,P2,P3 del prode y en un array de 3 drivers reales
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


// driverInList revisa si un driverID está en la lista de drivers del top real
func driverInList(driverID int, realTop []prodes.TopDriverDTO) bool {
    for _, driver := range realTop {
        if driver.DriverID == driverID {
            return true
        }
    }
    return false
}



//Función auxiliar para mayor modularidad y me devuelve el bool de si es session name y type = race. 
func isRaceSession(sessionName string, sessionType string) bool {
    return sessionName == "Race" && sessionType == "Race"
}

