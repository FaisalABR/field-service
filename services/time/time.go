package services

import (
	"context"
	"field-service/domain/dto"
	"field-service/domain/models"
	"field-service/repositories"

	"github.com/google/uuid"
)

type TimeService struct {
	repositories repositories.IRepostitoryRegistry
}

type ITimeService interface {
	GetAll(context.Context) ([]dto.TimeResponse, error)
	GetByUUID(context.Context, string) (*dto.TimeResponse, error)
	Create(context.Context, *dto.TimeRequest) (*dto.TimeResponse, error)
}

func NewTimeService(repositories repositories.IRepostitoryRegistry) ITimeService {
	return &TimeService{
		repositories: repositories,
	}
}

func (t *TimeService) GetAll(ctx context.Context) ([]dto.TimeResponse, error) {
	times, err := t.repositories.GetTimeRepository().FindAll(ctx)
	if err != nil {
		return nil, err
	}

	timeResult := make([]dto.TimeResponse, 0, len(times))
	for _, time := range times {
		timeResult = append(timeResult, dto.TimeResponse{
			UUID:      time.UUID,
			StartTime: time.StartTime,
			EndTime:   time.EndTime,
			CreatedAt: time.CreatedAt,
			UpdatedAt: time.UpdatedAt,
		})
	}

	return timeResult, nil
}

func (t *TimeService) GetByUUID(ctx context.Context, uuid string) (*dto.TimeResponse, error) {
	time, err := t.repositories.GetTimeRepository().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	response := dto.TimeResponse{
		UUID:      time.UUID,
		StartTime: time.StartTime,
		EndTime:   time.EndTime,
		CreatedAt: time.CreatedAt,
		UpdatedAt: time.UpdatedAt,
	}

	return &response, nil
}

func (t *TimeService) Create(ctx context.Context, request *dto.TimeRequest) (*dto.TimeResponse, error) {
	time, err := t.repositories.GetTimeRepository().Create(ctx, &models.Time{
		UUID:      uuid.New(),
		StartTime: request.StartTime,
		EndTime:   request.Endtime,
	})

	if err != nil {
		return nil, err
	}

	response := dto.TimeResponse{
		UUID:      time.UUID,
		StartTime: time.StartTime,
		EndTime:   time.EndTime,
		CreatedAt: time.CreatedAt,
		UpdatedAt: time.UpdatedAt,
	}

	return &response, nil
}
