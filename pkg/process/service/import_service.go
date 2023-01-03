package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
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

	resourceChan := make(chan ResourceAdapter)
	responseChan := make(chan error)
	isGoroutineRunning := false
	defer func() {
		close(resourceChan)
		<-responseChan
		close(forResponse)
	}()

	for req := range forRequests {
		// must be inside loop because of context instance
		if !isGoroutineRunning {
			go is.ProcessStorageClient.SaveFiles(req.Ctx, resourceChan, responseChan)
			isGoroutineRunning = true
			if <-responseChan != nil {
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

		// check version
		foundVersion := is.MetadataService.FindNewestVersionByURI(req.Ctx, m.URIWithoutVersion)
		m.UpdateVersion(foundVersion + 1)

		// resolve and save deps
		// todo

		// save file in process-storage
		resource := ResourceAdapter{
			Metadata: m,
			FileData: req.ProcessData.Bytes(),
		}
		resourceChan <- resource
		if <-responseChan != nil {
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

func (is *ImportService) transformRequestDataToValidationData(reqData models.ImportRequestData) *valModels.ImportProcessValidationData {
	return &valModels.ImportProcessValidationData{
		ReqData:         reqData,
		ValidationFlags: valModels.ImportValidationFlags{},
	}
}
