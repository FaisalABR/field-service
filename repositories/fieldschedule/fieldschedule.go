package repositories

import (
	"context"
	"errors"
	errorWrap "field-service/common/error"
	"field-service/constants"
	errConstants "field-service/constants/error"
	errField "field-service/constants/error/field"
	"field-service/domain/dto"
	"field-service/domain/models"
	"fmt"

	"gorm.io/gorm"
)

type FieldScheduleRepository struct {
	db *gorm.DB
}

type IFieldScheduleRepository interface {
	FindAllWithPagination(context.Context, *dto.FieldScheduleRequestParam) ([]models.FieldSchedule, int64, error)
	FindAllByFieldIDAndDate(context.Context, int, string) ([]models.FieldSchedule, error)
	FindByUUID(context.Context, string) (*models.FieldSchedule, error)
	FindByDateAndTimeID(context.Context, string, int, int) (*models.FieldSchedule, error)
	Create(context.Context, []models.FieldSchedule) error
	Update(context.Context, string, *models.FieldSchedule) (*models.FieldSchedule, error)
	UpdateStatus(context.Context, constants.FieldScheduleStatus, string) error
	Delete(context.Context, string) error
}

func NewFieldScheduleRepository(db *gorm.DB) IFieldScheduleRepository {
	return &FieldScheduleRepository{db: db}
}

func (f *FieldScheduleRepository) FindAllWithPagination(ctx context.Context, params *dto.FieldScheduleRequestParam) ([]models.FieldSchedule, int64, error) {
	var (
		fields []models.FieldSchedule
		sort   string
		total  int64
	)

	if params.SortColumn != nil {
		sort = fmt.Sprintf("%s %s", *params.SortColumn, *params.SortOrder)
	} else {
		sort = "created_at desc"
	}

	limit := params.Limit
	offset := (params.Page - 1) * params.Limit
	err := f.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order(sort).
		Find(&fields).Error

	if err != nil {
		return nil, 0, errorWrap.WrapError(errConstants.ErrSqlQuery)
	}

	err = f.db.WithContext(ctx).
		Model(&fields).Count(&total).Error

	if err != nil {
		return nil, 0, errorWrap.WrapError(errConstants.ErrSqlQuery)
	}

	return fields, total, nil
}

func (f *FieldScheduleRepository) FindAllByFieldIDAndDate(ctx context.Context, fieldID int, date string) ([]models.FieldSchedule, error) {
	var fields []models.FieldSchedule

	err := f.db.
		WithContext(ctx).
		Preload("Field").
		Preload("Time").
		Where("field_id = ?", fieldID).
		Where("date = ?", date).
		Joins("LEFT JOIN times ON times.id = field_schedules.time_id").
		Order("times.start_time ASC").
		Find(&fields).
		Error

	if err != nil {
		return nil, errorWrap.WrapError(errConstants.ErrSqlQuery)
	}

	return fields, nil
}

func (f *FieldScheduleRepository) FindByUUID(ctx context.Context, uuid string) (*models.FieldSchedule, error) {
	var field models.FieldSchedule
	err := f.db.WithContext(ctx).
		Preload("Field").
		Preload("Time").
		Where("uuid = ?", uuid).
		First(&field).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorWrap.WrapError(errField.ErrFieldNotFound)
		}
		return nil, errorWrap.WrapError(errConstants.ErrSqlQuery)
	}
	return &field, nil
}

func (f *FieldScheduleRepository) FindByDateAndTimeID(ctx context.Context, date string, timeID int, fieldID int) (*models.FieldSchedule, error) {
	var field models.FieldSchedule

	err := f.db.WithContext(ctx).
		Where("date = ?", date).
		Where("time_id = ?", timeID).
		Where("field_id = ?", fieldID).
		First(&field).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errorWrap.WrapError(errConstants.ErrSqlQuery)
	}

	return &field, nil

}

func (f *FieldScheduleRepository) Create(ctx context.Context, request []models.FieldSchedule) error {
	err := f.db.WithContext(ctx).Create(&request).Error
	if err != nil {
		return errorWrap.WrapError(errConstants.ErrSqlQuery)
	}

	return nil
}

func (f *FieldScheduleRepository) Update(ctx context.Context, uuid string, request *models.FieldSchedule) (*models.FieldSchedule, error) {
	fieldSchedule, err := f.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	fieldSchedule.Date = request.Date
	err = f.db.WithContext(ctx).Save(&fieldSchedule).Error
	if err != nil {
		return nil, errorWrap.WrapError(errConstants.ErrSqlQuery)
	}

	return fieldSchedule, nil
}

func (f *FieldScheduleRepository) UpdateStatus(ctx context.Context, status constants.FieldScheduleStatus, uuid string) error {
	fieldSchedule, err := f.FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	fieldSchedule.Status = status
	err = f.db.WithContext(ctx).Save(&fieldSchedule).Error
	if err != nil {
		return errorWrap.WrapError(errConstants.ErrSqlQuery)
	}

	return nil
}

func (f *FieldScheduleRepository) Delete(ctx context.Context, uuid string) error {
	err := f.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&models.FieldSchedule{}).Error
	if err != nil {
		return errorWrap.WrapError(errConstants.ErrSqlQuery)
	}
	return nil
}
