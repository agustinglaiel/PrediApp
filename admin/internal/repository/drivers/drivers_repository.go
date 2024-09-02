package repository

import (
	"admin/internal/model/drivers"
	e "admin/pkg/utils"
	"context"

	"gorm.io/gorm"
)

type driverRepository struct {
	db *gorm.DB
}

type DriverRepository interface {
	CreateDriver(ctx context.Context, driver *drivers.Driver) e.ApiError
	GetDriverByID(ctx context.Context, driverID uint) (*drivers.Driver, e.ApiError)
	UpdateDriver(ctx context.Context, driver *drivers.Driver) e.ApiError
	DeleteDriver(ctx context.Context, driverID uint) e.ApiError
	ListDrivers(ctx context.Context) ([]*drivers.Driver, e.ApiError)
	GetDriverByName(ctx context.Context, firstName, lastName string) (*drivers.Driver, e.ApiError)
	GetDriverByNumber(ctx context.Context, driverNumber int) (*drivers.Driver, e.ApiError)
	GetDriversByTeam(ctx context.Context, teamName string) ([]*drivers.Driver, e.ApiError)
	GetDriversByCountry(ctx context.Context, countryCode string) ([]*drivers.Driver, e.ApiError)
	GetDriversByFullName(ctx context.Context, fullName string) ([]*drivers.Driver, e.ApiError)
	GetDriversByAcronym(ctx context.Context, acronym string) ([]*drivers.Driver, e.ApiError)
}

func NewDriverRepository(db *gorm.DB) DriverRepository {
	return &driverRepository{db: db}
}

func (r *driverRepository) CreateDriver(ctx context.Context, driver *drivers.Driver) e.ApiError {
	if err := r.db.WithContext(ctx).Create(driver).Error; err != nil {
		return e.NewInternalServerApiError("Error creando el piloto", err)
	}
	return nil
}

func (r *driverRepository) GetDriverByID(ctx context.Context, driverID uint) (*drivers.Driver, e.ApiError) {
	var driver drivers.Driver
	if err := r.db.WithContext(ctx).First(&driver, driverID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, e.NewNotFoundApiError("Piloto no encontrado")
		}
		return nil, e.NewInternalServerApiError("Error encontrando el piloto", err)
	}
	return &driver, nil
}

func (r *driverRepository) UpdateDriver(ctx context.Context, driver *drivers.Driver) e.ApiError {
	if err := r.db.WithContext(ctx).Save(driver).Error; err != nil {
		return e.NewInternalServerApiError("Error actualizando el piloto", err)
	}
	return nil
}

func (r *driverRepository) DeleteDriver(ctx context.Context, driverID uint) e.ApiError {
	if err := r.db.WithContext(ctx).Where("id = ?", driverID).Delete(&drivers.Driver{}).Error; err != nil {
		return e.NewInternalServerApiError("Error eliminando el piloto", err)
	}
	return nil
}

func (r *driverRepository) ListDrivers(ctx context.Context) ([]*drivers.Driver, e.ApiError) {
	var drivers []*drivers.Driver
	if err := r.db.WithContext(ctx).Find(&drivers).Error; err != nil {
		return nil, e.NewInternalServerApiError("Error listando los pilotos", err)
	}
	return drivers, nil
}

func (r *driverRepository) GetDriverByName(ctx context.Context, firstName, lastName string) (*drivers.Driver, e.ApiError) {
	var driver drivers.Driver
	if err := r.db.WithContext(ctx).Where("first_name = ? AND last_name = ?", firstName, lastName).First(&driver).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No encontrado, no es un error
		}
		return nil, e.NewInternalServerApiError("Error encontrando piloto por nombre y apellido", err)
	}
	return &driver, nil
}

func (r *driverRepository) GetDriverByNumber(ctx context.Context, driverNumber int) (*drivers.Driver, e.ApiError) {
	var driver drivers.Driver
	if err := r.db.WithContext(ctx).Where("driver_number = ?", driverNumber).First(&driver).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No encontrado, no es un error
		}
		return nil, e.NewInternalServerApiError("Error encontrando piloto por número", err)
	}
	return &driver, nil
}

func (r *driverRepository) GetDriversByTeam(ctx context.Context, teamName string) ([]*drivers.Driver, e.ApiError) {
	var drivers []*drivers.Driver
	if err := r.db.WithContext(ctx).Where("team_name = ?", teamName).Find(&drivers).Error; err != nil {
		return nil, e.NewInternalServerApiError("Error encontrando pilotos por equipo", err)
	}
	return drivers, nil
}

func (r *driverRepository) GetDriversByCountry(ctx context.Context, countryCode string) ([]*drivers.Driver, e.ApiError) {
	var drivers []*drivers.Driver
	if err := r.db.WithContext(ctx).Where("country_code = ?", countryCode).Find(&drivers).Error; err != nil {
		return nil, e.NewInternalServerApiError("Error encontrando pilotos por país", err)
	}
	return drivers, nil
}

func (r *driverRepository) GetDriversByFullName(ctx context.Context, fullName string) ([]*drivers.Driver, e.ApiError) {
	var drivers []*drivers.Driver
	if err := r.db.WithContext(ctx).Where("full_name LIKE ?", "%"+fullName+"%").Find(&drivers).Error; err != nil {
		return nil, e.NewInternalServerApiError("Error encontrando pilotos por nombre completo", err)
	}
	return drivers, nil
}

func (r *driverRepository) GetDriversByAcronym(ctx context.Context, acronym string) ([]*drivers.Driver, e.ApiError) {
	var drivers []*drivers.Driver
	if err := r.db.WithContext(ctx).Where("name_acronym LIKE ?", "%"+acronym+"%").Find(&drivers).Error; err != nil {
		return nil, e.NewInternalServerApiError("Error encontrando pilotos por siglas", err)
	}
	return drivers, nil
}
