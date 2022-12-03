package models

import (
	"google.golang.org/grpc/codes"
	"pirs.io/process/domain"
)

type DownloadProcessResponseData struct {
	Status   codes.Code
	Metadata []domain.Metadata
}
