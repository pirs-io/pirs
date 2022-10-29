package models

import (
	"bytes"
	"google.golang.org/grpc/codes"
)

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
