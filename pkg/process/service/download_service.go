package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"pirs.io/process/domain"
	metadata "pirs.io/process/metadata/service"
	"pirs.io/process/service/models"
	valModels "pirs.io/process/validation/models"
	validation "pirs.io/process/validation/service"
)

// A DownloadService handles parsed request data from GRPC server to download process metadata
type DownloadService struct {
	MetadataService   *metadata.MetadataService
	ValidationService *validation.ValidationService
}

// DownloadProcess handles models.DownloadRequestData. If success, it returns models.DownloadResponseData
// with codes.OK and found metadata. Otherwise, a response with error code and nil metadata are returned.
func (ds *DownloadService) DownloadProcess(req *models.DownloadRequestData) *models.DownloadResponseData {
	createResponse := func(code codes.Code, m []domain.Metadata) *models.DownloadResponseData {
		return &models.DownloadResponseData{
			Status:   code,
			Metadata: m,
		}
	}

	// validate
	valData := &valModels.DownloadValidationData{
		ReqData:         *req,
		ValidationFlags: valModels.DownloadValidationFlags{},
	}
	isValid := ds.ValidationService.ValidateDownloadData(valData)
	if !isValid {
		return createResponse(codes.InvalidArgument, nil)
	}

	// find metadata
	foundMetadata := ds.MetadataService.FindByURI(req.Ctx, req.TargetUri)
	if foundMetadata.ID == primitive.NilObjectID {
		return createResponse(codes.NotFound, nil)
	}
	// resolve deps
	// todo

	// find additional dependent metadata
	// todo

	// return result
	return createResponse(codes.OK, []domain.Metadata{foundMetadata})
}

// DownloadPackage todo
func (ds *DownloadService) DownloadPackage(req *models.DownloadRequestData) *models.DownloadResponseData {
	createResponse := func(code codes.Code, m []domain.Metadata) *models.DownloadResponseData {
		return &models.DownloadResponseData{
			Status:   code,
			Metadata: m,
		}
	}

	// validate
	// todo

	// find metadata
	// todo

	// resolve deps
	// todo

	// find additional dependent metadata
	// todo

	// return result
	return createResponse(codes.OK, []domain.Metadata{})
}
