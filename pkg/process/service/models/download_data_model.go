package models

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"pirs.io/process/domain"
)

// A DownloadProcessRequestData holds parsed Uri from GRPC request. It's created in GRPC server.
type DownloadProcessRequestData struct {
	Ctx context.Context
	Uri string
}

// A DownloadProcessResponseData represents response from DownloadService to GRPC server. It's streamed in GRPC server.
type DownloadProcessResponseData struct {
	Status   codes.Code
	Metadata []domain.Metadata
}

// A DownloadPackageRequestData todo
type DownloadPackageRequestData struct {
	Ctx        context.Context
	PartialUri string
}

// A DownloadPackageResponseData todo
type DownloadPackageResponseData struct {
	Status   codes.Code
	Metadata []domain.Metadata
}
