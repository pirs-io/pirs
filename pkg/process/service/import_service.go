package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"pirs.io/commons"
	"pirs.io/process/db/mongo"
	metadata "pirs.io/process/metadata/service"
	"pirs.io/process/mocks"
	"pirs.io/process/service/models"
	valModels "pirs.io/process/validation/models"
	validation "pirs.io/process/validation/service"
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

func (is *ImportService) ImportProcess(req *models.ImportProcessRequestData) *models.ImportProcessResponseData {
	createResponse := func(code codes.Code) *models.ImportProcessResponseData {
		return &models.ImportProcessResponseData{
			Status: code,
		}
	}
	// validate process data
	valData := transformRequestDataToValidationData(*req)
	isValid := is.ValidationService.ValidateProcessData(valData)
	if !isValid {
		return createResponse(codes.InvalidArgument)
	}
	// create metadata
	m := is.MetadataService.CreateMetadata(*req)
	if m.ID == primitive.NilObjectID {
		return createResponse(codes.Internal)
	}
	// resolve and save deps
	// save file in storage
	_, err := is.ProcessStorageClient.SaveProcessFile(req.ProcessData)
	if err != nil {
		log.Error().Msg(status.Errorf(codes.Internal, "cannot store the process: %v", err).Error())
		return createResponse(codes.Internal)
	}
	// check version
	foundVersion := is.MetadataService.FindNewestVersionByURI(req.Ctx, m.URIWithoutVersion)
	m.UpdateVersion(foundVersion + 1)

	// save metadata
	insertedMetadata := is.MetadataService.InsertOne(req.Ctx, &m)
	if insertedMetadata.ID == primitive.NilObjectID {
		return createResponse(codes.Internal)
	}
	// apply grace period
	// create response
	return createResponse(codes.OK)
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
