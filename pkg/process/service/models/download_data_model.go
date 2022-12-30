package models

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"pirs.io/process/domain"
)

// A DownloadRequestData holds parsed URI from GRPC request. It's created in GRPC server. To download process metadata,
// example URI is "org.tenant.project.process:1". To download package metadata, example URI is "org.tenant.project".
type DownloadRequestData struct {
	Ctx       context.Context
	TargetUri string
	IsPackage bool
}

// A DownloadResponseData represents response from DownloadService to GRPC server. It's streamed in GRPC server. Contains
// list metadata, that correspond to a process or a package.
type DownloadResponseData struct {
	Status   codes.Code
	Metadata []domain.Metadata
}
