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

// A ValidationService contains validation chain, that is used to validate models.ImportProcessValidationData.
type ValidationService struct {
	chainStart models.Validator
}

// NewValidationService creates instance of ValidationService with validation chain. This validation chain is created by
// function buildValidationChains.
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

// ValidateProcessData validates models.ImportProcessValidationData by models.Validator implementations. It returns true,
// if all the models.ValidationFlags are set to true. Otherwise, false is returned.
func (vs *ValidationService) ValidateProcessData(data *models.ImportProcessValidationData) bool {
	vs.chainStart.Validate(data)
	validationFlags := reflect.ValueOf(data.ValidationFlags)
	// all validations must pass
	for i := 0; i < validationFlags.NumField(); i++ {
		if validationFlags.Field(i).Bool() == false {
			log.Error().Msg(status.Errorf(codes.InvalidArgument, "process file %s is invalid: %s is false", data.ReqData.ProcessFileName, validationFlags.Type().Field(i).Name).Error())
			return false
		}
	}
	return true
}

func (vs *ValidationService) ValidatePackageData(data *models.ImportPackageValidationData) {
	// todo
	panic("not implemented")
}
