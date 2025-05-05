package services

import (
	"context"
	"field-service/common/util"
	"field-service/constants"
	errFieldSchedule "field-service/constants/error/field_schedule"
	"field-service/domain/dto"
	"field-service/domain/models"
	"field-service/repositories"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type FieldScheduleService struct {
	repositories repositories.IRepostitoryRegistry
}

type IFieldScheduleService interface {
	GetAllWithPagination(context.Context, *dto.FieldScheduleRequestParam) (*util.PaginationResult, error)
	GetAllByFieldIDAndDate(context.Context, string, string) ([]dto.FieldScheduleForBookingResponse, error)
	GetByUUID(context.Context, string) (*dto.FieldScheduleReponse, error)
	GenerateScheduleForOneMonth(context.Context, *dto.GenerateFieldScheduleForOneMonthRequest) error
	Create(context.Context, *dto.FieldScheduleRequest) error
	Update(context.Context, string, *dto.UpdateFieldScheduleRequest) (*dto.FieldScheduleReponse, error)
	UpdateStatus(context.Context, *dto.UpdateStatusFieldScheduleRequest) error
	Delete(context.Context, string) error
}

func NewFieldScheduleService(repositories repositories.IRepostitoryRegistry) IFieldScheduleService {
	return &FieldScheduleService{repositories: repositories}
}

func (f *FieldScheduleService) GetAllWithPagination(ctx context.Context, param *dto.FieldScheduleRequestParam) (*util.PaginationResult, error) {
	fieldSchedules, total, err := f.repositories.GetFieldScheduleRepository().FindAllWithPagination(ctx, param)
	if err != nil {
		return nil, err
	}

	fieldSchedulesResults := make([]dto.FieldScheduleReponse, 0, len(fieldSchedules))
	for _, schedule := range fieldSchedules {
		fieldSchedulesResults = append(fieldSchedulesResults, dto.FieldScheduleReponse{
			UUID:         schedule.UUID,
			FieldName:    schedule.Field.Name,
			Date:         schedule.Date.Format("2006-01-02"),
			PricePerHour: schedule.Field.PricePerHour,
			Status:       schedule.Status.GetStatusString(),
			Time:         fmt.Sprintf("%s - %s", schedule.Time.StartTime, schedule.Time.EndTime),
			CreatedAt:    schedule.CreatedAt,
			UpdatedAt:    schedule.UpdatedAt,
		})
	}

	pagination := &util.PaginationParam{
		Page:  param.Page,
		Limit: param.Limit,
		Count: total,
		Data:  fieldSchedulesResults,
	}

	response := util.GeneratePagination(*pagination)

	return &response, nil
}

func (f *FieldScheduleService) convertMonthName(inputDate string) string {
	date, err := time.Parse(time.DateOnly, inputDate)
	if err != nil {
		return ""
	}

	indonesiaMonth := map[string]string{
		"Jan": "Jan",
		"Feb": "Feb",
		"Mar": "Mar",
		"Apr": "Apr",
		"May": "Mei",
		"Jun": "Jun",
		"Jul": "Jul",
		"Aug": "Agu",
		"Sep": "Sep",
		"Oct": "Okt",
		"Nov": "Nov",
		"Dec": "Des",
	}

	formattedDate := date.Format("02 Jan")
	day := formattedDate[0:3]
	month := formattedDate[3:]
	formattedDate = fmt.Sprintf("%s %s", day, indonesiaMonth[month])
	return formattedDate
}

func (f *FieldScheduleService) GetAllByFieldIDAndDate(ctx context.Context, uuid, date string) ([]dto.FieldScheduleForBookingResponse, error) {
	field, err := f.repositories.GetFieldRepository().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	fieldSchedules, err := f.repositories.GetFieldScheduleRepository().FindAllByFieldIDAndDate(ctx, int(field.ID), date)
	if err != nil {
		return nil, err
	}

	fieldSchedulesResults := make([]dto.FieldScheduleForBookingResponse, 0, len(fieldSchedules))
	for _, schedule := range fieldSchedules {
		pricePerHour := float64(schedule.Field.PricePerHour)
		startTime, _ := time.Parse("15:04:05", schedule.Time.StartTime)
		endTime, _ := time.Parse("15:04:05", schedule.Time.EndTime)
		fieldSchedulesResults = append(fieldSchedulesResults, dto.FieldScheduleForBookingResponse{
			UUID:         schedule.UUID,
			PricePerHour: util.FormatRupiah(&pricePerHour),
			Date:         f.convertMonthName(schedule.Date.Format(time.DateOnly)),
			Status:       schedule.Status.GetStatusString(),
			Time:         fmt.Sprintf("%s - %s", startTime.Format("15:04"), endTime.Format("15:04")),
		})
	}

	return fieldSchedulesResults, nil
}

func (f *FieldScheduleService) GetByUUID(ctx context.Context, uuid string) (*dto.FieldScheduleReponse, error) {
	fieldSchedule, err := f.repositories.GetFieldScheduleRepository().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	response := dto.FieldScheduleReponse{
		UUID:         fieldSchedule.UUID,
		FieldName:    fieldSchedule.Field.Name,
		PricePerHour: fieldSchedule.Field.PricePerHour,
		Status:       fieldSchedule.Status.GetStatusString(),
		Time:         fmt.Sprintf("%s - %s", fieldSchedule.Time.StartTime, fieldSchedule.Time.EndTime),
		CreatedAt:    fieldSchedule.CreatedAt,
		UpdatedAt:    fieldSchedule.UpdatedAt,
	}

	return &response, nil
}

func (f *FieldScheduleService) Create(ctx context.Context, request *dto.FieldScheduleRequest) error {
	// cek field schedule berdasarkan uuid ada atau tidak
	field, err := f.repositories.GetFieldRepository().FindByUUID(ctx, request.FieldID)
	if err != nil {
		return err
	}
	// jika tidak, looping time untuk membuat fieldSchedule sebanyak time request.
	fieldSchedules := make([]models.FieldSchedule, 0, len(request.TimeIDs))
	dateParsed, _ := time.Parse(time.DateOnly, request.Date)
	for _, timeID := range request.TimeIDs {
		// a.check apakah waktu ada di time atau tidak, jika tidak return error
		scheduleTime, err := f.repositories.GetTimeRepository().FindByUUID(ctx, timeID)
		if err != nil {
			return err
		}

		// b. jika ada, cek apakah sudah ada fieldschedule yang memiliki waktu yang sedang diiterasi,
		schedule, err := f.repositories.GetFieldScheduleRepository().FindByDateAndTimeID(ctx, request.Date, int(scheduleTime.ID), int(field.ID))
		if err != nil {
			return err
		}

		if schedule != nil {
			// jika ada, return error field schedule sudah ada
			return errFieldSchedule.ErrFieldShceduleExist
		}

		// jika belum masukan kearray models fieldscheduels
		fieldSchedules = append(fieldSchedules, models.FieldSchedule{
			UUID:    uuid.New(),
			FieldID: field.ID,
			TimeID:  scheduleTime.ID,
			Date:    dateParsed,
			Status:  constants.Available,
		})
	}

	// masukan fieldSchedules yang ada diarray ke repository
	err = f.repositories.GetFieldScheduleRepository().Create(ctx, fieldSchedules)
	if err != nil {
		return nil
	}
	// return hasil
	return nil

}

func (f *FieldScheduleService) GenerateScheduleForOneMonth(ctx context.Context, request *dto.GenerateFieldScheduleForOneMonthRequest) error {
	// cek apakah field ada atau tidak
	field, err := f.repositories.GetFieldRepository().FindByUUID(ctx, request.FieldID)
	if err != nil {
		return nil
	}

	// ambil semua time
	times, err := f.repositories.GetTimeRepository().FindAll(ctx)
	if err != nil {
		return nil
	}

	// buat semua schedule untuk 1 bulan atau 30 hari kedepan
	// a. inisiasi numberOfDays
	numberOfDays := 30
	// b. buat array kosong fieldSchedules
	fieldSchedules := make([]models.FieldSchedule, 0, numberOfDays)
	now := time.Now().Add(time.Duration(1) * 24 * time.Hour)
	// looping
	for i := 0; i < numberOfDays; i++ {
		currentDate := now.AddDate(0, 0, i)
		// sama seperti membuat field schedule
		for _, item := range times {
			schedule, err := f.repositories.GetFieldScheduleRepository().FindByDateAndTimeID(
				ctx,
				currentDate.Format(time.DateOnly),
				int(item.ID),
				int(field.ID),
			)

			if err != nil {
				return err
			}

			if schedule != nil {
				return errFieldSchedule.ErrFieldShceduleExist
			}

			fieldSchedules = append(fieldSchedules, models.FieldSchedule{
				UUID:    uuid.New(),
				FieldID: field.ID,
				TimeID:  item.ID,
				Date:    currentDate,
				Status:  constants.Available,
			})

		}
	}
	// masukkan semua fieldSchedules ke repository create field schedule
	err = f.repositories.GetFieldScheduleRepository().Create(ctx, fieldSchedules)
	if err != nil {
		return err
	}
	// return hasil
	return nil
}

func (f *FieldScheduleService) Update(ctx context.Context, uuid string, request *dto.UpdateFieldScheduleRequest) (*dto.FieldScheduleReponse, error) {
	// cek apakah field schedule ada
	fieldSchedule, err := f.repositories.GetFieldScheduleRepository().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	// cek apakah time ada
	scheduleTime, err := f.repositories.GetTimeRepository().FindByUUID(ctx, request.TimeID)
	if err != nil {
		return nil, err
	}

	// cek apakah field schedule yang memiliki waktu dan tanggal sesuai request ada
	isTimeExist, err := f.repositories.GetFieldScheduleRepository().FindByDateAndTimeID(
		ctx,
		request.Date,
		int(scheduleTime.ID),
		int(fieldSchedule.FieldID),
	)
	if err != nil {
		return nil, err
	}

	// cek apakah field schedule dengan waktu dan tanggal tertentu ada, dan request date tidak sama dengan tanggal fieldSchedule
	if isTimeExist != nil && request.Date != fieldSchedule.Date.Format(time.DateOnly) {
		// jika iya return error fieldschedule sudah ada
		return nil, errFieldSchedule.ErrFieldShceduleExist
	}

	// parsing date request
	dateParsed, _ := time.Parse(time.DateOnly, request.Date)
	// masukan date request, models field schedule ke repository field schedule update
	fieldResult, err := f.repositories.GetFieldScheduleRepository().Update(ctx, uuid, &models.FieldSchedule{
		Date:   dateParsed,
		TimeID: scheduleTime.ID,
	})
	if err != nil {
		return nil, err
	}

	response := dto.FieldScheduleReponse{
		UUID:         fieldResult.UUID,
		FieldName:    fieldResult.Field.Name,
		Date:         fieldResult.Date.Format(time.DateOnly),
		PricePerHour: fieldResult.Field.PricePerHour,
		Status:       fieldResult.Status.GetStatusString(),
		Time:         fmt.Sprintf("%s - %s", scheduleTime.StartTime, scheduleTime.EndTime),
		CreatedAt:    fieldResult.CreatedAt,
		UpdatedAt:    fieldResult.UpdatedAt,
	}

	// return hasil update
	return &response, nil

}

func (f *FieldScheduleService) UpdateStatus(ctx context.Context, request *dto.UpdateStatusFieldScheduleRequest) error {
	for _, item := range request.FieldScheduleIDs {
		_, err := f.repositories.GetFieldScheduleRepository().FindByUUID(ctx, item)
		if err != nil {
			return err
		}

		err = f.repositories.GetFieldScheduleRepository().UpdateStatus(ctx, constants.Booked, item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *FieldScheduleService) Delete(ctx context.Context, uuid string) error {
	_, err := f.repositories.GetFieldScheduleRepository().FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = f.repositories.GetFieldScheduleRepository().Delete(ctx, uuid)
	if err != nil {
		return err
	}

	return nil
}
