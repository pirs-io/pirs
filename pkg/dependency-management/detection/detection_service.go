package detection

import (
	"google.golang.org/grpc/codes"
	"pirs.io/commons"
	"pirs.io/dependency-management/detection/detectors"
	"pirs.io/dependency-management/detection/models"
	"pirs.io/process/domain"
)

// DetectionService todo
type DetectionService struct {
	detectorChain models.Detector
}

func NewDetectionService() *DetectionService {
	service := DetectionService{}
	service.detectorChain = buildDetectorChain()
	return &service
}

func buildDetectorChain() models.Detector {
	pd := detectors.PetriflowDetector{}
	bd := detectors.BPMNDetector{}

	bd.SetNext(&pd)

	return &bd
}

// Detect todo
func (ds *DetectionService) Detect(request models.DetectRequestData) models.DetectResponseData {
	// validate bytes
	if !isValidChecksum(request.ProcessData.Bytes(), request.CheckSum) {
		return models.DetectResponseData{
			Status: codes.InvalidArgument,
		}
	}

	// call detector chain and get metadata
	// append empty metadata
	// return array

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
