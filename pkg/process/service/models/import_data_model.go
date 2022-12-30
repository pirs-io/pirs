package models

import (
	"bytes"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
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
