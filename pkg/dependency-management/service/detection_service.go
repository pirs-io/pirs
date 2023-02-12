package service

import (
	"google.golang.org/grpc/codes"
	"pirs.io/commons"
	"pirs.io/dependency-management/service/models"
	"pirs.io/process/domain"
)

// DetectionService todo
type DetectionService struct{}

// Detect todo
func (ds *DetectionService) Detect(request models.DetectRequestData) models.DetectResponseData {
	isValid := isValidChecksum(request.ProcessData.Bytes(), request.CheckSum)
	if isValid {
		println("detecting dependencies....")
		// todo
		return models.DetectResponseData{
			Status: codes.OK,
			Metadata: []domain.Metadata{
				*domain.NewMetadata(),
				*domain.NewMetadata(),
				*domain.NewMetadata(),
				{},
			},
		}
	} else {
		return models.DetectResponseData{
			Status: codes.InvalidArgument,
		}
	}
}

func isValidChecksum(data []byte, toChecksum string) bool {
	rawHash := commons.HashBytesToSHA256(data)
	checksum := commons.ConvertBytesToString(rawHash)
	if checksum == toChecksum {
		return true
	} else {
		return false
	}
}
