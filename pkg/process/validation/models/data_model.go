package models

import (
	importModels "pirs.io/process/service/models"
)

// A ImportProcessValidationFlags contains fields, that are filled by corresponding Validator implementations.
type ImportProcessValidationFlags struct {
	IsRequestValid  bool
	IsFileTypeValid bool
	IsSchemaValid   bool
}

// ImportProcessValidationData is wrapper of ImportProcessRequestData. It adds validation flags to request data.
type ImportProcessValidationData struct {
	// ReqData mustn't be a pointer. It has to be a copy.
	ReqData         importModels.ImportProcessRequestData
	ValidationFlags ImportProcessValidationFlags
}

type ImportPackageValidationData struct {
}

type DownloadProcessValidationFlags struct {
	IsRequestValid bool
}

type DownloadProcessValidationData struct {
	ReqData         importModels.DownloadProcessRequestData
	ValidationFlags DownloadProcessValidationFlags
}
