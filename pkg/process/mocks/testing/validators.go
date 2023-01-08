package testing

import (
	"pirs.io/process/validation/models"
)

// ImportRequestValidator request
type ImportRequestValidator struct {
	MockResult bool
	next       models.Validator
}

type DownloadRequestValidator struct {
	MockResult bool
	next       models.Validator
}

func (rv *ImportRequestValidator) Validate(data models.Validable) {
	typedData := data.(*models.ImportValidationData)
	typedData.ValidationFlags.IsRequestValid = rv.MockResult
	rv.ExecuteNextIfPresent(data)
}
func (rv *ImportRequestValidator) ExecuteNextIfPresent(data models.Validable) {
	if rv.next != nil {
		rv.next.Validate(data)
	}
}
func (rv *ImportRequestValidator) SetNext(validator models.Validator) {
	rv.next = validator
}

func (rv *DownloadRequestValidator) Validate(data models.Validable) {
	typedData := data.(*models.DownloadValidationData)
	typedData.ValidationFlags.IsRequestValid = rv.MockResult
	rv.ExecuteNextIfPresent(data)
}
func (rv *DownloadRequestValidator) ExecuteNextIfPresent(data models.Validable) {
	if rv.next != nil {
		rv.next.Validate(data)
	}
}
func (rv *DownloadRequestValidator) SetNext(validator models.Validator) {
	rv.next = validator
}

// FileTypeValidator file type
type FileTypeValidator struct {
	MockResult bool
	next       models.Validator
}

func (ft *FileTypeValidator) Validate(data models.Validable) {
	typedData := data.(*models.ImportValidationData)
	typedData.ValidationFlags.IsFileTypeValid = ft.MockResult
	ft.ExecuteNextIfPresent(data)
}
func (ft *FileTypeValidator) ExecuteNextIfPresent(data models.Validable) {
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

func (sv *SchemaValidator) Validate(data models.Validable) {
	typedData := data.(*models.ImportValidationData)
	typedData.ValidationFlags.IsSchemaValid = sv.MockResult
	sv.ExecuteNextIfPresent(data)
}
func (sv *SchemaValidator) ExecuteNextIfPresent(data models.Validable) {
	if sv.next != nil {
		sv.next.Validate(data)
	}
}
func (sv *SchemaValidator) SetNext(validator models.Validator) {
	sv.next = validator
}
