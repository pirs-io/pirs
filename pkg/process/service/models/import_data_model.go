package models

import (
	"bytes"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
)

// ImportProcessRequestData holds parsed values from GRPC stream. It's created in GRPC server.
type ImportProcessRequestData struct {
	Ctx context.Context
	// PartialUri's example is "org.tenant.project". No process id and version involved.
	PartialUri      string
	ProcessFileName string
	ProcessData     bytes.Buffer
	ProcessSize     int
}

// ImportProcessResponseData represents response from ImportService to GRPC server.
type ImportProcessResponseData struct {
	Status codes.Code
}

type ImportPackageRequestData struct {
}

type ImportPackageResponseData struct {
	Status codes.Code
}
