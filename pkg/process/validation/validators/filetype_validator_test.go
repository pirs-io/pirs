package validators

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileTypeValidator_Validate(t *testing.T) {
	// todo use prod config
	conf, err := GetTestingConfig("../../dev.env", &TestingConfig{})
	if err != nil {
		log_test.Fatal().Msgf("Unable to load application config for FileTypeValidatorTest! %s", err)
	}
	validator := NewFileTypeValidator(conf.AllowedFileExtensions, conf.IgnoreWrongExtension)

	valData := buildValidationData(RESOURCES+PF3, PF3)
	validator.Validate(&valData)
	assert.Equal(t, true, valData.ValidationFlags.IsFileTypeValid)

	valData = buildValidationData(RESOURCES+PF2, PF2)
	validator.Validate(&valData)
	assert.Equal(t, true, valData.ValidationFlags.IsFileTypeValid)

	valData = buildValidationData(RESOURCES+BPMN1, BPMN1)
	validator.Validate(&valData)
	assert.Equal(t, true, valData.ValidationFlags.IsFileTypeValid)

	valData = buildValidationData(RESOURCES+WRONG1, WRONG1)
	validator.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsFileTypeValid)

	valData = buildValidationData(RESOURCES+WRONG2, WRONG2)
	validator.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsFileTypeValid)
}
