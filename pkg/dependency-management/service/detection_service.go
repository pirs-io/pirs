package service

import (
	"google.golang.org/grpc/codes"
	"pirs.io/dependency-management/service/models"
)

type DetectionService struct {
}

func (ds *DetectionService) Detect(request models.DetectRequestData) models.DetectResponseData {
	return models.DetectResponseData{
		Status: codes.OK,
	}
}
