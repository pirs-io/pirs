package models

import (
	"bytes"
	"google.golang.org/grpc/codes"
	"pirs.io/process/domain"
)

// A DetectRequestData represents input inside service.DetectionService. The instance is created in GRPC server. A CheckSum
// is checksum of ProcessData.
type DetectRequestData struct {
	CheckSum    string
	ProcessData bytes.Buffer
}

// A DetectResponseData represents output for GRPC server from service.DetectionService. Metadata is array of metadata,
// that are dependent on input process data.
type DetectResponseData struct {
	Status   codes.Code
	Metadata []domain.Metadata
}
