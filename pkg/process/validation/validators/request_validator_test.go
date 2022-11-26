package validators

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestImportPackageRequestValidator_Validate(t *testing.T) {
	rv := &ImportProcessRequestValidator{}

	valData := buildValidationData(RESOURCES+PF1, PF1)
	rv.Validate(&valData)
	assert.Equal(t, true, valData.ValidationFlags.IsRequestValid)

	// ProcessSize
	valData = buildValidationData(RESOURCES+PF1, PF1)
	valData.ReqData.ProcessSize = -10
	rv.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsRequestValid)

	// ProcessFileName
	valData = buildValidationData(RESOURCES+PF1, PF1)
	valData.ReqData.ProcessFileName = "!=+-.pdf"
	rv.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsRequestValid)

	valData = buildValidationData(RESOURCES+PF1, PF1)
	valData.ReqData.ProcessFileName = "my_process.xml"
	rv.Validate(&valData)
	assert.Equal(t, true, valData.ValidationFlags.IsRequestValid)

	// PartialUri
	valData = buildValidationData(RESOURCES+PF1, PF1)
	valData.ReqData.PartialUri = "awd.awd.a_wd"
	rv.Validate(&valData)
	assert.Equal(t, true, valData.ValidationFlags.IsRequestValid)

	valData = buildValidationData(RESOURCES+PF1, PF1)
	valData.ReqData.PartialUri = "awd.aw-d.awd"
	rv.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsRequestValid)

	valData = buildValidationData(RESOURCES+PF1, PF1)
	valData.ReqData.PartialUri = "awd..awdawd"
	rv.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsRequestValid)

	valData = buildValidationData(RESOURCES+PF1, PF1)
	valData.ReqData.PartialUri = "awdawdawd"
	rv.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsRequestValid)

	valData = buildValidationData(RESOURCES+PF1, PF1)
	valData.ReqData.PartialUri = "awd.awd.awd.awd"
	rv.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsRequestValid)
}
