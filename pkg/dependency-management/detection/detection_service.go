package detection

import (
	"google.golang.org/grpc/codes"
	"pirs.io/commons"
	"pirs.io/commons/db/mongo"
	"pirs.io/commons/domain"
	"pirs.io/dependency-management/detection/detectors"
	"pirs.io/dependency-management/detection/models"
)

// A DetectionService is a service to handle requests processed by the GRPC server. It contains field detectorChain,
// which is starting point of chain of responsibility pattern and field repository for metadata.
type DetectionService struct {
	repository    mongo.Repository
	detectorChain models.Detector
}

// NewDetectionService creates instance of DetectionService with initialized chain.
func NewDetectionService(repo mongo.Repository, petriflowApi map[string][]string) *DetectionService {
	service := DetectionService{}
	service.detectorChain = buildDetectorChain(repo, petriflowApi)
	return &service
}

func buildDetectorChain(repo mongo.Repository, petriflowApi map[string][]string) models.Detector {
	pd := detectors.NewPetriflowDetector(repo, petriflowApi)
	bd := detectors.NewBPMNDetector(repo)

	bd.SetNext(pd)

	return bd
}

// Detect handles models.DetectRequestData. The request is validated and sent to handlers via chain of responsibility
// pattern. On success a models.DetectResponseData is returned with codes.OK. Otherwise, an error code is returned.
func (ds *DetectionService) Detect(request models.DetectRequestData) models.DetectResponseData {
	// validate bytes
	if !isValidChecksum(request.ProcessData.Bytes(), request.CheckSum) {
		return models.DetectResponseData{
			Status: codes.InvalidArgument,
		}
	}
	// find dependencies
	dependencies := ds.detectorChain.Detect(request)

	// return dependencies
	if len(dependencies) == 0 {
		dependencies = append(dependencies, domain.Metadata{})
		return models.DetectResponseData{
			Status:   codes.NotFound,
			Metadata: dependencies,
		}
	} else {
		dependencies = append(dependencies, domain.Metadata{})
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
