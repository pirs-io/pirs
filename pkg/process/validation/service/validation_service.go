package service

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"pirs.io/commons"
	"pirs.io/process/validation/models"
	"pirs.io/process/validation/validators"
	"reflect"
)

var (
	log = commons.GetLoggerFor("ValidationService")
)

type ValidationService struct {
	chainStart models.Validator
}

func NewValidationService(rawExtensions string, ignoreWrongExtension bool) *ValidationService {
	chainStart := buildValidationChains(rawExtensions, ignoreWrongExtension)
	return &ValidationService{
		chainStart: chainStart,
	}
}

func buildValidationChains(rawExtensions string, ignoreWrongExtension bool) models.Validator {
	// define validators
	requestValidator := &validators.ImportProcessRequestValidator{}
	fileTypeValidator := validators.NewFileTypeValidator(rawExtensions, ignoreWrongExtension)
	schemaValidator := &validators.SchemaValidator{}
	// create chain
	requestValidator.SetNext(fileTypeValidator)
	fileTypeValidator.SetNext(schemaValidator)

	return requestValidator
}

func (vs *ValidationService) ValidateProcessData(data *models.ImportProcessValidationData) bool {
	vs.chainStart.Validate(data)
	validationFlags := reflect.ValueOf(data.ValidationFlags)
	// all validations must pass
	for i := 0; i < validationFlags.NumField(); i++ {
		if validationFlags.Field(i).Bool() == false {
			log.Error().Msg(status.Errorf(codes.InvalidArgument, "process file %s is invalid: %s is false", data.ReqData.ProcessFileName, validationFlags.Field(i).String()).Error())
			return false
		}
	}
	return true
}

func (vs *ValidationService) ValidatePackageData(data *models.ImportPackageValidationData) {
	// todo
	panic("not implemented")
}
