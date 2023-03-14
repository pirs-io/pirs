package models

import (
	"google.golang.org/grpc/codes"
	"pirs.io/commons/domain"
)

// A ResolveResponseData todo
type ResolveResponseData struct {
	Status   codes.Code
	Metadata []domain.Metadata
}
