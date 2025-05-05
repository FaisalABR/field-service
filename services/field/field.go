package services

import (
	"bytes"
	"context"
	gcs "field-service/common/gcs"
	"field-service/common/util"
	errConstant "field-service/constants/error"
	"field-service/domain/dto"
	"field-service/domain/models"
	"field-service/repositories"
	"fmt"
	"io"
	"mime/multipart"
	"path"
	"time"

	"github.com/sirupsen/logrus"
)

type FieldService struct {
	repositories repositories.IRepostitoryRegistry
	gcs          gcs.IGCSClient
}

type IFieldService interface {
	GetAllWithPagination(context.Context, *dto.FieldRequestParam) (*util.PaginationResult, error)
	GetAllWithoutPagination(context.Context) ([]dto.FieldResponse, error)
	GetByUUID(context.Context, string) (*dto.FieldResponse, error)
	Create(context.Context, *dto.FieldRequest) (*dto.FieldResponse, error)
	Update(context.Context, string, *dto.UpdateFieldRequest) (*dto.FieldResponse, error)
	Delete(context.Context, string) error
}

func NewFieldService(
	repositories repositories.IRepostitoryRegistry,
	gcs gcs.IGCSClient,
) IFieldService {
	return &FieldService{
		repositories: repositories,
		gcs:          gcs,
	}
}

func (f *FieldService) GetAllWithPagination(ctx context.Context, param *dto.FieldRequestParam) (*util.PaginationResult, error) {
	fields, total, err := f.repositories.GetFieldRepository().FindAllWithPagination(ctx, param)
	if err != nil {
		return nil, err
	}

	fieldResults := make([]dto.FieldResponse, 0, len(fields))
	for _, field := range fields {
		fieldResults = append(fieldResults, dto.FieldResponse{
			UUID:         field.UUID,
			Code:         field.Code,
			Name:         field.Name,
			PricePerHour: field.PricePerHour,
			Images:       field.Images,
			CreatedAt:    field.CreatedAt,
			UpdatedAt:    field.UpdatedAt,
		})
	}

	pagination := &util.PaginationParam{
		Page:  param.Page,
		Count: total,
		Limit: param.Limit,
		Data:  fieldResults,
	}

	responses := util.GeneratePagination(*pagination)

	return &responses, nil
}

func (f *FieldService) GetAllWithoutPagination(ctx context.Context) ([]dto.FieldResponse, error) {
	fields, err := f.repositories.GetFieldRepository().FindAllWithoutPagination(ctx)
	if err != nil {
		return nil, err
	}

	fieldResults := make([]dto.FieldResponse, 0, len(fields))
	for _, field := range fields {
		fieldResults = append(fieldResults, dto.FieldResponse{
			UUID:         field.UUID,
			Code:         field.Code,
			Name:         field.Name,
			PricePerHour: field.PricePerHour,
			Images:       field.Images,
			CreatedAt:    field.CreatedAt,
			UpdatedAt:    field.UpdatedAt,
		})
	}

	return fieldResults, nil
}

func (f *FieldService) GetByUUID(ctx context.Context, uuid string) (*dto.FieldResponse, error) {
	field, err := f.repositories.GetFieldRepository().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	pricePerHour := float64(field.PricePerHour)
	fieldResult := dto.FieldResponse{
		UUID:         field.UUID,
		Code:         field.Code,
		Name:         field.Name,
		PricePerHour: util.FormatRupiah(&pricePerHour),
		Images:       field.Images,
		CreatedAt:    field.CreatedAt,
		UpdatedAt:    field.UpdatedAt,
	}

	return &fieldResult, nil
}

func (f *FieldService) validateUpload(images []multipart.FileHeader) error {
	if images == nil || len(images) == 0 {
		return errConstant.ErrInvalidUploadFile
	}

	for _, image := range images {
		if image.Size > 5*1024*1024 {
			return errConstant.ErrSizeTooBig
		}
	}

	return nil
}

func (f *FieldService) processAndUploadImage(ctx context.Context, image multipart.FileHeader) (string, error) {
	file, err := image.Open()
	if err != nil {
		return "", err
	}

	defer file.Close()

	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, file)
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("images/%s-%s-%s", time.Now().Format("20060102150405"), image.Filename, path.Ext(image.Filename))
	url, err := f.gcs.UploadFile(ctx, filename, buffer.Bytes())
	if err != nil {
		return "", err
	}

	return url, nil
}

func (f *FieldService) uploadImage(ctx context.Context, images []multipart.FileHeader) ([]string, error) {
	err := f.validateUpload(images)
	if err != nil {
		return nil, err
	}

	urls := make([]string, 0, len(images))
	for _, image := range images {
		url, err := f.processAndUploadImage(ctx, image)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	return urls, nil
}

func (f *FieldService) Create(ctx context.Context, req *dto.FieldRequest) (*dto.FieldResponse, error) {
	imageUrl, err := f.uploadImage(ctx, req.Images)
	if err != nil {
		logrus.Errorf("Fieldservice Create - 1 %v", err)
		return nil, err
	}

	field, err := f.repositories.GetFieldRepository().Create(ctx, &models.Field{
		Name:         req.Name,
		Code:         req.Code,
		PricePerHour: req.PricePerHour,
		Images:       imageUrl,
	})
	if err != nil {
		logrus.Errorf("Fieldservice Create - 2 %v", err)
		return nil, err
	}

	response := dto.FieldResponse{
		UUID:         field.UUID,
		Code:         field.Code,
		Name:         field.Name,
		PricePerHour: field.PricePerHour,
		Images:       field.Images,
		CreatedAt:    field.CreatedAt,
		UpdatedAt:    field.UpdatedAt,
	}

	return &response, nil
}

func (f *FieldService) Update(ctx context.Context, uuid string, req *dto.UpdateFieldRequest) (*dto.FieldResponse, error) {
	field, err := f.repositories.GetFieldRepository().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if req.Images != nil && len(req.Images) > 0 {
		imageUrl, err := f.uploadImage(ctx, req.Images)
		if err != nil {
			return nil, err
		}
		field.Images = imageUrl
	}

	fieldResult, err := f.repositories.GetFieldRepository().Update(ctx, uuid, &models.Field{
		Name:         req.Name,
		Code:         req.Code,
		PricePerHour: req.PricePerHour,
		Images:       field.Images,
	})
	if err != nil {
		return nil, err
	}

	response := dto.FieldResponse{
		UUID:         fieldResult.UUID,
		Code:         fieldResult.Code,
		Name:         fieldResult.Name,
		PricePerHour: fieldResult.PricePerHour,
		Images:       fieldResult.Images,
		CreatedAt:    fieldResult.CreatedAt,
		UpdatedAt:    fieldResult.UpdatedAt,
	}

	return &response, nil
}

func (f *FieldService) Delete(ctx context.Context, uuid string) error {
	_, err := f.repositories.GetFieldRepository().FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = f.repositories.GetFieldRepository().Delete(ctx, uuid)
	if err != nil {
		return err
	}

	return nil
}
