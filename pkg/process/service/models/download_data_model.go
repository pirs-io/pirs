package models

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"pirs.io/process/domain"
)

type DownloadProcessRequestData struct {
	Ctx context.Context
	Uri string
}

type DownloadProcessResponseData struct {
	Status   codes.Code
	Metadata []domain.Metadata
}
