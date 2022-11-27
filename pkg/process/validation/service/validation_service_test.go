package service

import (
	"github.com/stretchr/testify/assert"
	"pirs.io/process/mocks"
	"pirs.io/process/service/models"
	valModels "pirs.io/process/validation/models"
	"testing"
)

func TestValidationService_ValidateProcessData(t *testing.T) {
	valData := &valModels.ImportProcessValidationData{
		ReqData:         models.ImportProcessRequestData{},
		ValidationFlags: valModels.ValidationFlags{},
	}

	// is valid
	chainSucess := buildValidationChainsTest(false)
	vsSucess := &ValidationService{
		chainStart: chainSucess,
	}
	result := vsSucess.ValidateProcessData(valData)
	assert.Equal(t, true, result)

	// is invalid
	chainFail := buildValidationChainsTest(true)
	vsFail := &ValidationService{
		chainStart: chainFail,
	}
	result = vsFail.ValidateProcessData(valData)
	assert.Equal(t, false, result)
}

func buildValidationChainsTest(isFail bool) valModels.Validator {
	// define validators
	requestValidator := &mocks.ImportProcessRequestValidator{
		MockResult: !isFail,
	}
	fileTypeValidator := &mocks.FileTypeValidator{
		MockResult: true,
	}
	schemaValidator := &mocks.SchemaValidator{
		MockResult: true,
	}
	// create chain
	requestValidator.SetNext(fileTypeValidator)
	fileTypeValidator.SetNext(schemaValidator)

	return requestValidator
}
