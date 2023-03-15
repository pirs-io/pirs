package resolution

import (
	"google.golang.org/grpc/codes"
	"pirs.io/commons/db/mongo"
	"pirs.io/commons/domain"
	"pirs.io/dependency-management/detection/models"
)

// A ResolutionService todo
type ResolutionService struct {
	repository mongo.Repository
}

// NewResolutionService todo
func NewResolutionService(repo mongo.Repository) *ResolutionService {
	service := ResolutionService{
		repository: repo,
	}
	return &service
}

// Resolve todo last must be with nil ID
func (rs *ResolutionService) Resolve(resolveURI string) models.ResolveResponseData {
	return models.ResolveResponseData{
		Status:   codes.OK,
		Metadata: []domain.Metadata{*domain.NewMetadata(), *domain.NewMetadata(), domain.Metadata{}},
	}
}
