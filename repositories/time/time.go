package repositories

import (
	"context"
	"errors"
	errorWrap "field-service/common/error"
	errConstants "field-service/constants/error"
	errTime "field-service/constants/error/time"
	"field-service/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TimeRepository struct {
	db *gorm.DB
}

type ITimeRepository interface {
	FindAll(context.Context) ([]models.Time, error)
	FindByUUID(context.Context, string) (*models.Time, error)
	FindByID(context.Context, int) (*models.Time, error)
	Create(context.Context, *models.Time) (*models.Time, error)
}

func NewTimeRepository(db *gorm.DB) ITimeRepository {
	return &TimeRepository{db: db}
}

func (t *TimeRepository) FindAll(ctx context.Context) ([]models.Time, error) {
	var times []models.Time
	err := t.db.WithContext(ctx).Find(&times).Error
	if err != nil {
		return nil, errorWrap.WrapError(errConstants.ErrSqlQuery)
	}

	return times, nil
}

func (t *TimeRepository) FindByUUID(ctx context.Context, uuid string) (*models.Time, error) {
	var time models.Time
	err := t.db.WithContext(ctx).Where("uuid = ?", uuid).First(&time).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorWrap.WrapError(errTime.ErrTimeNotFound)
		}

		return nil, errorWrap.WrapError(errConstants.ErrSqlQuery)
	}

	return &time, nil
}

func (t *TimeRepository) FindByID(ctx context.Context, id int) (*models.Time, error) {
	var time models.Time

	err := t.db.WithContext(ctx).Where("id = ?", id).First(&time).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorWrap.WrapError(errTime.ErrTimeNotFound)
		}

		return nil, errorWrap.WrapError(errConstants.ErrSqlQuery)
	}

	return &time, nil
}

func (t *TimeRepository) Create(ctx context.Context, time *models.Time) (*models.Time, error) {
	time.UUID = uuid.New()
	err := t.db.WithContext(ctx).Create(time).Error
	if err != nil {
		return nil, errorWrap.WrapError(errConstants.ErrSqlQuery)
	}
	return time, nil
}
