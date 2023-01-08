package models

import (
	importModels "pirs.io/process/service/models"
)

// ImportValidationData is wrapper of ImportRequestData. It adds validation flags to request data.
type ImportValidationData struct {
	// ReqData mustn't be a pointer. It has to be a copy.
	ReqData         importModels.ImportRequestData
	ValidationFlags ImportValidationFlags
}

func (vd *ImportValidationData) IsValidable() {}

// A ImportValidationFlags contains fields, that are filled by corresponding Validator implementations.
type ImportValidationFlags struct {
	IsRequestValid  bool
	IsFileTypeValid bool
	IsSchemaValid   bool
}

// A DownloadValidationData is wrapped of DownloadRequestData. It adds validation flags to request data.
type DownloadValidationData struct {
	ReqData         importModels.DownloadRequestData
	ValidationFlags DownloadValidationFlags
}

func (vd *DownloadValidationData) IsValidable() {}

// A DownloadValidationFlags contains fields, that are filled by corresponding Validator implementations.
type DownloadValidationFlags struct {
	IsRequestValid bool
}
