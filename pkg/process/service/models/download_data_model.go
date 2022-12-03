package models

import (
	"google.golang.org/grpc/codes"
	"pirs.io/process/domain"
)

type DownloadProcessRequestData struct {
	Uri string
}

type DownloadProcessResponseData struct {
	Status   codes.Code
	Metadata []domain.Metadata
}
