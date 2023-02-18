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

// NewDetectionService todo
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
	// find dependencies
	dependencies := ds.detectorChain.Detect(request.ProcessType, request.ProcessData)
	dependencies = append(dependencies, domain.Metadata{})

	// return dependencies
	if len(dependencies) == 0 {
		return models.DetectResponseData{
			Status:   codes.NotFound,
			Metadata: dependencies,
		}
	} else {
		return models.DetectResponseData{
			Status:   codes.OK,
			Metadata: dependencies,
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
