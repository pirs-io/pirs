package detection

import (
	"google.golang.org/grpc/codes"
	"pirs.io/commons"
	"pirs.io/dependency-management/detection/models"
	"pirs.io/process/domain"
)

// DetectionService todo
type DetectionService struct{}

// Detect todo
func (ds *DetectionService) Detect(request models.DetectRequestData) models.DetectResponseData {
	// validate bytes
	isValid := isValidChecksum(request.ProcessData.Bytes(), request.CheckSum)
	if !isValid {
		return models.DetectResponseData{
			Status: codes.InvalidArgument,
		}
	}

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
