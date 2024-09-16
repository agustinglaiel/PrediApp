package service

import (
	"context"
	dto "drivers/internal/dto/drivers"
	model "drivers/internal/model/drivers"
	repository "drivers/internal/repository/drivers"
	e "drivers/pkg/utils"
)

type driverService struct {
	driverRepo repository.DriverRepository
}

type DriverService interface {
	CreateDriver(ctx context.Context, request dto.CreateDriverDTO) (dto.ResponseDriverDTO, e.ApiError)
	GetDriverByID(ctx context.Context, driverID uint) (dto.ResponseDriverDTO, e.ApiError)
	UpdateDriver(ctx context.Context, driverID uint, request dto.UpdateDriverDTO) (dto.ResponseDriverDTO, e.ApiError)
	DeleteDriver(ctx context.Context, driverID uint) e.ApiError
	ListDrivers(ctx context.Context) ([]dto.ResponseDriverDTO, e.ApiError)
	ListDriversByTeam(ctx context.Context, teamName string) ([]dto.ResponseDriverDTO, e.ApiError)
	ListDriversByCountry(ctx context.Context, countryCode string) ([]dto.ResponseDriverDTO, e.ApiError)
	ListDriversByFullName(ctx context.Context, fullName string) ([]dto.ResponseDriverDTO, e.ApiError)
	ListDriversByAcronym(ctx context.Context, acronym string) ([]dto.ResponseDriverDTO, e.ApiError)
}

func NewDriverService(driverRepo repository.DriverRepository) DriverService {
	return &driverService{
		driverRepo: driverRepo,
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
		BroadcastName:  request.BroadcastName,
		CountryCode:    request.CountryCode,
		DriverNumber:   request.DriverNumber,
		FirstName:      request.FirstName,
		LastName:       request.LastName,
		FullName:       request.FullName,
		NameAcronym:    request.NameAcronym,
		TeamName:       request.TeamName,
	}

	if err := s.driverRepo.CreateDriver(ctx, newDriver); err != nil {
		return dto.ResponseDriverDTO{}, e.NewInternalServerApiError("Error creando el piloto", err)
	}

	// Convert Model to Response DTO
	response := dto.ResponseDriverDTO{
		ID:             newDriver.ID,
		BroadcastName:  newDriver.BroadcastName,
		CountryCode:    newDriver.CountryCode,
		DriverNumber:   newDriver.DriverNumber,
		FirstName:      newDriver.FirstName,
		LastName:       newDriver.LastName,
		FullName:       newDriver.FullName,
		NameAcronym:    newDriver.NameAcronym,
		TeamName:       newDriver.TeamName,
	}

	return response, nil
}

func (s *driverService) GetDriverByID(ctx context.Context, driverID uint) (dto.ResponseDriverDTO, e.ApiError) {
	driver, err := s.driverRepo.GetDriverByID(ctx, driverID)
	if err != nil {
		return dto.ResponseDriverDTO{}, err
	}

	// Convert Model to Response DTO
	response := dto.ResponseDriverDTO{
		ID:             driver.ID,
		BroadcastName:  driver.BroadcastName,
		CountryCode:    driver.CountryCode,
		DriverNumber:   driver.DriverNumber,
		FirstName:      driver.FirstName,
		LastName:       driver.LastName,
		FullName:       driver.FullName,
		NameAcronym:    driver.NameAcronym,
		TeamName:       driver.TeamName,
	}

	return response, nil
}

func (s *driverService) UpdateDriver(ctx context.Context, driverID uint, request dto.UpdateDriverDTO) (dto.ResponseDriverDTO, e.ApiError) {
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
	if request.TeamName != "" {
		driver.TeamName = request.TeamName
	}

	if err := s.driverRepo.UpdateDriver(ctx, driver); err != nil {
		return dto.ResponseDriverDTO{}, e.NewInternalServerApiError("Error actualizando el piloto", err)
	}

	// Convert Model to Response DTO
	response := dto.ResponseDriverDTO{
		ID:             driver.ID,
		BroadcastName:  driver.BroadcastName,
		CountryCode:    driver.CountryCode,
		DriverNumber:   driver.DriverNumber,
		FirstName:      driver.FirstName,
		LastName:       driver.LastName,
		FullName:       driver.FullName,
		NameAcronym:    driver.NameAcronym,
		TeamName:       driver.TeamName,
	}

	return response, nil
}

func (s *driverService) DeleteDriver(ctx context.Context, driverID uint) e.ApiError {
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
			ID:             driver.ID,
			BroadcastName:  driver.BroadcastName,
			CountryCode:    driver.CountryCode,
			DriverNumber:   driver.DriverNumber,
			FirstName:      driver.FirstName,
			LastName:       driver.LastName,
			FullName:       driver.FullName,
			NameAcronym:    driver.NameAcronym,
			TeamName:       driver.TeamName,
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
			ID:             driver.ID,
			BroadcastName:  driver.BroadcastName,
			CountryCode:    driver.CountryCode,
			DriverNumber:   driver.DriverNumber,
			FirstName:      driver.FirstName,
			LastName:       driver.LastName,
			FullName:       driver.FullName,
			NameAcronym:    driver.NameAcronym,
			TeamName:       driver.TeamName,
		})
	}

	return response, nil
}

func (s *driverService) ListDriversByCountry(ctx context.Context, countryCode string) ([]dto.ResponseDriverDTO, e.ApiError) {
	drivers, err := s.driverRepo.GetDriversByCountry(ctx, countryCode)
	if err != nil {
		return nil, err
	}

	var response []dto.ResponseDriverDTO
	for _, driver := range drivers {
		response = append(response, dto.ResponseDriverDTO{
			ID:             driver.ID,
			BroadcastName:  driver.BroadcastName,
			CountryCode:    driver.CountryCode,
			DriverNumber:   driver.DriverNumber,
			FirstName:      driver.FirstName,
			LastName:       driver.LastName,
			FullName:       driver.FullName,
			NameAcronym:    driver.NameAcronym,
			TeamName:       driver.TeamName,
		})
	}

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
			ID:             driver.ID,
			BroadcastName:  driver.BroadcastName,
			CountryCode:    driver.CountryCode,
			DriverNumber:   driver.DriverNumber,
			FirstName:      driver.FirstName,
			LastName:       driver.LastName,
			FullName:       driver.FullName,
			NameAcronym:    driver.NameAcronym,
			TeamName:       driver.TeamName,
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
			ID:             driver.ID,
			BroadcastName:  driver.BroadcastName,
			CountryCode:    driver.CountryCode,
			DriverNumber:   driver.DriverNumber,
			FirstName:      driver.FirstName,
			LastName:       driver.LastName,
			FullName:       driver.FullName,
			NameAcronym:    driver.NameAcronym,
			TeamName:       driver.TeamName,
		})
	}

	return response, nil
}