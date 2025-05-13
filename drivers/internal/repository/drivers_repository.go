package repository

import (
	"context"
	model "drivers/internal/model"
	e "drivers/pkg/utils"
	"fmt"

	"gorm.io/gorm"
)

type driverRepository struct {
	db *gorm.DB
}

type DriverRepository interface {
	CreateDriver(ctx context.Context, driver *model.Driver) e.ApiError
	GetDriverByID(ctx context.Context, driverID int) (*model.Driver, e.ApiError)
	UpdateDriver(ctx context.Context, driver *model.Driver) e.ApiError
	DeleteDriver(ctx context.Context, driverID int) e.ApiError
	ListDrivers(ctx context.Context) ([]*model.Driver, e.ApiError)
	GetDriverByName(ctx context.Context, firstName, lastName string) (*model.Driver, e.ApiError)
	GetDriverByNumber(ctx context.Context, driverNumber int) (*model.Driver, e.ApiError)
	GetDriversByTeam(ctx context.Context, teamName string) ([]*model.Driver, e.ApiError)
	GetDriversByCountry(ctx context.Context, countryCode string) ([]*model.Driver, e.ApiError)
	GetDriversByFullName(ctx context.Context, fullName string) ([]*model.Driver, e.ApiError)
	GetDriversByAcronym(ctx context.Context, acronym string) ([]*model.Driver, e.ApiError)
	BulkInsertDrivers(ctx context.Context, drivers []*model.Driver) e.ApiError
}

func NewDriverRepository(db *gorm.DB) DriverRepository {
	return &driverRepository{db: db}
}

func (r *driverRepository) CreateDriver(ctx context.Context, driver *model.Driver) e.ApiError {
    if err := r.db.WithContext(ctx).Create(driver).Error; err != nil {
        return e.NewInternalServerApiError("Error creando el piloto", err)
    }                 
    return nil
}

func (r *driverRepository) GetDriverByID(ctx context.Context, driverID int) (*model.Driver, e.ApiError) {
	var driver model.Driver
	if err := r.db.WithContext(ctx).First(&driver, driverID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, e.NewNotFoundApiError("Piloto no encontrado")
		}
		return nil, e.NewInternalServerApiError("Error encontrando el piloto", err)
	}
	return &driver, nil
}

func (r *driverRepository) UpdateDriver(ctx context.Context, driver *model.Driver) e.ApiError {
	if err := r.db.WithContext(ctx).Save(driver).Error; err != nil {
		return e.NewInternalServerApiError("Error actualizando el piloto", err)
	}
	return nil
}

func (r *driverRepository) DeleteDriver(ctx context.Context, driverID int) e.ApiError {
	if err := r.db.WithContext(ctx).Where("id = ?", driverID).Delete(&model.Driver{}).Error; err != nil {
		return e.NewInternalServerApiError("Error eliminando el piloto", err)
	}
	return nil
}

func (r *driverRepository) ListDrivers(ctx context.Context) ([]*model.Driver, e.ApiError) {
	var model []*model.Driver
	if err := r.db.WithContext(ctx).Find(&model).Error; err != nil {
		return nil, e.NewInternalServerApiError("Error listando los pilotos", err)
	}
	return model, nil
}

func (r *driverRepository) GetDriverByName(ctx context.Context, firstName string, lastName string) (*model.Driver, e.ApiError) {
	var driver model.Driver
	if err := r.db.WithContext(ctx).Where("first_name = ? AND last_name = ?", firstName, lastName).First(&driver).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No encontrado, no es un error
		}
		return nil, e.NewInternalServerApiError("Error encontrando piloto por nombre y apellido", err)
	}
	return &driver, nil
}

func (r *driverRepository) GetDriverByNumber(ctx context.Context, driverNumber int) (*model.Driver, e.ApiError) {
	var driver model.Driver
	err := r.db.WithContext(ctx).Where("driver_number = ?", driverNumber).First(&driver).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, e.NewNotFoundApiError(fmt.Sprintf("Piloto con número %d no encontrado", driverNumber))
		}
		return nil, e.NewInternalServerApiError("Error encontrando piloto por número", err)
	}
	return &driver, nil
}


func (r *driverRepository) GetDriversByTeam(ctx context.Context, teamName string) ([]*model.Driver, e.ApiError) {
	var model []*model.Driver
	if err := r.db.WithContext(ctx).Where("team_name = ?", teamName).Find(&model).Error; err != nil {
		return nil, e.NewInternalServerApiError("Error encontrando pilotos por equipo", err)
	}
	return model, nil
}

func (r *driverRepository) GetDriversByCountry(ctx context.Context, countryCode string) ([]*model.Driver, e.ApiError) {
	var model []*model.Driver
	if err := r.db.WithContext(ctx).Where("country_code = ?", countryCode).Find(&model).Error; err != nil {
		return nil, e.NewInternalServerApiError("Error encontrando pilotos por país", err)
	}
	return model, nil
}

func (r *driverRepository) GetDriversByFullName(ctx context.Context, fullName string) ([]*model.Driver, e.ApiError) {
	var model []*model.Driver
	if err := r.db.WithContext(ctx).Where("full_name LIKE ?", "%"+fullName+"%").Find(&model).Error; err != nil {
		return nil, e.NewInternalServerApiError("Error encontrando pilotos por nombre completo", err)
	}
	return model, nil
}

func (r *driverRepository) GetDriversByAcronym(ctx context.Context, acronym string) ([]*model.Driver, e.ApiError) {
	var model []*model.Driver
	if err := r.db.WithContext(ctx).Where("name_acronym LIKE ?", "%"+acronym+"%").Find(&model).Error; err != nil {
		return nil, e.NewInternalServerApiError("Error encontrando pilotos por siglas", err)
	}
	return model, nil
}

func (r *driverRepository) BulkInsertDrivers(ctx context.Context, drivers []*model.Driver) e.ApiError {
	// Utilizamos un bulk insert para insertar múltiples pilotos a la vez
	if err := r.db.WithContext(ctx).Create(&drivers).Error; err != nil {
		return e.NewInternalServerApiError("Error inserting drivers in bulk", err)
	}
	return nil
}

