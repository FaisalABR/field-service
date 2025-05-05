package repositories

import (
	"context"
	"errors"
	errorWrap "field-service/common/error"
	errConstants "field-service/constants/error"
	errField "field-service/constants/error/field"
	"field-service/domain/dto"
	"field-service/domain/models"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FieldRepository struct {
	db *gorm.DB
}

type IFieldRepository interface {
	FindAllWithPagination(context.Context, *dto.FieldRequestParam) ([]models.Field, int64, error)
	FindAllWithoutPagination(context.Context) ([]models.Field, error)
	FindByUUID(context.Context, string) (*models.Field, error)
	Create(context.Context, *models.Field) (*models.Field, error)
	Update(context.Context, string, *models.Field) (*models.Field, error)
	Delete(context.Context, string) error
}

func NewFieldRepository(db *gorm.DB) IFieldRepository {
	return &FieldRepository{db: db}
}

func (f *FieldRepository) FindAllWithPagination(ctx context.Context, params *dto.FieldRequestParam) ([]models.Field, int64, error) {
	var (
		fields []models.Field
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

func (f *FieldRepository) FindAllWithoutPagination(ctx context.Context) ([]models.Field, error) {
	var fields []models.Field

	err := f.db.
		WithContext(ctx).
		Find(&fields).
		Error

	if err != nil {
		return nil, errorWrap.WrapError(errConstants.ErrSqlQuery)
	}

	return fields, nil
}

func (f *FieldRepository) FindByUUID(ctx context.Context, uuid string) (*models.Field, error) {
	var field models.Field

	err := f.db.WithContext(ctx).
		Where("uuid = ?", uuid).First(&field).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorWrap.WrapError(errField.ErrFieldNotFound)
		}

		return nil, errorWrap.WrapError(errConstants.ErrSqlQuery)
	}

	return &field, nil
}

func (f *FieldRepository) Create(ctx context.Context, request *models.Field) (*models.Field, error) {
	field := models.Field{
		UUID:         uuid.New(),
		Code:         request.Code,
		Name:         request.Name,
		Images:       request.Images,
		PricePerHour: request.PricePerHour,
	}

	err := f.db.WithContext(ctx).Create(&field).Error
	if err != nil {
		return nil, errorWrap.WrapError(errConstants.ErrSqlQuery)
	}

	return &field, nil
}

func (f *FieldRepository) Update(ctx context.Context, uuid string, request *models.Field) (*models.Field, error) {
	field := models.Field{
		Code:         request.Code,
		Name:         request.Name,
		Images:       request.Images,
		PricePerHour: request.PricePerHour,
	}

	err := f.db.WithContext(ctx).Where("uuid = ?", uuid).Updates(&field).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorWrap.WrapError(errField.ErrFieldNotFound)
		}

		return nil, errorWrap.WrapError(errConstants.ErrSqlQuery)
	}

	return &field, nil
}

func (f *FieldRepository) Delete(ctx context.Context, uuid string) error {
	err := f.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&models.Field{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorWrap.WrapError(errField.ErrFieldNotFound)
		}

		return errorWrap.WrapError(errConstants.ErrSqlQuery)
	}

	return nil
}
