package service

import "os"

const (
	IMPORT_PROCESS_ROLE = "PROCESS_WRITE"
)

type ImportService struct {
}

type ImportProcessRequestData struct {
	ProcessFile *os.File
}

type ImportProcessResponseData struct {
	Status int8
}

type ImportPackageRequestData struct {
}

type ImportPackageResponseData struct {
	Status int8
}

func (is *ImportService) ImportProcess(req *ImportProcessRequestData) (*ImportProcessResponseData, error) {
	// validate process file
	// file pre-processing
	// create metadata
	// resolve deps
	// save file in storage
	// save metadata
	// apply grace period
	// create response
	return &ImportProcessResponseData{Status: 0}, nil
}

func (is *ImportService) ImportPackage() (*ImportPackageResponseData, error) {
	panic("Not implemented")
}
