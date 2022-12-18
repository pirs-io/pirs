package models

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"pirs.io/process/domain"
)

// A DownloadRequestData holds parsed Uri from GRPC request. It's created in GRPC server. todo
type DownloadRequestData struct {
	Ctx       context.Context
	TargetUri string
}

// A DownloadResponseData represents response from DownloadService to GRPC server. It's streamed in GRPC server. todo
type DownloadResponseData struct {
	Status   codes.Code
	Metadata []domain.Metadata
}
