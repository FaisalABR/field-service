package dto

import (
	"field-service/constants"
	"time"

	"github.com/google/uuid"
)

// Field schedyle request
type FieldScheduleRequest struct {
	FieldID string   `json:"fieldID" validate:"required"`
	Date    string   `json:"date" validate:"required"`
	TimeIDs []string `json:"timeIDs" validate:"required"`
}

// Generate field schedule for one month request
type GenerateFieldScheduleForOneMonthRequest struct {
	FieldID string `json:"fieldID" validate:"required"`
}

// Update field schedule request
type UpdateFieldScheduleRequest struct {
	Date   string `json:"date" validate:"required"`
	TimeID string `json:"timeID" validate:"required"`
}

// Update status field schedule request
type UpdateStatusFieldScheduleRequest struct {
	FieldScheduleIDs []string `json:"fieldScheduleIDs" validate:"required"`
}

// Field schedule response
type FieldScheduleReponse struct {
	UUID         uuid.UUID                         `json:"uuid"`
	FieldName    string                            `json:"fieldName"`
	Date         string                            `json:"code"`
	PricePerHour int                               `json:"pricePerHour"`
	Status       constants.FieldScheduleStatusName `json:"status"`
	Time         string                            `json:"time"`
	CreatedAt    *time.Time                        `json:"createdAt"`
	UpdatedAt    *time.Time                        `json:"updatedAt"`
}

// field schedule for booking response

type FieldScheduleForBookingResponse struct {
	UUID         uuid.UUID                         `json:"uuid"`
	PricePerHour string                            `json:"pricePerHour"`
	Date         string                            `json:"date"`
	Status       constants.FieldScheduleStatusName `json:"status"`
	Time         string                            `json:"time"`
}

// field schedule request params
type FieldScheduleRequestParam struct {
	Page       int     `form:"page" validate:"required"`
	Limit      int     `form:"limit" validate:"required"`
	SortColumn *string `form:"sortColumn"`
	SortOrder  *string `form:"sortOrder"`
}

type FieldScheduleByFieldIDAndDateRequestParam struct {
	Date string `form:"date" validate:"required"`
}
