package service

import (
	"context"
	"fmt"
	"results/internal/client"
	"results/internal/dto"
	"results/internal/model"
	"results/internal/repository"
	e "results/pkg/utils"
)

type resultService struct {
	resultRepo repository.ResultRepository
	client     *client.HttpClient
}

type ResultService interface {
	FetchResultsFromExternalAPI(ctx context.Context, sessionId uint) ([]dto.ResponseResultDTO, e.ApiError)
	UpdateResult(ctx context.Context, resultID uint, request dto.UpdateResultDTO) (dto.ResponseResultDTO, e.ApiError)
	GetResultsOrderedByPosition(ctx context.Context, sessionID uint) ([]dto.ResponseResultDTO, e.ApiError)
	GetFastestLapInSession(ctx context.Context, sessionID uint) (dto.ResponseResultDTO, e.ApiError)
	CreateResult(ctx context.Context, request dto.CreateResultDTO) (dto.ResponseResultDTO, e.ApiError)
	GetResultByID(ctx context.Context, resultID uint) (dto.ResponseResultDTO, e.ApiError)
	DeleteResult(ctx context.Context, resultID uint) e.ApiError
	GetAllResults(ctx context.Context) ([]dto.ResponseResultDTO, e.ApiError)
	GetResultsForDriverAcrossSessions(ctx context.Context, driverID uint) ([]dto.ResponseResultDTO, e.ApiError)
	GetBestPositionForDriver(ctx context.Context, driverID uint) (dto.ResponseResultDTO, e.ApiError)
	GetTopNDriversInSession(ctx context.Context, sessionID uint, n int) ([]dto.ResponseResultDTO, e.ApiError)
	DeleteAllResultsForSession(ctx context.Context, sessionID uint) e.ApiError
	GetResultsForSessionByDriverName(ctx context.Context, sessionID uint, driverName string) ([]dto.ResponseResultDTO, e.ApiError)
	GetTotalFastestLapsForDriver(ctx context.Context, driverID uint) (int, e.ApiError)
	GetLastResultForDriver(ctx context.Context, driverID uint) (dto.ResponseResultDTO, e.ApiError)
}

func NewResultService(resultRepo repository.ResultRepository, client *client.HttpClient) ResultService {
	return &resultService{
		resultRepo: resultRepo,
		client:     client,
	}
}

// FetchResultsFromExternalAPI obtiene los resultados de una API externa y los inserta o actualiza en la base de datos
func (s *resultService) FetchResultsFromExternalAPI(ctx context.Context, sessionID uint) ([]dto.ResponseResultDTO, e.ApiError) {
	fmt.Println("Service: Iniciando FetchResultsFromExternalAPI")
	
	// Llamar al microservicio de sessions para obtener el sessionKey
	sessionKey, err := s.client.GetSessionKeyBySessionID(sessionID)
	if err != nil {
		fmt.Println("Service: Error obteniendo sessionKey SERVICE", err)
		return nil, e.NewInternalServerApiError("Error obteniendo session key", err)
	}

	fmt.Println("Service: sessionKey obtenido:", sessionKey)

	// 2. Usar la sessionKey para hacer la solicitud a la API externa y obtener las posiciones
	positions, err := s.client.GetPositions(sessionKey)
	if err != nil {
		fmt.Println("Service: Error obteniendo posiciones de la API externa", err)
		return nil, e.NewInternalServerApiError("Error fetching positions from external API", err)
	}

	// 3. Crear un slice para almacenar los DTOs de respuesta
	var responseResults []dto.ResponseResultDTO

	// 4. Crear un mapa para eliminar duplicados y quedarnos con la última posición registrada para cada piloto
	finalPositions := make(map[int]dto.Position)
	for _, pos := range positions {
		finalPositions[pos.DriverNumber] = pos
	}

	// 5. Obtener las vueltas más rápidas de los pilotos y guardarlas en la base de datos
	for _, pos := range finalPositions {
		// Obtener las vueltas del piloto usando sessionKey y driver_number
		laps, err := s.client.GetLaps(int(sessionKey), pos.DriverNumber)
		if err != nil {
			fmt.Printf("Error obteniendo vueltas para el piloto %d: %v\n", pos.DriverNumber, err)
			continue
		}

		// Encontrar la vuelta más rápida del piloto
		var fastestLap float64
		for _, lap := range laps {
			if fastestLap == 0 || lap.LapDuration < fastestLap {
				fastestLap = lap.LapDuration
			}
		}

		// Llamar al microservicio de drivers para obtener la información completa del piloto
		driverInfo, err := s.client.GetDriverByNumber(pos.DriverNumber)
		if err != nil {
			fmt.Printf("Error obteniendo piloto para el driver_number %d: %v\n", pos.DriverNumber, err)
			continue
		}

		// Verificar si el resultado ya existe
		existingPosition, _ := s.resultRepo.GetDriverPositionInSession(ctx, driverInfo.ID, sessionID)

		// Crear el nuevo resultado o actualizar si ya existe
		newResult := &model.Result{
			SessionID:      sessionID, 
			DriverID:       uint(driverInfo.ID),  // Aquí debes usar el driver_id en lugar del driver_number
			Position:       pos.Position,         // Guardar la posición final del piloto
			FastestLapTime: fastestLap,           // Guardar la vuelta más rápida del piloto
		}

		if existingPosition == 0 {
			// Si no existe, insertarlo en la base de datos
			if err := s.resultRepo.CreateResult(ctx, newResult); err != nil {
				return nil, e.NewInternalServerApiError("Error inserting result into database", err)
			}
		} else {
			// Actualizar el resultado si ya existe
			newResult.Position = pos.Position
			newResult.FastestLapTime = fastestLap
			if err := s.resultRepo.UpdateResult(ctx, newResult); err != nil {
				return nil, e.NewInternalServerApiError("Error updating existing result", err)
			}
		}

		// Convertir el modelo a DTO y agregarlo a la respuesta
		responseResult := dto.ResponseResultDTO{
			ID:             newResult.ID,
			Position:       newResult.Position,
			FastestLapTime: newResult.FastestLapTime,
			Driver: dto.ResponseDriverDTO{
				ID:          driverInfo.ID,  // Usar la información del microservicio de drivers
				FirstName:   driverInfo.FirstName,
				LastName:    driverInfo.LastName,
				FullName:    driverInfo.FullName,
				NameAcronym: driverInfo.NameAcronym,
				TeamName:    driverInfo.TeamName,
			},
			Session: dto.ResponseSessionDTO{
				ID:               newResult.SessionID,
				CircuitShortName: newResult.Session.CircuitShortName,
				CountryName:      newResult.Session.CountryName,
				Location:         newResult.Session.Location,
				SessionName:      newResult.Session.SessionName,
				SessionType:      newResult.Session.SessionType,
				DateStart:        newResult.Session.DateStart,
			},
			CreatedAt: newResult.CreatedAt,
			UpdatedAt: newResult.UpdatedAt,
		}
		responseResults = append(responseResults, responseResult)
	}

	// 6. Retornar los resultados procesados
	return responseResults, nil
}

// UpdateResult actualiza un resultado existente
func (s *resultService) UpdateResult(ctx context.Context, resultID uint, request dto.UpdateResultDTO) (dto.ResponseResultDTO, e.ApiError) {
	result, err := s.resultRepo.GetResultByID(ctx, resultID)
	if err != nil {
		return dto.ResponseResultDTO{}, e.NewBadRequestApiError("error al obtener el resultado por su ID")
	}

	// Actualizar los campos que estén presentes en el DTO
	if request.Position != 0 {
		result.Position = request.Position
	}
	if request.FastestLapTime != 0 {
		result.FastestLapTime = request.FastestLapTime
	}

	if err := s.resultRepo.UpdateResult(ctx, result); err != nil {
		return dto.ResponseResultDTO{}, e.NewInternalServerApiError("Error updating result", err)
	}

	// Convertir el modelo actualizado a DTO
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

	return response, nil
}

// GetResultsOrderedByPosition obtiene los resultados de una sesión específica ordenados por posición
func (s *resultService) GetResultsOrderedByPosition(ctx context.Context, sessionID uint) ([]dto.ResponseResultDTO, e.ApiError) {
	results, err := s.resultRepo.GetResultsOrderedByPosition(ctx, sessionID)
	if err != nil {
		return nil, e.NewBadRequestApiError("Error al mostrar las posiciones de la sesión")
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

	return responseResults, nil
}

// GetFastestLapInSession obtiene el piloto con la vuelta más rápida en una sesión específica
func (s *resultService) GetFastestLapInSession(ctx context.Context, sessionID uint) (dto.ResponseResultDTO, e.ApiError) {
	result, err := s.resultRepo.GetFastestLapInSession(ctx, sessionID)
	if err != nil {
		return dto.ResponseResultDTO{}, e.NewBadRequestApiError("Error al obtener la vuelta mas rápida de la sesión")
	}

	// Convertir el resultado más rápido a DTO
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

	return response, nil
}

// CreateResult crea un nuevo resultado
func (s *resultService) CreateResult(ctx context.Context, request dto.CreateResultDTO) (dto.ResponseResultDTO, e.ApiError) {
	newResult := &model.Result{
		SessionID:      request.SessionID,
		DriverID:       request.DriverID,
		Position:       request.Position,
		FastestLapTime: request.FastestLapTime,
	}

	// Insertar el nuevo resultado en la base de datos
	if err := s.resultRepo.CreateResult(ctx, newResult); err != nil {
		return dto.ResponseResultDTO{}, e.NewInternalServerApiError("Error creating result", err)
	}

	// Convertir el modelo en DTO
	response := dto.ResponseResultDTO{
		ID:             newResult.ID,
		Position:       newResult.Position,
		FastestLapTime: newResult.FastestLapTime,
		Driver: dto.ResponseDriverDTO{
			ID:          newResult.DriverID,
			FirstName:   newResult.Driver.FirstName,
			LastName:    newResult.Driver.LastName,
			FullName:    newResult.Driver.FullName,
			NameAcronym: newResult.Driver.NameAcronym,
			TeamName:    newResult.Driver.TeamName,
		},
		Session: dto.ResponseSessionDTO{
			ID:               newResult.SessionID,
			CircuitShortName: newResult.Session.CircuitShortName,
			CountryName:      newResult.Session.CountryName,
			Location:         newResult.Session.Location,
			SessionName:      newResult.Session.SessionName,
			SessionType:      newResult.Session.SessionType,
			DateStart:        newResult.Session.DateStart,
		},
		CreatedAt: newResult.CreatedAt,
		UpdatedAt: newResult.UpdatedAt,
	}

	return response, nil
}

// GetResultByID obtiene un resultado específico por su ID
func (s *resultService) GetResultByID(ctx context.Context, resultID uint) (dto.ResponseResultDTO, e.ApiError) {
	result, err := s.resultRepo.GetResultByID(ctx, resultID)
	if err != nil {
		return dto.ResponseResultDTO{}, e.NewBadRequestApiError("Error al obtener el resultado por su ID")
	}

	// Convertir el modelo en DTO
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

	return response, nil
}

// DeleteResult elimina un resultado específico
func (s *resultService) DeleteResult(ctx context.Context, resultID uint) e.ApiError {
    // Verificar si el resultado existe
    _, err := s.resultRepo.GetResultByID(ctx, resultID)
    if err != nil {
        if err == e.NewNotFoundApiError("Result not found") {
            return e.NewNotFoundApiError("Result not found, cannot delete")
        }
        return e.NewInternalServerApiError("Error checking result existence", err)
    }

    // Eliminar el resultado de la base de datos
    if err := s.resultRepo.DeleteResult(ctx, resultID); err != nil {
        return e.NewInternalServerApiError("Error deleting result", err)
    }

    return nil
}

// GetAllResults obtiene todos los resultados de la base de datos
func (s *resultService) GetAllResults(ctx context.Context) ([]dto.ResponseResultDTO, e.ApiError) {
	results, err := s.resultRepo.GetAllResults(ctx)
	if err != nil {
		return nil, err
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

	return responseResults, nil
}

//Esta función obtiene todos los resultados de un piloto específico en todas las sesiones.
func (s *resultService) GetResultsForDriverAcrossSessions(ctx context.Context, driverID uint) ([]dto.ResponseResultDTO, e.ApiError) {
    results, err := s.resultRepo.GetResultsByDriverID(ctx, driverID)
    if err != nil {
        return nil, e.NewInternalServerApiError("Error fetching results for driver across sessions", err)
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

    return responseResults, nil
}

//Esta función obtiene la mejor posición que un piloto ha obtenido en cualquier sesión.
//Propósito: Puede ser útil para mostrar el rendimiento máximo de un piloto a lo largo de las temporadas.
func (s *resultService) GetBestPositionForDriver(ctx context.Context, driverID uint) (dto.ResponseResultDTO, e.ApiError) {
    results, err := s.resultRepo.GetResultsByDriverID(ctx, driverID)
    if err != nil {
        return dto.ResponseResultDTO{}, e.NewInternalServerApiError("Error fetching results for driver", err)
    }

    var bestResult *model.Result
    for _, result := range results {
        if bestResult == nil || result.Position < bestResult.Position {
            bestResult = result
        }
    }

    if bestResult == nil {
        return dto.ResponseResultDTO{}, e.NewNotFoundApiError("No results found for driver")
    }

    response := dto.ResponseResultDTO{
        ID:             bestResult.ID,
        Position:       bestResult.Position,
        FastestLapTime: bestResult.FastestLapTime,
        Driver: dto.ResponseDriverDTO{
            ID:          bestResult.Driver.ID,
            FirstName:   bestResult.Driver.FirstName,
            LastName:    bestResult.Driver.LastName,
            FullName:    bestResult.Driver.FullName,
            NameAcronym: bestResult.Driver.NameAcronym,
            TeamName:    bestResult.Driver.TeamName,
        },
        Session: dto.ResponseSessionDTO{
            ID:               bestResult.Session.ID,
            CircuitShortName: bestResult.Session.CircuitShortName,
            CountryName:      bestResult.Session.CountryName,
            Location:         bestResult.Session.Location,
            SessionName:      bestResult.Session.SessionName,
            SessionType:      bestResult.Session.SessionType,
            DateStart:        bestResult.Session.DateStart,
        },
        CreatedAt: bestResult.CreatedAt,
        UpdatedAt: bestResult.UpdatedAt,
    }

    return response, nil
}

//Esta función obtiene los mejores N pilotos de una sesión específica.
func (s *resultService) GetTopNDriversInSession(ctx context.Context, sessionID uint, n int) ([]dto.ResponseResultDTO, e.ApiError) {
    results, err := s.resultRepo.GetResultsOrderedByPosition(ctx, sessionID)
    if err != nil {
        return nil, e.NewInternalServerApiError("Error fetching session results", err)
    }

    if len(results) == 0 {
        return nil, e.NewNotFoundApiError("No results found for session")
    }

    if n > len(results) {
        n = len(results) // Ajustar si N es mayor que el número total de pilotos
    }

    var responseResults []dto.ResponseResultDTO
    for i := 0; i < n; i++ {
        result := results[i]
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

    return responseResults, nil
}

//Esta función elimina todos los resultados asociados a una sesión específica.
func (s *resultService) DeleteAllResultsForSession(ctx context.Context, sessionID uint) e.ApiError {
    results, err := s.resultRepo.GetResultsBySessionID(ctx, sessionID)
    if err != nil {
        return e.NewInternalServerApiError("Error fetching results for session", err)
    }

    for _, result := range results {
        if err := s.resultRepo.DeleteResult(ctx, result.ID); err != nil {
            return e.NewInternalServerApiError("Error deleting result", err)
        }
    }

    return nil
}

func (s *resultService) GetResultsForSessionByDriverName(ctx context.Context, sessionID uint, driverName string) ([]dto.ResponseResultDTO, e.ApiError) {
    results, err := s.resultRepo.GetResultsBySessionID(ctx, sessionID)
    if err != nil {
        return nil, e.NewInternalServerApiError("Error fetching results for session", err)
    }

    var filteredResults []dto.ResponseResultDTO
    for _, result := range results {
        if result.Driver.FullName == driverName || result.Driver.NameAcronym == driverName {
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
            filteredResults = append(filteredResults, response)
        }
    }

    if len(filteredResults) == 0 {
        return nil, e.NewNotFoundApiError("No results found for the specified driver in this session")
    }

    return filteredResults, nil
}

func (s *resultService) GetTotalFastestLapsForDriver(ctx context.Context, driverID uint) (int, e.ApiError) {
    results, err := s.resultRepo.GetResultsByDriverID(ctx, driverID)
    if err != nil {
        return 0, e.NewInternalServerApiError("Error fetching results for driver", err)
    }

    fastestLapCount := 0
    for _, result := range results {
        fastestLap, err := s.resultRepo.GetFastestLapInSession(ctx, result.SessionID)
        if err == nil && fastestLap.DriverID == driverID {
            fastestLapCount++
        }
    }

    return fastestLapCount, nil
}

func (s *resultService) GetLastResultForDriver(ctx context.Context, driverID uint) (dto.ResponseResultDTO, e.ApiError) {
    results, err := s.resultRepo.GetResultsByDriverID(ctx, driverID)
    if err != nil {
        return dto.ResponseResultDTO{}, e.NewInternalServerApiError("Error fetching results for driver", err)
    }

    if len(results) == 0 {
        return dto.ResponseResultDTO{}, e.NewNotFoundApiError("No results found for driver")
    }

    // Ordenar los resultados por fecha y devolver el último
    latestResult := results[len(results)-1]
    response := dto.ResponseResultDTO{
        ID:             latestResult.ID,
        Position:       latestResult.Position,
        FastestLapTime: latestResult.FastestLapTime,
        Driver: dto.ResponseDriverDTO{
            ID:          latestResult.Driver.ID,
            FirstName:   latestResult.Driver.FirstName,
            LastName:    latestResult.Driver.LastName,
            FullName:    latestResult.Driver.FullName,
            NameAcronym: latestResult.Driver.NameAcronym,
            TeamName:    latestResult.Driver.TeamName,
        },
        Session: dto.ResponseSessionDTO{
            ID:               latestResult.Session.ID,
            CircuitShortName: latestResult.Session.CircuitShortName,
            CountryName:      latestResult.Session.CountryName,
            Location:         latestResult.Session.Location,
            SessionName:      latestResult.Session.SessionName,
            SessionType:      latestResult.Session.SessionType,
            DateStart:        latestResult.Session.DateStart,
        },
        CreatedAt: latestResult.CreatedAt,
        UpdatedAt: latestResult.UpdatedAt,
    }

    return response, nil
}