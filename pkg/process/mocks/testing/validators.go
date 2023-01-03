package testing

import (
	"pirs.io/process/validation/models"
)

// ImportProcessRequestValidator request
type ImportProcessRequestValidator struct {
	MockResult bool
	next       models.Validator
}

type DownloadProcessRequestValidator struct {
	MockResult bool
	next       models.Validator
}

func (rv *ImportProcessRequestValidator) Validate(data interface{}) {
	typedData := data.(*models.ImportProcessValidationData)
	typedData.ValidationFlags.IsRequestValid = rv.MockResult
	rv.ExecuteNextIfPresent(data)
}
func (rv *ImportProcessRequestValidator) ExecuteNextIfPresent(data interface{}) {
	if rv.next != nil {
		rv.next.Validate(data)
	}
}
func (rv *ImportProcessRequestValidator) SetNext(validator models.Validator) {
	rv.next = validator
}

func (rv *DownloadProcessRequestValidator) Validate(data interface{}) {
	typedData := data.(*models.DownloadValidationData)
	typedData.ValidationFlags.IsRequestValid = rv.MockResult
	rv.ExecuteNextIfPresent(data)
}
func (rv *DownloadProcessRequestValidator) ExecuteNextIfPresent(data interface{}) {
	if rv.next != nil {
		rv.next.Validate(data)
	}
}
func (rv *DownloadProcessRequestValidator) SetNext(validator models.Validator) {
	rv.next = validator
}

// FileTypeValidator file type
type FileTypeValidator struct {
	MockResult bool
	next       models.Validator
}

func (ft *FileTypeValidator) Validate(data interface{}) {
	typedData := data.(*models.ImportProcessValidationData)
	typedData.ValidationFlags.IsFileTypeValid = ft.MockResult
	ft.ExecuteNextIfPresent(data)
}
func (ft *FileTypeValidator) ExecuteNextIfPresent(data interface{}) {
	if ft.next != nil {
		ft.next.Validate(data)
	}
}
func (ft *FileTypeValidator) SetNext(validator models.Validator) {
	ft.next = validator
}

// SchemaValidator schema
type SchemaValidator struct {
	MockResult bool
	next       models.Validator
}

func (sv *SchemaValidator) Validate(data interface{}) {
	typedData := data.(*models.ImportProcessValidationData)
	typedData.ValidationFlags.IsSchemaValid = sv.MockResult
	sv.ExecuteNextIfPresent(data)
}
func (sv *SchemaValidator) ExecuteNextIfPresent(data interface{}) {
	if sv.next != nil {
		sv.next.Validate(data)
	}
}
func (sv *SchemaValidator) SetNext(validator models.Validator) {
	sv.next = validator
}
