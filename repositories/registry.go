package repositories

import (
	"gorm.io/gorm"

	fieldRepo "field-service/repositories/field"
	fieldScheduleRepo "field-service/repositories/fieldschedule"
	timeRepo "field-service/repositories/time"
)

type Registry struct {
	db *gorm.DB
}

type IRepostitoryRegistry interface {
	GetFieldRepository() fieldRepo.IFieldRepository
	GetFieldScheduleRepository() fieldScheduleRepo.IFieldScheduleRepository
	GetTimeRepository() timeRepo.ITimeRepository
}

func NewRepositoryRegistry(db *gorm.DB) IRepostitoryRegistry {
	return &Registry{
		db: db,
	}
}

func (r *Registry) GetFieldRepository() fieldRepo.IFieldRepository {
	return fieldRepo.NewFieldRepository(r.db)
}

func (r *Registry) GetFieldScheduleRepository() fieldScheduleRepo.IFieldScheduleRepository {
	return fieldScheduleRepo.NewFieldScheduleRepository(r.db)
}

func (r *Registry) GetTimeRepository() timeRepo.ITimeRepository {
	return timeRepo.NewTimeRepository(r.db)
}
