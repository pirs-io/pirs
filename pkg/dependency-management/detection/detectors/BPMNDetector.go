package detectors

import (
	"bytes"
	"pirs.io/commons/db/mongo"
	"pirs.io/commons/domain"
	"pirs.io/commons/enums"
	"pirs.io/dependency-management/detection/models"
)

// A BPMNDetector represents structure for dependency detection of process type enums.BPMN. It contains field next,
// which is a pointer on the next models.Detector within chain of responsibility pattern.
type BPMNDetector struct {
	repository mongo.Repository
	next       models.Detector
}

// NewBPMNDetector return pointer of instance of BPMNDetector. It contains initialized metadata repository.
func NewBPMNDetector(repo mongo.Repository) *BPMNDetector {
	return &BPMNDetector{
		repository: repo,
	}
}

// Detect todo
func (bd *BPMNDetector) Detect(processType enums.ProcessType, data bytes.Buffer) []domain.Metadata {
	if !bd.IsProcessTypeEqual(processType) {
		return bd.ExecuteNextIfPresent(processType, data)
	}

	return []domain.Metadata{}
}

// SetNext sets next models.Detector within chain of responsibility.
func (bd *BPMNDetector) SetNext(detector models.Detector) {
	bd.next = detector
}

// ExecuteNextIfPresent executes next models.Detector if exists.
func (bd *BPMNDetector) ExecuteNextIfPresent(processType enums.ProcessType, data bytes.Buffer) []domain.Metadata {
	if bd.next != nil {
		return bd.next.Detect(processType, data)
	}
	return []domain.Metadata{}
}

// IsProcessTypeEqual checks if enums.ProcessType is equal to handler type
func (bd *BPMNDetector) IsProcessTypeEqual(toCheck enums.ProcessType) bool {
	return enums.BPMN == toCheck
}
