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

// An ImportService handles parsed request data from GRPC server to import process files.
type ImportService struct {
	ProcessStorageClient *StorageService
	ValidationService    *validation.ValidationService
	MetadataService      *metadata.MetadataService
	DependencyService    *DependencyService
}

// ImportProcesses runs in a separate goroutine. It's created from GRPC server endpoint method. It waits for requests
// coming through forRequests channel. Once a request is received, it handles it and sends to the lower level services.
// On success, codes.OK is sent through the forResponse channel. On failure, non-codes.OK is sent through the channel.
func (is *ImportService) ImportProcesses(forRequests <-chan models.ImportRequestData, forResponse chan<- models.ImportResponseData) {
	createResponse := func(code codes.Code) models.ImportResponseData {
		return models.ImportResponseData{
			Status: code,
		}
	}

	resourceChanSS := make(chan models.ResourceAdapter)
	responseChanSS := make(chan error)
	resourceChanDS := make(chan models.DetectResourceAdapter)
	responseChanDS := make(chan models.ResponseAdapter)
	isGoroutineRunning := false
	defer func() {
		close(resourceChanDS)
		<-responseChanDS
		close(resourceChanSS)
		<-responseChanSS
		close(forResponse)
	}()

	for req := range forRequests {
		// must be inside loop because of context instance
		if !isGoroutineRunning {
			go is.ProcessStorageClient.SaveFiles(req.Ctx, resourceChanSS, responseChanSS)
			go is.DependencyService.Detect(req.Ctx, resourceChanDS, responseChanDS)
			isGoroutineRunning = true

			if (<-responseChanDS).Err != nil {
				forResponse <- createResponse(codes.Unavailable)
				return
			}
			if <-responseChanSS != nil {
				forResponse <- createResponse(codes.Unavailable)
				return
			}
		}
		// validate process data
		valData := is.transformRequestDataToValidationData(req)
		isValid := is.ValidationService.ValidateProcessData(valData)
		if !isValid {
			forResponse <- createResponse(codes.InvalidArgument)
			return
		}

		// create metadata
		m := is.MetadataService.CreateMetadata(req)
		if m.ID == primitive.NilObjectID {
			forResponse <- createResponse(codes.Internal)
			return
		}

		// resolve and save deps
		resourceChanDS <- models.DetectResourceAdapter{
			ProcessType: m.ProcessType,
			FileData:    req.ProcessData.Bytes(),
		}
		var currentDependencies []domain.NestedMetadata
		for {
			respDs := <-responseChanDS
			if respDs.Err != nil {
				forResponse <- createResponse(codes.Aborted)
				return
			}
			if respDs.Metadata.ID == primitive.NilObjectID {
				break
			} else {
				nested := respDs.Metadata.TransformToNestedMetadata()
				currentDependencies = append(currentDependencies, *nested)
			}
		}
		m.DependencyData = domain.DependencyMetadata{Dependencies: currentDependencies}

		// check version
		foundVersion := is.MetadataService.FindNewestVersionByURI(req.Ctx, m.URIWithoutVersion)
		m.UpdateVersion(foundVersion + 1)

		// save file in process-storage
		resource := models.ResourceAdapter{
			Metadata: m,
			FileData: req.ProcessData.Bytes(),
		}
		resourceChanSS <- resource
		if <-responseChanSS != nil {
			forResponse <- createResponse(codes.Aborted)
			return
		}

		// save metadata
		insertedID := is.MetadataService.InsertOne(req.Ctx, &m)
		if insertedID == primitive.NilObjectID {
			forResponse <- createResponse(codes.Internal)
			return
		} else {
			forResponse <- createResponse(codes.OK)
		}
	}
}

func (is *ImportService) transformRequestDataToValidationData(reqData models.ImportRequestData) *valModels.ImportValidationData {
	return &valModels.ImportValidationData{
		ReqData:         reqData,
		ValidationFlags: valModels.ImportValidationFlags{},
	}
}
