package service

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"pirs.io/commons"
	"pirs.io/commons/mongo"
	"pirs.io/process/enums"
	"pirs.io/process/metadata"
	"pirs.io/process/mocks"
	"pirs.io/process/service/models"
	"pirs.io/process/validation"
	valModels "pirs.io/process/validation/models"
)

var (
	log = commons.GetLoggerFor("processGrpc")
)

type ImportService struct {
	// todo mockup
	ProcessStorageClient *mocks.DiskProcessStore
	MongoClient          *mongo.Client
	ValidationService    *validation.ValidationService
	MetadataService      *metadata.MetadataService
}

func (is *ImportService) ImportProcess(req *models.ImportProcessRequestData) (*models.ImportProcessResponseData, error) {
	// validate process data
	valData := transformRequestDataToValidationData(*req)
	isValid, err := is.ValidationService.ValidateProcessData(valData)
	if err != nil {
		return &models.ImportProcessResponseData{
			Status: codes.Internal,
		}, err
	}
	if !isValid {
		return &models.ImportProcessResponseData{
			Status: codes.InvalidArgument,
		}, nil
	}
	// file pre-processing
	// todo determine new or update (version,...)(by uri)
	// create metadata
	m, _ := is.MetadataService.CreateMetadata(enums.Petriflow, 0, *req)
	// resolve and save deps
	// save file in storage
	_, err = is.ProcessStorageClient.SaveProcessFile(req.ProcessData)
	if err != nil {
		log.Error().Msg(status.Errorf(codes.Internal, "cannot store the process: %v", err).Error())
		return nil, err
	}
	// save metadata
	_, err = is.MetadataService.InsertOne(req.Ctx, &m)
	if err != nil {
		return nil, err
	}
	// apply grace period
	// create response
	return &models.ImportProcessResponseData{
		Status: codes.OK,
	}, nil
}

func transformRequestDataToValidationData(reqData models.ImportProcessRequestData) *valModels.ImportProcessValidationData {
	return &valModels.ImportProcessValidationData{
		ReqData:         reqData,
		ValidationFlags: valModels.ValidationFlags{},
	}
}

func (is *ImportService) ImportPackage() (*models.ImportPackageResponseData, error) {
	panic("Not implemented")
}
