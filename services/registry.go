package services

import (
	gcs "field-service/common/gcs"
	"field-service/repositories"
	fieldService "field-service/services/field"
	fieldScheduleService "field-service/services/fieldschedule"
	timeService "field-service/services/time"
)

type ServiceRegistry struct {
	repositories repositories.IRepostitoryRegistry
	gcs          gcs.IGCSClient
}

type IServiceRegistry interface {
	GetField() fieldService.IFieldService
	GetFieldSchedule() fieldScheduleService.IFieldScheduleService
	GetTime() timeService.ITimeService
}

func NewServiceRegistry(repositories repositories.IRepostitoryRegistry, gcs gcs.IGCSClient) IServiceRegistry {
	return &ServiceRegistry{
		repositories: repositories,
		gcs:          gcs,
	}
}

func (s *ServiceRegistry) GetField() fieldService.IFieldService {
	return fieldService.NewFieldService(s.repositories, s.gcs)
}

func (s *ServiceRegistry) GetFieldSchedule() fieldScheduleService.IFieldScheduleService {
	return fieldScheduleService.NewFieldScheduleService(s.repositories)
}

func (s *ServiceRegistry) GetTime() timeService.ITimeService {
	return timeService.NewTimeService(s.repositories)
}
