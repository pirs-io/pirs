package models

import (
	"bytes"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"pirs.io/process/domain"
)

// ImportRequestData holds parsed values from GRPC stream. It's created in GRPC server.
type ImportRequestData struct {
	Ctx context.Context
	// PartialUri's example is "org.tenant.project". No process id and version involved.
	PartialUri      string
	ProcessFileName string
	ProcessData     bytes.Buffer
	ProcessSize     int
	IsLast          bool
}

// ImportResponseData represents response from ImportService to GRPC server.
type ImportResponseData struct {
	Status codes.Code
}

// A ResourceAdapter is wrapper for metadata and file data.
type ResourceAdapter struct {
	Metadata domain.Metadata
	FileData []byte
}

// A ResponseAdapter is wrapper for metadata and error.
type ResponseAdapter struct {
	Metadata domain.Metadata
	Err      error
}
