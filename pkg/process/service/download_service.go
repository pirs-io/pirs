package service

import (
	"google.golang.org/grpc/codes"
	"pirs.io/process/domain"
	metadata "pirs.io/process/metadata/service"
	"pirs.io/process/service/models"
	validation "pirs.io/process/validation/service"
)

type DownloadService struct {
	MetadataService   *metadata.MetadataService
	ValidationService *validation.ValidationService
}

func (ds *DownloadService) DownloadProcess(uri string) *models.DownloadProcessResponseData {
	createResponse := func(code codes.Code, m []domain.Metadata) *models.DownloadProcessResponseData {
		return &models.DownloadProcessResponseData{
			Status:   code,
			Metadata: m,
		}
	}

	// validate
	// find in repo metadata
	// resolve deps
	// find additional dependent metadata

	// return result
	customM := domain.NewMetadata()
	customM.CustomData = domain.PetriflowMetadata{
		ProcessIdentifier: "awdawdawd",
	}
	return createResponse(codes.OK, []domain.Metadata{*customM, *domain.NewMetadata(), *domain.NewMetadata(), *domain.NewMetadata(), *domain.NewMetadata(), *domain.NewMetadata()})
}
