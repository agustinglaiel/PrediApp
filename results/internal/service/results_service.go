package service

import (
	"context"
	"errors"
	"fmt"
	"results/internal/client"
	"results/internal/dto"
	"results/internal/model"
	"results/internal/repository"
	e "results/pkg/utils"

	"gorm.io/gorm"
)

type resultService struct {
	resultRepo repository.ResultRepository
	client     *client.HttpClient
}

type ResultService interface {
	FetchResultsFromExternalAPI(ctx context.Context, sessionId int) ([]dto.ResponseResultDTO, e.ApiError)
	UpdateResult(ctx context.Context, resultID int, request dto.UpdateResultDTO) (dto.ResponseResultDTO, e.ApiError)
	GetResultsOrderedByPosition(ctx context.Context, sessionID int) ([]dto.ResponseResultDTO, e.ApiError)
	GetFastestLapInSession(ctx context.Context, sessionID int) (dto.ResponseResultDTO, e.ApiError)
	CreateResult(ctx context.Context, request dto.CreateResultDTO) (dto.ResponseResultDTO, e.ApiError)
	// GetResultByID(ctx context.Context, resultID int) (dto.ResponseResultDTO, e.ApiError)
	DeleteResult(ctx context.Context, resultID int) e.ApiError
	GetAllResults(ctx context.Context) ([]dto.ResponseResultDTO, e.ApiError)
	// GetResultsForDriverAcrossSessions(ctx context.Context, driverID int) ([]dto.ResponseResultDTO, e.ApiError)
	// GetBestPositionForDriver(ctx context.Context, driverID int) (dto.ResponseResultDTO, e.ApiError)
	GetTopNDriversInSession(ctx context.Context, sessionID int, n int) ([]dto.TopDriverDTO, e.ApiError)
	DeleteAllResultsForSession(ctx context.Context, sessionID int) e.ApiError
	// GetResultsForSessionByDriverName(ctx context.Context, sessionID int, driverName string) ([]dto.ResponseResultDTO, e.ApiError)
	// GetTotalFastestLapsForDriver(ctx context.Context, driverID int) (int, e.ApiError)
	// GetLastResultForDriver(ctx context.Context, driverID int) (dto.ResponseResultDTO, e.ApiError)
    CreateSessionResultsAdmin(ctx context.Context, bulkRequest dto.CreateBulkResultsDTO) ([]dto.ResponseResultDTO, e.ApiError)
}

func NewResultService(resultRepo repository.ResultRepository, client *client.HttpClient) ResultService {
	return &resultService{
		resultRepo: resultRepo,
		client:     client,
	}
}

// FetchResultsFromExternalAPI obtiene los resultados de una API externa y los inserta o actualiza en la base de datos
func (s *resultService) FetchResultsFromExternalAPI(ctx context.Context, sessionID int) ([]dto.ResponseResultDTO, e.ApiError) {
    // 1. Obtener sessionKey llamando al otro microservicio
    sessionKey, err := s.client.GetSessionKeyBySessionID(sessionID)
    if err != nil {
        return nil, e.NewInternalServerApiError("Error obteniendo session key", err)
    }

    // 2. Obtener las "positions" desde la API externa
    positions, err := s.client.GetPositions(sessionKey)
    if err != nil {
        return nil, e.NewInternalServerApiError("Error fetching positions from external API", err)
    }

    // Mapa para quedarnos con la última posición reportada (driverNumber -> *int)
    finalPositions := make(map[int]*int)
    for _, pos := range positions {
        // pos.Position es *int
        finalPositions[pos.DriverNumber] = pos.Position
    }

    var responseResults []dto.ResponseResultDTO

    // 3. Para cada driverNumber, determinamos la vuelta más rápida y actualizamos/insertamos en DB
    for driverNumber, positionPtr := range finalPositions {
        laps, err := s.client.GetLaps(sessionKey, driverNumber)
        if err != nil {
            // Simplemente mostramos en consola y continuamos
            fmt.Printf("Error obteniendo vueltas para driver %d: %v\n", driverNumber, err)
            continue
        }

        // 4. Encontrar la vuelta más rápida o asignar DNS si no hay vueltas
        var fastestLap float64
        if len(laps) == 0 {
            // Si no hay vueltas válidas, asignamos DNS
            fastestLap = 0 // O podrías usar nil si prefieres, pero 0 es más explícito
        } else {
            for _, lap := range laps {
                if fastestLap == 0 || lap.LapDuration < fastestLap {
                    fastestLap = lap.LapDuration
                }
            }
        }

        // 5. Obtener info completa del driver desde microservicio de drivers
        driverInfo, err := s.client.GetDriverByNumber(driverNumber)
        if err != nil {
            fmt.Printf("Error obteniendo info del driver_number %d: %v\n", driverNumber, err)
            continue
        }

        // 6. Ver si ya existe un resultado para (driver, session)
        existingResult, _ := s.resultRepo.GetResultByDriverAndSession(ctx, driverInfo.ID, sessionID)

        // Definir status según si hay o no position o vueltas
        status := "DNS" // Default a DNS si no hay vueltas
        if positionPtr != nil {
            status = "FINISHED"
        } else if len(laps) > 0 {
            status = "DNF" // Solo DNF si hay vueltas pero no posición
        }

        if existingResult == nil {
            // Crear nuevo result
            newResult := &model.Result{
                SessionID:      sessionID,
                DriverID:       driverInfo.ID,
                Position:       positionPtr,   // *int
                Status:         status,        // "FINISHED", "DNF" o "DNS"
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

        // 7. Obtener la info de la sesión para armar el DTO de respuesta
        sessionData, err := s.client.GetSessionByID(sessionID)
        if err != nil {
            return nil, e.NewInternalServerApiError("Error fetching session data", err)
        }

        // 8. Construir el DTO
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

//ESTO SOLO SIRVE PARA CREAR UN RESULTADO A LA VEZ
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

    // 7. Construir DTO de respuesta
    response := dto.ResponseResultDTO{
        ID:             newResult.ID,
        Position:       newResult.Position,
        Status:         newResult.Status,
        FastestLapTime: newResult.FastestLapTime,
        // Driver y Session lo puedes rellenar si has hecho Preload en repo, o usando otros microservicios
        Driver: dto.ResponseDriverDTO{ ID: newResult.DriverID },
        Session: dto.ResponseSessionDTO{ ID: newResult.SessionID },
        CreatedAt: newResult.CreatedAt,
        UpdatedAt: newResult.UpdatedAt,
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
    } else {
        // Si Position es nil en request, no necesariamente la borras:
        //  - Podrías dejarla tal cual, o
        //  - Si explicitamente deseas "quitársela", significaría un status != FINISHED
        //    Depende de tu negocio. Un approach:
        // if request.Status == "DNF" { 
        //   result.Position = nil 
        // }
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

    // 6. Construir respuesta
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

	return responseResults, nil
}

// GetFastestLapInSession obtiene el piloto con la vuelta más rápida en una sesión específica
func (s *resultService) GetFastestLapInSession(ctx context.Context, sessionID int) (dto.ResponseResultDTO, e.ApiError) {
    // Verificar si existe el sessionID en la tabla de resultados
    exists, err := s.resultRepo.ExistsSessionInResults(ctx, sessionID)
    if err != nil {
        return dto.ResponseResultDTO{}, e.NewInternalServerApiError("Error verifying session existence in results", err)
    }
    if !exists {
        return dto.ResponseResultDTO{}, e.NewNotFoundApiError("No results found for the given session ID")
    }

    // Obtener la vuelta más rápida de la sesión
    results, err := s.resultRepo.GetResultsBySessionID(ctx, sessionID)  // Obtenemos todos los resultados
    if err != nil {
        return dto.ResponseResultDTO{}, e.NewInternalServerApiError("Error fetching session results", err)
    }

    var fastestResult *model.Result
    for _, result := range results {
        // Ignorar tiempos nulos o 0
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

    return response, nil
}

// GetResultByID obtiene un resultado específico por su ID
// func (s *resultService) GetResultByID(ctx context.Context, resultID int) (dto.ResponseResultDTO, e.ApiError) {
// 	// Validar que el resultID no sea 0
// 	if resultID == 0 {
// 		return dto.ResponseResultDTO{}, e.NewBadRequestApiError("El ID del resultado no puede ser 0")
// 	}

// 	// Intentar obtener el resultado por su ID
// 	result, err := s.resultRepo.GetResultByID(ctx, resultID)
// 	if err != nil {
// 		// Manejar el caso donde no se encuentra el resultado
// 		if err == e.NewNotFoundApiError("Result not found") {
// 			return dto.ResponseResultDTO{}, e.NewNotFoundApiError("No se encontró el resultado con el ID proporcionado")
// 		}
// 		// Manejar cualquier otro error
// 		return dto.ResponseResultDTO{}, e.NewInternalServerApiError("Error al obtener el resultado por su ID", err)
// 	}

// 	// Convertir el modelo en DTO
// 	response := dto.ResponseResultDTO{
// 		ID:             result.ID,
// 		Position:       result.Position,
// 		FastestLapTime: result.FastestLapTime,
// 		Driver: dto.ResponseDriverDTO{
// 			ID:          result.Driver.ID,
// 			FirstName:   result.Driver.FirstName,
// 			LastName:    result.Driver.LastName,
// 			FullName:    result.Driver.FullName,
// 			NameAcronym: result.Driver.NameAcronym,
// 			TeamName:    result.Driver.TeamName,
// 		},
// 		Session: dto.ResponseSessionDTO{
// 			ID:               result.Session.ID,
// 			CircuitShortName: result.Session.CircuitShortName,
// 			CountryName:      result.Session.CountryName,
// 			Location:         result.Session.Location,
// 			SessionName:      result.Session.SessionName,
// 			SessionType:      result.Session.SessionType,
// 			DateStart:        result.Session.DateStart,
// 		},
// 		CreatedAt: result.CreatedAt,
// 		UpdatedAt: result.UpdatedAt,
// 	}

// 	return response, nil
// }

// DeleteResult elimina un resultado específico
func (s *resultService) DeleteResult(ctx context.Context, resultID int) e.ApiError {
    // Validar que el resultID no sea 0
    if resultID == 0 {
        return e.NewBadRequestApiError("El ID del resultado no puede ser 0")
    }

    // Verificar si el resultado existe antes de intentar eliminarlo
    result, err := s.resultRepo.GetResultByID(ctx, resultID)
    if err != nil {
        if err == e.NewNotFoundApiError("Result not found") {
            return e.NewNotFoundApiError("El resultado con el ID proporcionado no existe, no se puede eliminar")
        }
        return e.NewInternalServerApiError("Error al verificar la existencia del resultado", err)
    }

    // Imprimir información sobre el resultado que será eliminado (opcional)
    fmt.Printf("Eliminando resultado: ID=%d, DriverID=%d, SessionID=%d\n", result.ID, result.DriverID, result.SessionID)

    // Eliminar el resultado de la base de datos
    if err := s.resultRepo.DeleteResult(ctx, resultID); err != nil {
        return e.NewInternalServerApiError("Error al eliminar el resultado", err)
    }

    return nil
}

// DeleteAllResultsForSession elimina todos los resultados asociados a una sesión específica
func (s *resultService) DeleteAllResultsForSession(ctx context.Context, sessionID int) e.ApiError {
    // Validar que el sessionID no sea 0
    if sessionID == 0 {
        return e.NewBadRequestApiError("El ID de la sesión no puede ser 0")
    }

    // Obtener todos los resultados de la sesión
    results, err := s.resultRepo.GetResultsBySessionID(ctx, sessionID)
    if err != nil {
        return e.NewInternalServerApiError("Error al obtener los resultados de la sesión", err)
    }

    // Verificar si existen resultados para esa sesión
    if len(results) == 0 {
        return e.NewNotFoundApiError("No se encontraron resultados para la sesión especificada")
    }

    // Eliminar cada resultado de la sesión
    for _, result := range results {
        if err := s.resultRepo.DeleteResult(ctx, result.ID); err != nil {
            return e.NewInternalServerApiError(fmt.Sprintf("Error al eliminar el resultado con ID %d", result.ID), err)
        }
    }

    // Retornar éxito si todos los resultados fueron eliminados
    return nil
}

// GetAllResults obtiene todos los resultados de la base de datos
func (s *resultService) GetAllResults(ctx context.Context) ([]dto.ResponseResultDTO, e.ApiError) {
	// Obtener todos los resultados del repositorio
	results, err := s.resultRepo.GetAllResults(ctx)
	if err != nil {
		return nil, e.NewInternalServerApiError("Error al obtener todos los resultados", err)
	}

	// Verificar si no se encontraron resultados
	if len(results) == 0 {
		return nil, e.NewNotFoundApiError("No se encontraron resultados en la base de datos")
	}

	// Crear un slice para almacenar los DTOs de respuesta
	var responseResults []dto.ResponseResultDTO

	// Iterar sobre los resultados obtenidos y convertirlos a DTOs
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
		// Agregar el DTO al slice de respuestas
		responseResults = append(responseResults, response)
	}

	// Retornar los resultados procesados
	return responseResults, nil
}

// GetResultsForDriverAcrossSessions obtiene todos los resultados de un piloto específico en todas las sesiones
// func (s *resultService) GetResultsForDriverAcrossSessions(ctx context.Context, driverID int) ([]dto.ResponseResultDTO, e.ApiError) {
//     // Verificar que el driverID sea válido
//     if driverID == 0 {
//         return nil, e.NewBadRequestApiError("El ID del piloto no puede ser 0")
//     }

//     // Obtener los resultados del piloto en todas las sesiones
//     results, err := s.resultRepo.GetResultsByDriverID(ctx, driverID)
//     if err != nil {
//         return nil, e.NewInternalServerApiError("Error al obtener los resultados del piloto en las sesiones", err)
//     }

//     // Verificar si no se encontraron resultados
//     if len(results) == 0 {
//         return nil, e.NewNotFoundApiError(fmt.Sprintf("No se encontraron resultados para el piloto con ID %d", driverID))
//     }

//     // Crear un slice para almacenar los DTOs de respuesta
//     var responseResults []dto.ResponseResultDTO

//     // Iterar sobre los resultados obtenidos y convertirlos a DTOs
//     for _, result := range results {
//         response := dto.ResponseResultDTO{
//             ID:             result.ID,
//             Position:       result.Position,
//             FastestLapTime: result.FastestLapTime,
//             Driver: dto.ResponseDriverDTO{
//                 ID:          result.Driver.ID,
//                 FirstName:   result.Driver.FirstName,
//                 LastName:    result.Driver.LastName,
//                 FullName:    result.Driver.FullName,
//                 NameAcronym: result.Driver.NameAcronym,
//                 TeamName:    result.Driver.TeamName,
//             },
//             Session: dto.ResponseSessionDTO{
//                 ID:               result.Session.ID,
//                 CircuitShortName: result.Session.CircuitShortName,
//                 CountryName:      result.Session.CountryName,
//                 Location:         result.Session.Location,
//                 SessionName:      result.Session.SessionName,
//                 SessionType:      result.Session.SessionType,
//                 DateStart:        result.Session.DateStart,
//             },
//             CreatedAt: result.CreatedAt,
//             UpdatedAt: result.UpdatedAt,
//         }
//         // Agregar el DTO al slice de respuestas
//         responseResults = append(responseResults, response)
//     }

//     // Retornar los resultados procesados
//     return responseResults, nil
// }

//Esta función obtiene la mejor posición que un piloto ha obtenido en cualquier sesión.
// func (s *resultService) GetBestPositionForDriver(ctx context.Context, driverID int) (dto.ResponseResultDTO, e.ApiError) {
//     // Validar que el driverID sea mayor a 0
//     if driverID == 0 {
//         return dto.ResponseResultDTO{}, e.NewBadRequestApiError("El ID del piloto no puede ser 0")
//     }

//     // Obtener los resultados del piloto
//     results, err := s.resultRepo.GetResultsByDriverID(ctx, driverID)
//     if err != nil {
//         return dto.ResponseResultDTO{}, e.NewInternalServerApiError("Error obteniendo resultados del piloto", err)
//     }

//     // Verificar si no se encontraron resultados
//     if len(results) == 0 {
//         return dto.ResponseResultDTO{}, e.NewNotFoundApiError(fmt.Sprintf("No se encontraron resultados para el piloto con ID %d", driverID))
//     }

//     // Encontrar el mejor resultado (posición más baja)
//     bestResult := results[0]
//     for _, result := range results {
//         if result.Position < bestResult.Position {
//             bestResult = result
//         }
//     }

//     // Convertir el mejor resultado a DTO
//     response := dto.ResponseResultDTO{
//         ID:             bestResult.ID,
//         Position:       bestResult.Position,
//         FastestLapTime: bestResult.FastestLapTime,
//         Driver: dto.ResponseDriverDTO{
//             ID:          bestResult.Driver.ID,
//             FirstName:   bestResult.Driver.FirstName,
//             LastName:    bestResult.Driver.LastName,
//             FullName:    bestResult.Driver.FullName,
//             NameAcronym: bestResult.Driver.NameAcronym,
//             TeamName:    bestResult.Driver.TeamName,
//         },
//         Session: dto.ResponseSessionDTO{
//             ID:               bestResult.Session.ID,
//             CircuitShortName: bestResult.Session.CircuitShortName,
//             CountryName:      bestResult.Session.CountryName,
//             Location:         bestResult.Session.Location,
//             SessionName:      bestResult.Session.SessionName,
//             SessionType:      bestResult.Session.SessionType,
//             DateStart:        bestResult.Session.DateStart,
//         },
//         CreatedAt: bestResult.CreatedAt,
//         UpdatedAt: bestResult.UpdatedAt,
//     }

//     return response, nil
// }

// GetTopNDriversInSession obtiene los mejores N pilotos de una sesión específica.
func (s *resultService) GetTopNDriversInSession(ctx context.Context, sessionID int, n int) ([]dto.TopDriverDTO, e.ApiError) {
    if sessionID == 0 {
        return nil, e.NewBadRequestApiError("El ID de la sesión no puede ser 0")
    }
    if n < 1 {
        return nil, e.NewBadRequestApiError("El número de pilotos a obtener debe ser mayor que 0")
    }

    // Obtener los resultados de la sesión ordenados por posición
    results, err := s.resultRepo.GetResultsOrderedByPosition(ctx, sessionID)
    if err != nil {
        return nil, e.NewInternalServerApiError("Error obteniendo resultados de la sesión", err)
    }

    // Verificar si no se encontraron resultados
    if len(results) == 0 {
        return nil, e.NewNotFoundApiError("No se encontraron resultados para la sesión")
    }

    // Filtrar sólo los que tengan Position != nil (o Status == FINISHED)
    var finishedResults []*model.Result
    for _, r := range results {
        if r.Position != nil {
            finishedResults = append(finishedResults, r)
        }
    }
    if len(finishedResults) == 0 {
        return nil, e.NewNotFoundApiError("Ningún piloto terminó la sesión")
    }

    // Ajustar n si excede
    if n > len(finishedResults) {
        n = len(finishedResults)
    }
    if n > 20 {
        n = 20
    }

    // Construir el slice final
    var topDrivers []dto.TopDriverDTO
    for i := 0; i < n; i++ {
        // Aquí *finishedResults[i].Position es seguro: no es nil
        topDrivers = append(topDrivers, dto.TopDriverDTO{
            Position: *finishedResults[i].Position, 
            DriverID: finishedResults[i].DriverID,
        })
    }
    return topDrivers, nil
}

// GetResultsForSessionByDriverName obtiene los resultados de un piloto en una sesión específica por nombre o acrónimo
// func (s *resultService) GetResultsForSessionByDriverName(ctx context.Context, sessionID int, driverName string) ([]dto.ResponseResultDTO, e.ApiError) {
//     // Validar que el sessionID no sea 0
//     if sessionID == 0 {
//         return nil, e.NewBadRequestApiError("El ID de la sesión no puede ser 0")
//     }

//     // Validar que el nombre del piloto no esté vacío
//     if driverName == "" {
//         return nil, e.NewBadRequestApiError("El nombre del piloto no puede estar vacío")
//     }

//     // Normalizar el nombre del piloto a minúsculas y eliminar espacios adicionales
//     driverName = strings.TrimSpace(strings.ToLower(driverName))

//     // Obtener los resultados de la sesión
//     results, err := s.resultRepo.GetResultsBySessionID(ctx, sessionID)
//     if err != nil {
//         return nil, e.NewInternalServerApiError("Error obteniendo los resultados de la sesión", err)
//     }

//     var filteredResults []dto.ResponseResultDTO
//     for _, result := range results {
//         // Comparar el nombre completo del piloto o su acrónimo con el nombre proporcionado, ignorando mayúsculas
//         if strings.ToLower(result.Driver.FullName) == driverName || strings.ToLower(result.Driver.NameAcronym) == driverName {
//             response := dto.ResponseResultDTO{
//                 ID:             result.ID,
//                 Position:       result.Position,
//                 FastestLapTime: result.FastestLapTime,
//                 Driver: dto.ResponseDriverDTO{
//                     ID:          result.Driver.ID,
//                     FirstName:   result.Driver.FirstName,
//                     LastName:    result.Driver.LastName,
//                     FullName:    result.Driver.FullName,
//                     NameAcronym: result.Driver.NameAcronym,
//                     TeamName:    result.Driver.TeamName,
//                 },
//                 Session: dto.ResponseSessionDTO{
//                     ID:               result.Session.ID,
//                     CircuitShortName: result.Session.CircuitShortName,
//                     CountryName:      result.Session.CountryName,
//                     Location:         result.Session.Location,
//                     SessionName:      result.Session.SessionName,
//                     SessionType:      result.Session.SessionType,
//                     DateStart:        result.Session.DateStart,
//                 },
//                 CreatedAt: result.CreatedAt,
//                 UpdatedAt: result.UpdatedAt,
//             }
//             filteredResults = append(filteredResults, response)
//         }
//     }

//     // Si no se encontraron resultados, devolver un error
//     if len(filteredResults) == 0 {
//         return nil, e.NewNotFoundApiError("No se encontraron resultados para el piloto especificado en esta sesión")
//     }

//     return filteredResults, nil
// }

// GetTotalFastestLapsForDriver cuenta cuántas veces un piloto ha registrado la vuelta más rápida en diferentes sesiones
// func (s *resultService) GetTotalFastestLapsForDriver(ctx context.Context, driverID int) (int, e.ApiError) {
//     // Validar que el driverID no sea 0
//     if driverID == 0 {
//         return 0, e.NewBadRequestApiError("El ID del piloto no puede ser 0")
//     }

//     // Obtener todos los resultados para el piloto dado
//     results, err := s.resultRepo.GetResultsByDriverID(ctx, driverID)
//     if err != nil {
//         return 0, e.NewInternalServerApiError("Error al obtener los resultados del piloto", err)
//     }

//     // Contador de vueltas más rápidas
//     fastestLapCount := 0

//     // Iterar sobre los resultados del piloto
//     for _, result := range results {
//         // Verificar si el piloto tuvo la vuelta más rápida en la sesión
//         fastestLap, err := s.resultRepo.GetFastestLapInSession(ctx, result.SessionID)
//         if err != nil {
//             // Si hay un error al obtener la vuelta más rápida, continuamos con la siguiente sesión
//             continue
//         }
//         // Si el piloto tiene la vuelta más rápida, incrementar el contador
//         if fastestLap.DriverID == driverID {
//             fastestLapCount++
//         }
//     }

//     return fastestLapCount, nil
// }

// GetLastResultForDriver obtiene el último resultado registrado de un piloto en cualquier sesión
// func (s *resultService) GetLastResultForDriver(ctx context.Context, driverID int) (dto.ResponseResultDTO, e.ApiError) {
//     // Validar que el driverID no sea 0
//     if driverID == 0 {
//         return dto.ResponseResultDTO{}, e.NewBadRequestApiError("El ID del piloto no puede ser 0")
//     }

//     // Obtener el último resultado del piloto, ordenando por fecha de creación directamente en la base de datos
//     results, err := s.resultRepo.GetResultsByDriverID(ctx, driverID)
//     if err != nil {
//         return dto.ResponseResultDTO{}, e.NewInternalServerApiError("Error al obtener los resultados del piloto", err)
//     }

//     // Verificar si no se encontraron resultados
//     if len(results) == 0 {
//         return dto.ResponseResultDTO{}, e.NewNotFoundApiError("No se encontraron resultados para el piloto")
//     }

//     // Ordenar los resultados por fecha y devolver el más reciente
//     latestResult := results[len(results)-1]

//     // Convertir el último resultado en un DTO
//     response := dto.ResponseResultDTO{
//         ID:             latestResult.ID,
//         Position:       latestResult.Position,
//         FastestLapTime: latestResult.FastestLapTime,
//         Driver: dto.ResponseDriverDTO{
//             ID:          latestResult.Driver.ID,
//             FirstName:   latestResult.Driver.FirstName,
//             LastName:    latestResult.Driver.LastName,
//             FullName:    latestResult.Driver.FullName,
//             NameAcronym: latestResult.Driver.NameAcronym,
//             TeamName:    latestResult.Driver.TeamName,
//         },
//         Session: dto.ResponseSessionDTO{
//             ID:               latestResult.Session.ID,
//             CircuitShortName: latestResult.Session.CircuitShortName,
//             CountryName:      latestResult.Session.CountryName,
//             Location:         latestResult.Session.Location,
//             SessionName:      latestResult.Session.SessionName,
//             SessionType:      latestResult.Session.SessionType,
//             DateStart:        latestResult.Session.DateStart,
//         },
//         CreatedAt: latestResult.CreatedAt,
//         UpdatedAt: latestResult.UpdatedAt,
//     }

//     return response, nil
// }

func (s *resultService) CreateSessionResultsAdmin(ctx context.Context, bulkRequest dto.CreateBulkResultsDTO) ([]dto.ResponseResultDTO, e.ApiError) {
    // 1. Validar session_id
    if bulkRequest.SessionID == 0 {
        return nil, e.NewBadRequestApiError("El session_id no puede ser 0")
    }

    var resultsToCreate []*model.Result

    // 2. Definir un set de status válidos (o puedes usar "FINISHED" por defecto si Position != nil)
    validStatuses := map[string]bool{"FINISHED": true, "DNF": true, "DNS": true, "DSQ": true}

    // 3. Recorrer cada ítem
    for _, item := range bulkRequest.Results {
        // item.Position es *int en tu `CreateResultItemDTO` si lo ajustaste; si sigue int, 
        //   tendrás que adaptarlo. Asumamos que lo cambiaste a *int.  
        //   OJO: en tu snippet actual `CreateResultItemDTO` está "Position int binding:'required'".
        //   Cambia a Position *int si quieres permitir nulos.

        // Validar si el status no viene, asumimos algo en base a la position
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

        // Validación de la combinación status/position
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
            // "DNF", "DNS", "DSQ" => position debe ser nil
            if item.Position != nil {
                return nil, e.NewBadRequestApiError(
                    fmt.Sprintf("No puede dar Position si el status es %s (driver_id %d)", item.Status, item.DriverID),
                )
            }
        }

        // Validar fastestLapTime
        if item.FastestLapTime != 0 && item.FastestLapTime < 30 {
            return nil, e.NewBadRequestApiError(
                fmt.Sprintf("FastestLapTime inválido para driver_id %d. Debe ser >30 o 0", item.DriverID),
            )
        }

        // Verificar si ya existe un resultado para (driver, session)
        existingResult, _ := s.resultRepo.GetResultByDriverAndSession(ctx, item.DriverID, bulkRequest.SessionID)
        if existingResult != nil {
            return nil, e.NewBadRequestApiError(
                fmt.Sprintf("Ya existe un resultado para driver_id %d en la sesión %d", item.DriverID, bulkRequest.SessionID),
            )
        }

        // Crear el objeto
        newResult := &model.Result{
            SessionID:      bulkRequest.SessionID,
            DriverID:       item.DriverID,
            Position:       item.Position,
            Status:         item.Status,
            FastestLapTime: item.FastestLapTime,
        }
        resultsToCreate = append(resultsToCreate, newResult)
    }

    // 4. Insertar en DB usando un método de repositorio con Transaction
    txErr := s.resultRepo.SessionCreateResultAdmin(ctx, resultsToCreate)
    if txErr != nil {
        return nil, e.NewInternalServerApiError("Error creando resultados masivamente", txErr)
    }

    // 5. Convertir a DTO
    var responseResults []dto.ResponseResultDTO
    for _, r := range resultsToCreate {
        responseResults = append(responseResults, dto.ResponseResultDTO{
            ID:             r.ID,
            Position:       r.Position,
            Status:         r.Status,
            FastestLapTime: r.FastestLapTime,
            Driver: dto.ResponseDriverDTO{ ID: r.DriverID },
            Session: dto.ResponseSessionDTO{ ID: r.SessionID },
            CreatedAt: r.CreatedAt,
            UpdatedAt: r.UpdatedAt,
        })
    }

    return responseResults, nil
}
