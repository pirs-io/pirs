package service

import (
	metadata "pirs.io/process/metadata/service"
	validation "pirs.io/process/validation/service"
)

type DownloadService struct {
	MetadataService   *metadata.MetadataService
	ValidationService *validation.ValidationService
}

func (ds *DownloadService) DownloadProcess() {
}
