package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"pirs.io/process/db/mongo"
	metadata "pirs.io/process/metadata/service"
	"pirs.io/process/service/models"
	valModels "pirs.io/process/validation/models"
	validation "pirs.io/process/validation/service"
)

type ImportService struct {
	ProcessStorageClient *StorageService
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
	// check version
	foundVersion := is.MetadataService.FindNewestVersionByURI(req.Ctx, m.URIWithoutVersion)
	m.UpdateVersion(foundVersion + 1)
	// resolve and save deps
	// todo
	// save file in process-storage
	err := is.ProcessStorageClient.SaveFile(req.Ctx, m, req.ProcessData.Bytes())
	if err != nil {
		return createResponse(codes.Aborted)
	}
	// save metadata
	insertedMetadata := is.MetadataService.InsertOne(req.Ctx, &m)
	if insertedMetadata.ID == primitive.NilObjectID {
		return createResponse(codes.Internal)
	}
	// apply grace period
	// todo
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
	// todo
	panic("Not implemented")
}
