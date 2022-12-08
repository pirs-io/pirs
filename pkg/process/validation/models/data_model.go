package models

import (
	importModels "pirs.io/process/service/models"
)

// A ValidationFlags contains fields, that are filled by corresponding Validator implementations.
type ValidationFlags struct {
	IsRequestValid  bool
	IsFileTypeValid bool
	IsSchemaValid   bool
}

// ImportProcessValidationData is wrapper of ImportProcessRequestData. It adds validation flags to request data.
type ImportProcessValidationData struct {
	// ReqData mustn't be a pointer. It has to be a copy.
	ReqData         importModels.ImportProcessRequestData
	ValidationFlags ValidationFlags
}

type ImportPackageValidationData struct {
}
