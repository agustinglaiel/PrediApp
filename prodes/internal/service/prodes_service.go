package service

import (
	"context"
	"encoding/json"
	"fmt"
	client "prodes/internal/client"
	prodes "prodes/internal/dto"
	model "prodes/internal/model"
	repository "prodes/internal/repository"
	e "prodes/pkg/utils"
)

type prodeService struct {
	prodeRepo  repository.ProdeRepository
	httpClient *client.HttpClient // Agregar el cliente HTTP
}

type ProdeServiceInterface interface {
	CreateProdeCarrera(ctx context.Context, request prodes.CreateProdeCarreraDTO) (prodes.ResponseProdeCarreraDTO, e.ApiError)
	CreateProdeSession(ctx context.Context, request prodes.CreateProdeSessionDTO) (prodes.ResponseProdeSessionDTO, e.ApiError)
	UpdateProdeCarrera(ctx context.Context, request prodes.UpdateProdeCarreraDTO) (prodes.ResponseProdeCarreraDTO, e.ApiError)
	UpdateProdeSession(ctx context.Context, request prodes.UpdateProdeSessionDTO) (prodes.ResponseProdeSessionDTO, e.ApiError)
	DeleteProdeById(ctx context.Context, prodeID int) e.ApiError 
	GetProdesByUserId(ctx context.Context, userID int) ([]prodes.ResponseProdeCarreraDTO, []prodes.ResponseProdeSessionDTO, e.ApiError)
}

// NewProdeService crea una nueva instancia de ProdeService con inyección de dependencias
func NewPrediService(prodeRepo repository.ProdeRepository, httpClient *client.HttpClient) ProdeServiceInterface {
	return &prodeService{
		prodeRepo:  prodeRepo,
		httpClient: httpClient, // Inyectar el cliente HTTP
	}
}

func (s *prodeService) CreateProdeCarrera(ctx context.Context, request prodes.CreateProdeCarreraDTO) (prodes.ResponseProdeCarreraDTO, e.ApiError) {
    // Llamar al cliente HTTP para obtener el nombre y tipo de sesión
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
        FastestLap: request.FastestLap,
        VSC:        request.VSC,
        SC:         request.SC,
        DNF:        request.DNF,
    }

    // Crear el pronóstico de carrera en la base de datos
    err = s.prodeRepo.CreateProdeCarrera(ctx, &prode)
    if err != nil {
        // Convertir el error estándar a ApiError al crear el prode
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
        FastestLap: prode.FastestLap,
        VSC:        prode.VSC,
        SC:         prode.SC,
        DNF:        prode.DNF,
    }

    return response, nil
}

func (s *prodeService) CreateProdeSession(ctx context.Context, request prodes.CreateProdeSessionDTO) (prodes.ResponseProdeSessionDTO, e.ApiError) {
    // Hacer la llamada HTTP al microservicio de sessions para obtener el nombre y tipo de sesión
    endpoint := fmt.Sprintf("/sessions/%d/name-type", request.SessionID)
    responseData, err := s.httpClient.Get(endpoint)
    if err != nil {
        return prodes.ResponseProdeSessionDTO{}, e.NewInternalServerApiError("Error en la solicitud HTTP a sessions", err)
    }

    // Parsear la respuesta JSON del microservicio de sessions
    var sessionInfo prodes.SessionNameAndTypeDTO
    err = json.Unmarshal(responseData, &sessionInfo)
    if err != nil {
        return prodes.ResponseProdeSessionDTO{}, e.NewInternalServerApiError("Error parseando respuesta de sessions", err)
    }

    // Verificar si la sesión es de tipo "Race"
    if isRaceSession(sessionInfo.SessionName, sessionInfo.SessionType) {
		return prodes.ResponseProdeSessionDTO{}, e.NewBadRequestApiError("La sesión asociada no es una carrera válida (Race), no se puede crear un ProdeCarrera")
	}

    // Convertir DTO a modelo
    prode := model.ProdeSession{
        UserID:  request.UserID,
        SessionID: request.SessionID,
        P1:      request.P1,
        P2:      request.P2,
        P3:      request.P3,
    }

    // Crear el pronóstico de sesión en la base de datos
    err = s.prodeRepo.CreateProdeSession(ctx, &prode)
    if err != nil {
        return prodes.ResponseProdeSessionDTO{}, e.NewInternalServerApiError("Error creando el pronóstico de sesión", err)
    }

    // Convertir el modelo a DTO de respuesta
    response := prodes.ResponseProdeSessionDTO{
        ID:      prode.ID,
        UserID:  prode.UserID,
        SessionID: prode.SessionID,
        P1:      prode.P1,
        P2:      prode.P2,
        P3:      prode.P3,
    }

    return response, nil
}

func (s *prodeService) UpdateProdeCarrera(ctx context.Context, request prodes.UpdateProdeCarreraDTO) (prodes.ResponseProdeCarreraDTO, e.ApiError) {
	prode := model.ProdeCarrera{
		ID:         request.ProdeID,
		UserID:     request.UserID,
		SessionID:  request.SessionID,
		P1:         request.P1,
		P2:         request.P2,
		P3:         request.P3,
		P4:         request.P4,
		P5:         request.P5,
		FastestLap: request.FastestLap,
		VSC:        request.VSC,
		SC:         request.SC,
		DNF:        request.DNF,
	}
	err := s.prodeRepo.UpdateProdeCarrera(ctx, &prode)
	if err != nil {
		return prodes.ResponseProdeCarreraDTO{}, err
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
		FastestLap: prode.FastestLap,
		VSC:        prode.VSC,
		SC:         prode.SC,
		DNF:        prode.DNF,
	}
	return response, nil
}

func (s *prodeService) UpdateProdeSession(ctx context.Context, request prodes.UpdateProdeSessionDTO) (prodes.ResponseProdeSessionDTO, e.ApiError) {
	prode := model.ProdeSession{
		ID:      request.ProdeID,
		UserID:  request.UserID,
		SessionID: request.SessionID,
		P1:      request.P1,
		P2:      request.P2,
		P3:      request.P3,
	}
	err := s.prodeRepo.UpdateProdeSession(ctx, &prode)
	if err != nil {
		return prodes.ResponseProdeSessionDTO{}, err
	}
	response := prodes.ResponseProdeSessionDTO{
		ID:      prode.ID,
		UserID:  prode.UserID,
		SessionID: prode.SessionID,
		P1:      prode.P1,
		P2:      prode.P2,
		P3:      prode.P3,
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
    // Usar el cliente HTTP para obtener el nombre y tipo de sesión
    sessionInfo, err := s.httpClient.GetSessionNameAndType(prodeID)
    if err != nil {
        return e.NewInternalServerApiError("Error fetching session name and type from sessions service", err)
    }

    // Verificar si la sesión es de tipo "Race" tanto en session_name como en session_type
    if isRaceSession(sessionInfo.SessionName, sessionInfo.SessionType) {
		//Es carrera, entonces elimina el prode en race_prode
		if err := s.DeleteProdeCarrera(ctx, prodeID); err != nil {
			return err
		}
	} else {
		//No es carrera, elimina en session_prode
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
			FastestLap: prode.FastestLap,
			VSC:        prode.VSC,
			SC:         prode.SC,
			DNF:        prode.DNF,
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
		})
	}

	return carreraResponses, sessionResponses, nil
}

//Función auxiliar para mayor modularidad y me devuelve el bool de si es session name y type = race. 
func isRaceSession(sessionName string, sessionType string) bool {
    return sessionName == "Race" && sessionType == "Race"
}