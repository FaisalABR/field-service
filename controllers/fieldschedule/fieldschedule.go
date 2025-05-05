package controllers

import (
	errValidation "field-service/common/error"
	"field-service/common/response"
	"field-service/domain/dto"
	"field-service/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type FieldScheduleController struct {
	service services.IServiceRegistry
}

type IFieldScheduleController interface {
	GetAllWithPagination(*gin.Context)
	GetAllByFieldIDAndDate(*gin.Context)
	GetByUUID(*gin.Context)
	GenerateScheduleForOneMonth(*gin.Context)
	Create(*gin.Context)
	Update(*gin.Context)
	UpdateStatus(*gin.Context)
	Delete(*gin.Context)
}

func NewFieldScheduleController(service services.IServiceRegistry) IFieldScheduleController {
	return &FieldScheduleController{
		service: service,
	}
}

func (f *FieldScheduleController) GetAllWithPagination(c *gin.Context) {
	var params dto.FieldScheduleRequestParam
	err := c.ShouldBindQuery(&params)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code:  http.StatusBadRequest,
			Error: err,
			Gin:   c,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(params)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPResp{
			Code:    http.StatusBadRequest,
			Message: &errMessage,
			Data:    errResponse,
			Error:   err,
			Gin:     c,
		})
		return
	}

	result, err := f.service.GetFieldSchedule().GetAllWithPagination(c, &params)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code:  http.StatusInternalServerError,
			Error: err,
			Gin:   c,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}

func (f *FieldScheduleController) GetAllByFieldIDAndDate(c *gin.Context) {
	var params dto.FieldScheduleByFieldIDAndDateRequestParam
	err := c.ShouldBindQuery(params)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code:  http.StatusBadRequest,
			Error: err,
			Gin:   c,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(params)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPResp{
			Code:    http.StatusBadRequest,
			Message: &errMessage,
			Data:    errResponse,
			Error:   err,
			Gin:     c,
		})
		return
	}

	result, err := f.service.GetFieldSchedule().GetAllByFieldIDAndDate(c, c.Param("uuid"), params.Date)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code:  http.StatusInternalServerError,
			Error: err,
			Gin:   c,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}

func (f *FieldScheduleController) GetByUUID(c *gin.Context) {
	result, err := f.service.GetFieldSchedule().GetByUUID(c, c.Param("uuid"))
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code:  http.StatusInternalServerError,
			Error: err,
			Gin:   c,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}

func (f *FieldScheduleController) GenerateScheduleForOneMonth(c *gin.Context) {
	var request dto.GenerateFieldScheduleForOneMonthRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code:  http.StatusBadRequest,
			Error: err,
			Gin:   c,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPResp{
			Code:    http.StatusBadRequest,
			Message: &errMessage,
			Data:    errResponse,
			Error:   err,
			Gin:     c,
		})
		return
	}

	err = f.service.GetFieldSchedule().GenerateScheduleForOneMonth(c, &request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code:  http.StatusInternalServerError,
			Error: err,
			Gin:   c,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Gin:  c,
	})

}

func (f *FieldScheduleController) Create(c *gin.Context) {
	var request dto.FieldScheduleRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code:  http.StatusBadRequest,
			Error: err,
			Gin:   c,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPResp{
			Code:    http.StatusBadRequest,
			Message: &errMessage,
			Data:    errResponse,
			Error:   err,
			Gin:     c,
		})
		return
	}

	err = f.service.GetFieldSchedule().Create(c, &request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code:  http.StatusInternalServerError,
			Error: err,
			Gin:   c,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusCreated,
		Gin:  c,
	})
}

func (f *FieldScheduleController) Update(c *gin.Context) {
	var request dto.UpdateFieldScheduleRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code:  http.StatusBadRequest,
			Error: err,
			Gin:   c,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPResp{
			Code:    http.StatusBadRequest,
			Message: &errMessage,
			Data:    errResponse,
			Error:   err,
			Gin:     c,
		})
		return
	}

	result, err := f.service.GetFieldSchedule().Update(c, c.Param("uuid"), &request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code:  http.StatusInternalServerError,
			Error: err,
			Gin:   c,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})

}

func (f *FieldScheduleController) UpdateStatus(c *gin.Context) {
	var request dto.UpdateStatusFieldScheduleRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code:  http.StatusBadRequest,
			Error: err,
			Gin:   c,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errValidation.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHTTPResp{
			Code:    http.StatusBadRequest,
			Message: &errMessage,
			Data:    errResponse,
			Error:   err,
			Gin:     c,
		})
		return
	}

	err = f.service.GetFieldSchedule().UpdateStatus(c, &request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code:  http.StatusInternalServerError,
			Error: err,
			Gin:   c,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Gin:  c,
	})
}

func (f *FieldScheduleController) Delete(c *gin.Context) {
	err := f.service.GetFieldSchedule().Delete(c, c.Param("uuid"))
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code:  http.StatusInternalServerError,
			Error: err,
			Gin:   c,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Gin:  c,
	})

}
