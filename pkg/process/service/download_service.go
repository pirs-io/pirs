package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"pirs.io/commons/domain"
	metadata "pirs.io/process/metadata/service"
	"pirs.io/process/service/models"
	valModels "pirs.io/process/validation/models"
	validation "pirs.io/process/validation/service"
)

// A DownloadService handles parsed request data from GRPC server to download process metadata
type DownloadService struct {
	MetadataService   *metadata.MetadataService
	ValidationService *validation.ValidationService
	DependencyService *DependencyService
}

// DownloadProcesses handles models.DownloadRequestData. If success, it returns models.DownloadResponseData
// with codes.OK and found metadata. Otherwise, a response with error code and nil metadata are returned. todo
func (ds *DownloadService) DownloadProcesses(forRequests <-chan models.DownloadRequestData, forResponse chan<- models.DownloadResponseData) {
	createResponse := func(code codes.Code, m []domain.Metadata) *models.DownloadResponseData {
		return &models.DownloadResponseData{
			Status:   code,
			Metadata: m,
		}
	}
	resourceChanDS := make(chan string)
	responseChanDS := make(chan models.ResponseAdapter)
	isGoroutineRunning := false
	defer func() {
		close(resourceChanDS)
		<-responseChanDS
	}()

	for req := range forRequests {
		// validate
		valData := &valModels.DownloadValidationData{
			ReqData:         req,
			ValidationFlags: valModels.DownloadValidationFlags{},
		}
		isValid := ds.ValidationService.ValidateDownloadData(valData)
		if !isValid {
			forResponse <- *createResponse(codes.InvalidArgument, nil)
			continue
		}

		var metadataList []domain.Metadata
		var dependencyList []domain.Metadata

		// find metadata
		if req.IsPackage {
			foundMetadataList := ds.MetadataService.FindAllInPackage(req.Ctx, req.TargetUri)
			if len(foundMetadataList) == 0 {
				forResponse <- *createResponse(codes.NotFound, nil)
				continue
			}
			metadataList = append(metadataList, foundMetadataList...)
		} else {
			foundMetadata := ds.MetadataService.FindByURI(req.Ctx, req.TargetUri)
			if foundMetadata.ID == primitive.NilObjectID {
				forResponse <- *createResponse(codes.NotFound, nil)
				continue
			}
			metadataList = append(metadataList, foundMetadata)
		}

		// resolve deps
		if !isGoroutineRunning {
			go ds.DependencyService.Resolve(req.Ctx, resourceChanDS, responseChanDS)
			isGoroutineRunning = true
			if (<-responseChanDS).Err != nil {
				forResponse <- *createResponse(codes.Unavailable, nil)
				return
			}
		}
		for _, foundMetadata := range metadataList {
			resourceChanDS <- foundMetadata.URI
			response := <-responseChanDS

			if response.Err != nil {
				forResponse <- *createResponse(codes.Internal, metadataList)
				continue
			}

			if foundMetadata.ID != primitive.NilObjectID {
				dependencyList = append(dependencyList, foundMetadata)
			}
		}
		metadataList = append(metadataList, dependencyList...)
		if len(metadataList) == 0 {
			forResponse <- *createResponse(codes.NotFound, metadataList)
		} else {
			forResponse <- *createResponse(codes.OK, metadataList)
		}
	}
}
