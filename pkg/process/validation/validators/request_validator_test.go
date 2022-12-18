package validators

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestImportProcessRequestValidator_Validate(t *testing.T) {
	rv := &ImportProcessRequestValidator{}

	valData := buildValidationDataForImportProcess(RESOURCES+PF1, PF1)
	rv.Validate(&valData)
	assert.Equal(t, true, valData.ValidationFlags.IsRequestValid)

	// ProcessSize
	valData = buildValidationDataForImportProcess(RESOURCES+PF1, PF1)
	valData.ReqData.ProcessSize = -10
	rv.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsRequestValid)

	// ProcessFileName
	valData = buildValidationDataForImportProcess(RESOURCES+PF1, PF1)
	valData.ReqData.ProcessFileName = "!=+-.pdf"
	rv.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsRequestValid)

	valData = buildValidationDataForImportProcess(RESOURCES+PF1, PF1)
	valData.ReqData.ProcessFileName = "my_process.xml"
	rv.Validate(&valData)
	assert.Equal(t, true, valData.ValidationFlags.IsRequestValid)

	// PartialUri
	valData = buildValidationDataForImportProcess(RESOURCES+PF1, PF1)
	valData.ReqData.PartialUri = "awd.awd.a_wd"
	rv.Validate(&valData)
	assert.Equal(t, true, valData.ValidationFlags.IsRequestValid)

	valData = buildValidationDataForImportProcess(RESOURCES+PF1, PF1)
	valData.ReqData.PartialUri = "awd.aw-d.awd"
	rv.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsRequestValid)

	valData = buildValidationDataForImportProcess(RESOURCES+PF1, PF1)
	valData.ReqData.PartialUri = "awd..awdawd"
	rv.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsRequestValid)

	valData = buildValidationDataForImportProcess(RESOURCES+PF1, PF1)
	valData.ReqData.PartialUri = "awdawdawd"
	rv.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsRequestValid)

	valData = buildValidationDataForImportProcess(RESOURCES+PF1, PF1)
	valData.ReqData.PartialUri = "awd.awd.awd.awd"
	rv.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsRequestValid)
}

func TestDownloadProcessRequestValidator_Validate(t *testing.T) {
	rv := &DownloadProcessRequestValidator{}

	valData := buildValidationDataForDownload("awd.awd.awd.awd:1")
	rv.Validate(&valData)
	assert.Equal(t, true, valData.ValidationFlags.IsRequestValid)

	valData = buildValidationDataForDownload("awdawd.awd.awd:1")
	rv.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsRequestValid)

	valData = buildValidationDataForDownload("awd..awd.awd.awd:1")
	rv.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsRequestValid)

	valData = buildValidationDataForDownload("awd.awd.awd:awd.1")
	rv.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsRequestValid)

	valData = buildValidationDataForDownload(":awd.awd.awd.awd:1")
	rv.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsRequestValid)

	valData = buildValidationDataForDownload("...:")
	rv.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsRequestValid)

	valData = buildValidationDataForDownload("")
	rv.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsRequestValid)

	valData = buildValidationDataForDownload("awd.awd.awd.awd:1d")
	rv.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsRequestValid)

	valData = buildValidationDataForDownload("awd.awd.awd.awd:two")
	rv.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsRequestValid)

	valData = buildValidationDataForDownload("awd.awd.awd:1.awd")
	rv.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsRequestValid)
}

func TestDownloadPackageRequestValidator_Validate(t *testing.T) {
	rv := &DownloadPackageRequestValidator{}

	valData := buildValidationDataForDownload("awd.awd.awd")
	rv.Validate(&valData)
	assert.Equal(t, true, valData.ValidationFlags.IsRequestValid)

	valData = buildValidationDataForDownload("awd.awd.")
	rv.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsRequestValid)

	valData = buildValidationDataForDownload("awd..awd")
	rv.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsRequestValid)

	valData = buildValidationDataForDownload(":awd.awd.awd.awd:1")
	rv.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsRequestValid)

	valData = buildValidationDataForDownload("..")
	rv.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsRequestValid)

	valData = buildValidationDataForDownload("")
	rv.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsRequestValid)
}
