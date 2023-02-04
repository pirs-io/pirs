package models

import (
	"bytes"
	"google.golang.org/grpc/codes"
	"pirs.io/process/domain"
)

// A DetectRequestData todo
type DetectRequestData struct {
	CheckSum    string
	ProcessData bytes.Buffer
}

// A DetectResponseData todo
type DetectResponseData struct {
	Status   codes.Code
	Metadata []domain.Metadata
}
