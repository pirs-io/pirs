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

// DownloadProcess handles models.DownloadProcessRequestData. If success, it returns models.DownloadProcessResponseData
// with codes.OK and found metadata. Otherwise, a response with error code and nil metadata are returned.
func (ds *DownloadService) DownloadProcess(req *models.DownloadProcessRequestData) *models.DownloadProcessResponseData {
	createResponse := func(code codes.Code, m []domain.Metadata) *models.DownloadProcessResponseData {
		return &models.DownloadProcessResponseData{
			Status:   code,
			Metadata: m,
		}
	}

	// validate
	valData := ds.transformRequestDataToValidationData(*req)
	isValid := ds.ValidationService.ValidateDownloadData(valData)
	if !isValid {
		return createResponse(codes.InvalidArgument, nil)
	}

	// find metadata
	foundMetadata := ds.MetadataService.FindByURI(req.Ctx, req.Uri)
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

func (ds *DownloadService) transformRequestDataToValidationData(reqData models.DownloadProcessRequestData) *valModels.DownloadProcessValidationData {
	return &valModels.DownloadProcessValidationData{
		ReqData:         reqData,
		ValidationFlags: valModels.DownloadProcessValidationFlags{},
	}
}
