package service

import (
	"github.com/stretchr/testify/assert"
	mockTesting "pirs.io/process/mocks/testing"
	"pirs.io/process/service/models"
	valModels "pirs.io/process/validation/models"
	"testing"
)

func TestValidationService_ValidateProcessData(t *testing.T) {
	valData := &valModels.ImportProcessValidationData{
		ReqData:         models.ImportProcessRequestData{},
		ValidationFlags: valModels.ImportProcessValidationFlags{},
	}

	// is valid
	chainSucess := buildValidationChainForImportProcessTest(false)
	vsSucess := &ValidationService{
		chainStartImportProcess: chainSucess,
	}
	result := vsSucess.ValidateProcessData(valData)
	assert.Equal(t, true, result)

	// is invalid
	chainFail := buildValidationChainForImportProcessTest(true)
	vsFail := &ValidationService{
		chainStartImportProcess: chainFail,
	}
	result = vsFail.ValidateProcessData(valData)
	assert.Equal(t, false, result)
}

func buildValidationChainForImportProcessTest(isFail bool) valModels.Validator {
	// define validators
	requestValidator := &mockTesting.ImportProcessRequestValidator{
		MockResult: !isFail,
	}
	fileTypeValidator := &mockTesting.FileTypeValidator{
		MockResult: true,
	}
	schemaValidator := &mockTesting.SchemaValidator{
		MockResult: true,
	}
	// create chain
	requestValidator.SetNext(fileTypeValidator)
	fileTypeValidator.SetNext(schemaValidator)

	return requestValidator
}

func TestValidationService_ValidateDownloadData(t *testing.T) {
	valData := &valModels.DownloadValidationData{
		ReqData:         models.DownloadRequestData{},
		ValidationFlags: valModels.DownloadValidationFlags{},
	}

	// is valid
	chainSucess := buildValidationChainForDownloadProcessTest(false)
	vsSucess := &ValidationService{
		chainStartDownloadProcess: chainSucess,
	}
	result := vsSucess.ValidateDownloadData(valData, false)
	assert.Equal(t, true, result)

	// is invalid
	chainFail := buildValidationChainForDownloadProcessTest(true)
	vsFail := &ValidationService{
		chainStartDownloadProcess: chainFail,
	}
	result = vsFail.ValidateDownloadData(valData, false)
	assert.Equal(t, false, result)
}

func buildValidationChainForDownloadProcessTest(isFail bool) valModels.Validator {
	// define validators
	requestValidator := &mockTesting.DownloadProcessRequestValidator{
		MockResult: !isFail,
	}

	return requestValidator
}
