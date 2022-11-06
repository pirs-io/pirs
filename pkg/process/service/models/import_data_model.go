package models

import (
	"bytes"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
)

type ImportProcessRequestData struct {
	Ctx             context.Context
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
