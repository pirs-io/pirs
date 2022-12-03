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

// A ValidationService contains validation chains, that are used to validate request sent from parent services, such as
// ImportService or DownloadService.
type ValidationService struct {
	chainStartImportProcess   models.Validator
	chainStartDownloadProcess models.Validator
}

// NewValidationService creates instance of ValidationService with validation chains.
func NewValidationService(rawExtensions string, ignoreWrongExtension bool) *ValidationService {
	chainStartImportProcess := buildValidationChainForImportProcess(rawExtensions, ignoreWrongExtension)
	chainStartDownloadProcess := buildValidationChainForDownloadProcess()

	return &ValidationService{
		chainStartImportProcess:   chainStartImportProcess,
		chainStartDownloadProcess: chainStartDownloadProcess,
	}
}

func buildValidationChainForImportProcess(rawExtensions string, ignoreWrongExtension bool) models.Validator {
	// define validators
	requestValidator := &validators.ImportProcessRequestValidator{}
	fileTypeValidator := validators.NewFileTypeValidator(rawExtensions, ignoreWrongExtension)
	schemaValidator := &validators.SchemaValidator{}
	// create chain
	requestValidator.SetNext(fileTypeValidator)
	fileTypeValidator.SetNext(schemaValidator)

	return requestValidator
}

func buildValidationChainForDownloadProcess() models.Validator {
	// define validators
	requestValidator := &validators.DownloadProcessRequestValidator{}

	return requestValidator
}

// ValidateProcessData validates models.ImportProcessValidationData by models.Validator implementations. It returns true,
// if all the models.ValidationFlags are set to true. Otherwise, false is returned.
func (vs *ValidationService) ValidateProcessData(data *models.ImportProcessValidationData) bool {
	vs.chainStartImportProcess.Validate(data)
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

// ValidateDownloadData todo
func (vs *ValidationService) ValidateDownloadData(data *models.DownloadProcessValidationData) bool {
	vs.chainStartDownloadProcess.Validate(data)
	validationFlags := reflect.ValueOf(data.ValidationFlags)
	// all validations must pass
	for i := 0; i < validationFlags.NumField(); i++ {
		if validationFlags.Field(i).Bool() == false {
			log.Error().Msg(status.Errorf(codes.InvalidArgument, "uri %s is invalid: %s is false", data.ReqData.Uri, validationFlags.Type().Field(i).Name).Error())
			return false
		}
	}
	return true
}
