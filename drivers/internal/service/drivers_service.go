package service

import (
	"context"
	"fmt"
	"log"

	model "prediapp.local/db/model"
	"prediapp.local/drivers/internal/client"
	dto "prediapp.local/drivers/internal/dto"
	repository "prediapp.local/drivers/internal/repository"
	e "prediapp.local/drivers/pkg/utils"
)

type driverService struct {
	driverRepo repository.DriverRepository
	client     *client.HttpClient
}

type DriverService interface {
	CreateDriver(ctx context.Context, request dto.CreateDriverDTO) (dto.ResponseDriverDTO, e.ApiError)
	GetDriverByID(ctx context.Context, driverID int) (dto.ResponseDriverDTO, e.ApiError)
	UpdateDriver(ctx context.Context, driverID int, request dto.UpdateDriverDTO) (dto.ResponseDriverDTO, e.ApiError)
	DeleteDriver(ctx context.Context, driverID int) e.ApiError
	ListDrivers(ctx context.Context) ([]dto.ResponseDriverDTO, e.ApiError)
	ListDriversByTeam(ctx context.Context, teamName string) ([]dto.ResponseDriverDTO, e.ApiError)
	ListDriversByCountry(ctx context.Context, countryCode string) ([]dto.ResponseDriverDTO, e.ApiError)
	ListDriversByFullName(ctx context.Context, fullName string) ([]dto.ResponseDriverDTO, e.ApiError)
	ListDriversByAcronym(ctx context.Context, acronym string) ([]dto.ResponseDriverDTO, e.ApiError)
	FetchAllDriversFromExternalAPI(ctx context.Context) ([]dto.ResponseDriverDTO, e.ApiError)
	GetDriverByNumber(ctx context.Context, driverNumber int) (dto.ResponseDriverDTO, e.ApiError)
	GetDriverByFirstAndLastName(ctx context.Context, firstName, lastName string) (dto.ResponseDriverDTO, e.ApiError)
}

func NewDriverService(driverRepo repository.DriverRepository, client *client.HttpClient) DriverService {
	return &driverService{
		driverRepo: driverRepo,
		client:     client, // Pasar el cliente HTTP
	}
}

func (s *driverService) CreateDriver(ctx context.Context, request dto.CreateDriverDTO) (dto.ResponseDriverDTO, e.ApiError) {
	// Verificar si un piloto con el mismo nombre y apellido ya existe
	existingDriver, _ := s.driverRepo.GetDriverByName(ctx, request.FirstName, request.LastName)
	if existingDriver != nil {
		return dto.ResponseDriverDTO{}, e.NewBadRequestApiError("Ya existe un piloto con el mismo nombre y apellido")
	}

	// Verificar si el número de piloto ya está en uso
	driverWithNumber, _ := s.driverRepo.GetDriverByNumber(ctx, request.DriverNumber)
	if driverWithNumber != nil {
		return dto.ResponseDriverDTO{}, e.NewBadRequestApiError("Ya existe un piloto con el mismo número")
	}

	// Convert DTO to Model
	newDriver := &model.Driver{
		BroadcastName: request.BroadcastName,
		CountryCode:   request.CountryCode,
		DriverNumber:  request.DriverNumber,
		FirstName:     request.FirstName,
		LastName:      request.LastName,
		FullName:      request.FullName,
		NameAcronym:   request.NameAcronym,
		HeadshotURL:   request.HeadshotURL,
		TeamName:      request.TeamName,
		Activo:        request.Activo,
	}

	if err := s.driverRepo.CreateDriver(ctx, newDriver); err != nil {
		return dto.ResponseDriverDTO{}, e.NewInternalServerApiError("Error creando el piloto", err)
	}

	// Convert Model to Response DTO
	response := dto.ResponseDriverDTO{
		ID:            newDriver.ID,
		BroadcastName: newDriver.BroadcastName,
		CountryCode:   newDriver.CountryCode,
		DriverNumber:  newDriver.DriverNumber,
		FirstName:     newDriver.FirstName,
		LastName:      newDriver.LastName,
		FullName:      newDriver.FullName,
		NameAcronym:   newDriver.NameAcronym,
		HeadshotURL:   newDriver.HeadshotURL,
		TeamName:      newDriver.TeamName,
		Activo:        newDriver.Activo,
	}

	return response, nil
}

func (s *driverService) GetDriverByID(ctx context.Context, driverID int) (dto.ResponseDriverDTO, e.ApiError) {
	driver, err := s.driverRepo.GetDriverByID(ctx, driverID)
	if err != nil {
		return dto.ResponseDriverDTO{}, err
	}

	// Convert Model to Response DTO
	response := dto.ResponseDriverDTO{
		ID:            driver.ID,
		BroadcastName: driver.BroadcastName,
		CountryCode:   driver.CountryCode,
		DriverNumber:  driver.DriverNumber,
		FirstName:     driver.FirstName,
		LastName:      driver.LastName,
		FullName:      driver.FullName,
		NameAcronym:   driver.NameAcronym,
		HeadshotURL:   driver.HeadshotURL,
		TeamName:      driver.TeamName,
		Activo:        driver.Activo,
	}

	return response, nil
}

func (s *driverService) UpdateDriver(ctx context.Context, driverID int, request dto.UpdateDriverDTO) (dto.ResponseDriverDTO, e.ApiError) {
	driver, err := s.driverRepo.GetDriverByID(ctx, driverID)
	if err != nil {
		return dto.ResponseDriverDTO{}, err
	}

	// Update only the fields present in the DTO
	if request.BroadcastName != "" {
		driver.BroadcastName = request.BroadcastName
	}
	if request.CountryCode != "" {
		driver.CountryCode = request.CountryCode
	}
	if request.DriverNumber != 0 {
		driver.DriverNumber = request.DriverNumber
	}
	if request.FirstName != "" {
		driver.FirstName = request.FirstName
	}
	if request.LastName != "" {
		driver.LastName = request.LastName
	}
	if request.FullName != "" {
		driver.FullName = request.FullName
	}
	if request.NameAcronym != "" {
		driver.NameAcronym = request.NameAcronym
	}
	if request.HeadshotURL != "" {
		driver.HeadshotURL = request.HeadshotURL
	}
	if request.TeamName != "" {
		driver.TeamName = request.TeamName
	}
	if request.Activo != nil {
		driver.Activo = *request.Activo
	}

	if err := s.driverRepo.UpdateDriver(ctx, driver); err != nil {
		return dto.ResponseDriverDTO{}, e.NewInternalServerApiError("Error actualizando el piloto", err)
	}

	// Convert Model to Response DTO
	response := dto.ResponseDriverDTO{
		ID:            driver.ID,
		BroadcastName: driver.BroadcastName,
		CountryCode:   driver.CountryCode,
		DriverNumber:  driver.DriverNumber,
		FirstName:     driver.FirstName,
		LastName:      driver.LastName,
		FullName:      driver.FullName,
		NameAcronym:   driver.NameAcronym,
		HeadshotURL:   driver.HeadshotURL,
		TeamName:      driver.TeamName,
		Activo:        driver.Activo,
	}

	return response, nil
}

func (s *driverService) DeleteDriver(ctx context.Context, driverID int) e.ApiError {
	// Check if the driver exists before attempting to delete
	driver, err := s.driverRepo.GetDriverByID(ctx, driverID)
	if err != nil {
		return err
	}

	// Delete the driver
	if err := s.driverRepo.DeleteDriver(ctx, driver.ID); err != nil {
		return e.NewInternalServerApiError("Error eliminando el piloto", err)
	}

	return nil
}

func (s *driverService) ListDrivers(ctx context.Context) ([]dto.ResponseDriverDTO, e.ApiError) {
	drivers, err := s.driverRepo.ListDrivers(ctx)
	if err != nil {
		return nil, err
	}

	// Convert Model list to Response DTO list
	var response []dto.ResponseDriverDTO
	for _, driver := range drivers {
		response = append(response, dto.ResponseDriverDTO{
			ID:            driver.ID,
			BroadcastName: driver.BroadcastName,
			CountryCode:   driver.CountryCode,
			DriverNumber:  driver.DriverNumber,
			FirstName:     driver.FirstName,
			LastName:      driver.LastName,
			FullName:      driver.FullName,
			NameAcronym:   driver.NameAcronym,
			HeadshotURL:   driver.HeadshotURL,
			TeamName:      driver.TeamName,
			Activo:        driver.Activo,
		})
	}

	return response, nil
}

func (s *driverService) ListDriversByTeam(ctx context.Context, teamName string) ([]dto.ResponseDriverDTO, e.ApiError) {
	drivers, err := s.driverRepo.GetDriversByTeam(ctx, teamName)
	if err != nil {
		return nil, err
	}

	var response []dto.ResponseDriverDTO
	for _, driver := range drivers {
		response = append(response, dto.ResponseDriverDTO{
			ID:            driver.ID,
			BroadcastName: driver.BroadcastName,
			CountryCode:   driver.CountryCode,
			DriverNumber:  driver.DriverNumber,
			FirstName:     driver.FirstName,
			LastName:      driver.LastName,
			FullName:      driver.FullName,
			NameAcronym:   driver.NameAcronym,
			HeadshotURL:   driver.HeadshotURL,
			TeamName:      driver.TeamName,
			Activo:        driver.Activo,
		})
	}

	return response, nil
}

func (s *driverService) ListDriversByCountry(ctx context.Context, countryCode string) ([]dto.ResponseDriverDTO, e.ApiError) {
	drivers, err := s.driverRepo.GetDriversByCountry(ctx, countryCode)
	if err != nil {
		return nil, err
	}
	log.Printf("Drivers: %v", drivers)

	var response []dto.ResponseDriverDTO
	for _, driver := range drivers {
		response = append(response, dto.ResponseDriverDTO{
			ID:            driver.ID,
			BroadcastName: driver.BroadcastName,
			CountryCode:   driver.CountryCode,
			DriverNumber:  driver.DriverNumber,
			FirstName:     driver.FirstName,
			LastName:      driver.LastName,
			FullName:      driver.FullName,
			NameAcronym:   driver.NameAcronym,
			HeadshotURL:   driver.HeadshotURL,
			TeamName:      driver.TeamName,
			Activo:        driver.Activo,
		})
	}
	log.Printf("Drivers despues: %v", response)

	return response, nil
}

func (s *driverService) ListDriversByFullName(ctx context.Context, fullName string) ([]dto.ResponseDriverDTO, e.ApiError) {
	drivers, err := s.driverRepo.GetDriversByFullName(ctx, fullName)
	if err != nil {
		return nil, err
	}

	var response []dto.ResponseDriverDTO
	for _, driver := range drivers {
		response = append(response, dto.ResponseDriverDTO{
			ID:            driver.ID,
			BroadcastName: driver.BroadcastName,
			CountryCode:   driver.CountryCode,
			DriverNumber:  driver.DriverNumber,
			FirstName:     driver.FirstName,
			LastName:      driver.LastName,
			FullName:      driver.FullName,
			NameAcronym:   driver.NameAcronym,
			HeadshotURL:   driver.HeadshotURL,
			TeamName:      driver.TeamName,
			Activo:        driver.Activo,
		})
	}

	return response, nil
}

func (s *driverService) ListDriversByAcronym(ctx context.Context, acronym string) ([]dto.ResponseDriverDTO, e.ApiError) {
	drivers, err := s.driverRepo.GetDriversByAcronym(ctx, acronym)
	if err != nil {
		return nil, err
	}

	var response []dto.ResponseDriverDTO
	for _, driver := range drivers {
		response = append(response, dto.ResponseDriverDTO{
			ID:            driver.ID,
			BroadcastName: driver.BroadcastName,
			CountryCode:   driver.CountryCode,
			DriverNumber:  driver.DriverNumber,
			FirstName:     driver.FirstName,
			LastName:      driver.LastName,
			FullName:      driver.FullName,
			NameAcronym:   driver.NameAcronym,
			HeadshotURL:   driver.HeadshotURL,
			TeamName:      driver.TeamName,
			Activo:        driver.Activo,
		})
	}

	return response, nil
}

func (s *driverService) FetchAllDriversFromExternalAPI(ctx context.Context) ([]dto.ResponseDriverDTO, e.ApiError) {
	// 1. Obtener los pilotos desde la API externa
	drivers, err := s.client.GetAllDriversFromExternalAPI()
	if err != nil {
		return nil, e.NewInternalServerApiError("Error fetching drivers from external API", err)
	}

	// 2. Eliminar duplicados por 'FirstName y LastName'
	uniqueDrivers := uniqueDrivers(drivers)

	// 3. Obtener todos los pilotos existentes en la base de datos
	existingDrivers, err := s.driverRepo.ListDrivers(ctx)
	if err != nil {
		return nil, e.NewInternalServerApiError("Error fetching drivers from database", err)
	}

	// 4. Crear un mapa para los pilotos existentes (para verificar más rápido)
	existingDriverMap := make(map[string]bool)
	for _, driver := range existingDrivers {
		key := driver.FirstName + "|" + driver.LastName
		existingDriverMap[key] = true
	}

	// 5. Preparar la lista de pilotos nuevos para insertar
	var newDrivers []*model.Driver
	for _, driver := range uniqueDrivers {
		key := driver.FirstName + "|" + driver.LastName
		if _, exists := existingDriverMap[key]; !exists {
			newDriver := &model.Driver{
				BroadcastName: driver.BroadcastName,
				CountryCode:   driver.CountryCode,
				DriverNumber:  driver.DriverNumber,
				FirstName:     driver.FirstName,
				LastName:      driver.LastName,
				FullName:      driver.FullName,
				NameAcronym:   driver.NameAcronym,
				HeadshotURL:   driver.HeadshotURL,
				TeamName:      driver.TeamName,
				Activo:        true,
			}
			newDrivers = append(newDrivers, newDriver)
			existingDriverMap[key] = true
		}
	}

	// 6. Insertar los pilotos nuevos usando una transacción en el repositorio
	insertedModels, err := s.driverRepo.CreateDriversTransaction(ctx, newDrivers)
	if err != nil {
		return nil, e.NewInternalServerApiError("Error inserting new drivers", err)
	}

	// 8. Convertir los modelos insertados a DTOs de respuesta
	var insertedDrivers []dto.ResponseDriverDTO
	for _, driver := range insertedModels {
		insertedDrivers = append(insertedDrivers, dto.ResponseDriverDTO{
			ID:            driver.ID,
			BroadcastName: driver.BroadcastName,
			CountryCode:   driver.CountryCode,
			DriverNumber:  driver.DriverNumber,
			FirstName:     driver.FirstName,
			LastName:      driver.LastName,
			FullName:      driver.FullName,
			NameAcronym:   driver.NameAcronym,
			HeadshotURL:   driver.HeadshotURL,
			TeamName:      driver.TeamName,
			Activo:        false,
		})
	}

	return insertedDrivers, nil
}

// Función auxiliar para eliminar pilotos duplicados basados en 'First name y Last name'
func uniqueDrivers(drivers []dto.ResponseDriverDTO) []dto.ResponseDriverDTO {
	seen := make(map[string]bool)
	unique := []dto.ResponseDriverDTO{}

	for _, driver := range drivers {
		// Crear una clave única combinando FirstName y LastName
		key := driver.FirstName + "|" + driver.LastName
		if _, ok := seen[key]; !ok {
			seen[key] = true
			unique = append(unique, driver)
		}
	}

	return unique
}

func (s *driverService) GetDriverByNumber(ctx context.Context, driverNumber int) (dto.ResponseDriverDTO, e.ApiError) {
	driver, err := s.driverRepo.GetDriverByNumber(ctx, driverNumber)
	if err != nil {
		return dto.ResponseDriverDTO{}, err
	}

	// Verificar si el driver es nil
	if driver == nil {
		return dto.ResponseDriverDTO{}, e.NewNotFoundApiError(fmt.Sprintf("Driver with number %d not found", driverNumber))
	}

	// Convert Model to Response DTO
	response := dto.ResponseDriverDTO{
		ID:            driver.ID,
		BroadcastName: driver.BroadcastName,
		CountryCode:   driver.CountryCode,
		DriverNumber:  driver.DriverNumber,
		FirstName:     driver.FirstName,
		LastName:      driver.LastName,
		FullName:      driver.FullName,
		NameAcronym:   driver.NameAcronym,
		HeadshotURL:   driver.HeadshotURL,
		TeamName:      driver.TeamName,
		Activo:        driver.Activo,
	}

	return response, nil
}

func (s *driverService) GetDriverByFirstAndLastName(ctx context.Context, firstName, lastName string) (dto.ResponseDriverDTO, e.ApiError) {
	// Llamar al repositorio para buscar el piloto por first_name y last_name
	driver, err := s.driverRepo.GetDriverByName(ctx, firstName, lastName)
	if err != nil {
		return dto.ResponseDriverDTO{}, err
	}

	// Verificar si el driver es nil (no encontrado)
	if driver == nil {
		return dto.ResponseDriverDTO{}, e.NewNotFoundApiError(fmt.Sprintf("Driver with name %s %s not found", firstName, lastName))
	}

	// Convertir el modelo a DTO de respuesta
	response := dto.ResponseDriverDTO{
		ID:            driver.ID,
		BroadcastName: driver.BroadcastName,
		CountryCode:   driver.CountryCode,
		DriverNumber:  driver.DriverNumber,
		FirstName:     driver.FirstName,
		LastName:      driver.LastName,
		FullName:      driver.FullName,
		NameAcronym:   driver.NameAcronym,
		HeadshotURL:   driver.HeadshotURL,
		TeamName:      driver.TeamName,
		Activo:        driver.Activo,
	}

	return response, nil
}
