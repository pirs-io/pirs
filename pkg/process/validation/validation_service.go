package validation

import (
	"pirs.io/process/validation/models"
	"pirs.io/process/validation/validators"
	"reflect"
)

type ValidationService struct {
	chainStart models.Validator
}

func NewValidationService(rawExtensions string) *ValidationService {
	chainStart := buildValidationChains(rawExtensions)
	return &ValidationService{chainStart: chainStart}
}

func buildValidationChains(rawExtensions string) models.Validator {
	// define validators
	requestValidator := &validators.ImportProcessRequestValidator{}
	extensionValidator := validators.NewExtensionValidator(rawExtensions)
	// create chain
	requestValidator.SetNext(extensionValidator)

	return requestValidator
}

func (vs *ValidationService) ValidateProcessData(data *models.ImportProcessValidationData) (bool, error) {
	vs.chainStart.Validate(data)
	validationFlags := reflect.ValueOf(data.ValidationFlags)
	// all validations must pass
	for i := 0; i < validationFlags.NumField(); i++ {
		if validationFlags.Field(i).Bool() == false {
			// todo use error as description msg
			return false, nil
		}
	}
	return true, nil
}

func (vs *ValidationService) ValidatePackageData(data *models.ImportPackageValidationData) {
	panic("not implemented")
}
