package service

import (
	"bytes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"pirs.io/commons"
	"pirs.io/process/mock"
)

var (
	log = commons.GetLoggerFor("processGrpc")
)

const (
	IMPORT_PROCESS_ROLE = "PROCESS_WRITE"
)

type ImportService struct {
	// todo mockup
	ProcessStorageClient *mock.DiskProcessStore
	ValidationService    *ValidationService
}

type ImportProcessRequestData struct {
	ProcessFileName string
	ProcessData     bytes.Buffer
	ProcessSize     int
}

type ImportProcessResponseData struct {
	Status codes.Code
}

type ImportPackageRequestData struct {
}

type ImportPackageResponseData struct {
	Status codes.Code
}

func (is *ImportService) ImportProcess(req *ImportProcessRequestData) (*ImportProcessResponseData, error) {
	// validate process data
	isValid, err := is.ValidationService.ValidateProcessData(*req)
	if err != nil {
		return &ImportProcessResponseData{
			Status: codes.Internal,
		}, err
	}
	if !isValid {
		return &ImportProcessResponseData{
			Status: codes.InvalidArgument,
		}, nil
	}
	// file pre-processing
	// create metadata
	// resolve deps
	// save file in storage
	_, err = is.ProcessStorageClient.SaveProcessFile(req.ProcessData)
	if err != nil {
		log.Error().Msg(status.Errorf(codes.Internal, "cannot store the process: %v", err).Error())
		return nil, err
	}
	// save metadata
	// apply grace period
	// create response
	return &ImportProcessResponseData{
		Status: codes.OK,
	}, nil
}

func (is *ImportService) ImportPackage() (*ImportPackageResponseData, error) {
	panic("Not implemented")
}
