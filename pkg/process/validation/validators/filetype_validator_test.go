package validators

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileTypeValidator_Validate(t *testing.T) {
	// todo use prod config
	conf, err := GetTestingConfig("../../.env", &TestingConfig{})
	if err != nil {
		log_test.Fatal().Msgf("Unable to load application config for FileTypeValidatorTest! %s", err)
	}
	validator := NewFileTypeValidator(conf.AllowedFileExtensions, conf.IgnoreWrongExtension)

	valData := buildValidationDataForImportProcess(RESOURCES+PF3, PF3)
	validator.Validate(&valData)
	assert.Equal(t, true, valData.ValidationFlags.IsFileTypeValid)

	valData = buildValidationDataForImportProcess(RESOURCES+PF2, PF2)
	validator.Validate(&valData)
	assert.Equal(t, true, valData.ValidationFlags.IsFileTypeValid)

	valData = buildValidationDataForImportProcess(RESOURCES+BPMN1, BPMN1)
	validator.Validate(&valData)
	assert.Equal(t, true, valData.ValidationFlags.IsFileTypeValid)

	valData = buildValidationDataForImportProcess(RESOURCES+WRONG1, WRONG1)
	validator.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsFileTypeValid)

	valData = buildValidationDataForImportProcess(RESOURCES+WRONG2, WRONG2)
	validator.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsFileTypeValid)
}
