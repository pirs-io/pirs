package models

import (
	importModels "pirs.io/process/service/models"
)

type ValidationFlags struct {
	IsRequestValid   bool
	IsExtensionValid bool
	IsSchemaValid    bool
}

type ImportProcessValidationData struct {
	// cannot be a pointer, we want a copy
	ReqData         importModels.ImportProcessRequestData
	ValidationFlags ValidationFlags
}

type ImportPackageValidationData struct {
}
