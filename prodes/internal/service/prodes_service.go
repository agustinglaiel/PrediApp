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
    // Hacer la llamada HTTP al microservicio de sessions para obtener el nombre y tipo de sesión
    endpoint := fmt.Sprintf("/sessions/%d/name-type", request.EventID)
    responseData, err := s.httpClient.Get(endpoint)
    if err != nil {
        // Convertir el error estándar a ApiError utilizando la función de errores personalizada
        return prodes.ResponseProdeCarreraDTO{}, e.NewInternalServerApiError("Error en la solicitud HTTP a sessions", err)
    }

    // Parsear la respuesta JSON del microservicio de sessions
    var sessionInfo prodes.SessionNameAndTypeDTO //Aca defini este dto y en ese archivo explico porqué!
    err = json.Unmarshal(responseData, &sessionInfo)
    if err != nil {
        // Convertir el error estándar a ApiError si hay un problema al parsear la respuesta
        return prodes.ResponseProdeCarreraDTO{}, e.NewInternalServerApiError("Error parseando respuesta de sessions", err)
    }

    // Verificar si la sesión es una carrera
    if sessionInfo.SessionName != "Race" {
        return prodes.ResponseProdeCarreraDTO{}, e.NewBadRequestApiError("La sesión asociada no es una carrera, no se puede crear un ProdeCarrera")
    }

    // Convertir DTO a modelo
    prode := model.ProdeCarrera{
        UserID:     request.UserID,
        EventID:    request.EventID,
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
        EventID:    prode.EventID,
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
    endpoint := fmt.Sprintf("/sessions/%d/name-type", request.EventID)
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
    if sessionInfo.SessionType == "Race" {
        return prodes.ResponseProdeSessionDTO{}, e.NewBadRequestApiError("La sesión asociada es una carrera, no se puede crear un ProdeSession")
    }

    // Convertir DTO a modelo
    prode := model.ProdeSession{
        UserID:  request.UserID,
        EventID: request.EventID,
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
        EventID: prode.EventID,
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
		EventID:    request.EventID,
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
		EventID:    prode.EventID,
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
		EventID: request.EventID,
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
		EventID: prode.EventID,
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
    // Hacer la llamada HTTP al microservicio de sessions para obtener el nombre y tipo de sesión
    endpoint := fmt.Sprintf("/sessions/%d/name-type", prodeID)
    responseData, err := s.httpClient.Get(endpoint)
    if err != nil {
        return e.NewInternalServerApiError("Error en la solicitud HTTP a sessions", err)
    }

    // Parsear la respuesta JSON del microservicio de sessions
    var sessionInfo prodes.SessionNameAndTypeDTO
    err = json.Unmarshal(responseData, &sessionInfo)
    if err != nil {
        return e.NewInternalServerApiError("Error parseando respuesta de sessions", err)
    }

    // Verificar si la sesión es de tipo "Race"
    if sessionInfo.SessionName == "Race" {
        // Es una carrera, eliminar el prode de carrera
        if err := s.DeleteProdeCarrera(ctx, prodeID); err != nil {
            return err
        }
    } else {
        // Es una sesión, eliminar el prode de sesión
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
			EventID:    prode.EventID,
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
			EventID: prode.EventID,
			P1:      prode.P1,
			P2:      prode.P2,
			P3:      prode.P3,
		})
	}

	return carreraResponses, sessionResponses, nil
}