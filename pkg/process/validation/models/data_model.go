package models

import (
	importModels "pirs.io/process/service/models"
)

// ImportProcessValidationData is wrapper of ImportProcessRequestData. It adds validation flags to request data.
type ImportProcessValidationData struct {
	// ReqData mustn't be a pointer. It has to be a copy.
	ReqData         importModels.ImportProcessRequestData
	ValidationFlags ImportProcessValidationFlags
}

// A ImportProcessValidationFlags contains fields, that are filled by corresponding Validator implementations.
type ImportProcessValidationFlags struct {
	IsRequestValid  bool
	IsFileTypeValid bool
	IsSchemaValid   bool
}

type ImportPackageValidationData struct {
}

// A DownloadValidationData todo
type DownloadValidationData struct {
	ReqData         importModels.DownloadRequestData
	ValidationFlags DownloadValidationFlags
}

// A DownloadValidationFlags todo
type DownloadValidationFlags struct {
	IsRequestValid bool
}
