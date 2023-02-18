package detectors

import (
	"bytes"
	"pirs.io/commons/enums"
	"pirs.io/dependency-management/detection/models"
	"pirs.io/process/domain"
)

// A BPMNDetector todo
type BPMNDetector struct {
	next models.Detector
}

// Detect todo
func (bd *BPMNDetector) Detect(processType enums.ProcessType, data bytes.Buffer) []domain.Metadata {
	if !bd.IsProcessTypeEqual(processType) {
		return bd.ExecuteNextIfPresent(processType, data)
	}

	return []domain.Metadata{}
}

// SetNext todo
func (bd *BPMNDetector) SetNext(detector models.Detector) {
	bd.next = detector
}

// ExecuteNextIfPresent todo
func (bd *BPMNDetector) ExecuteNextIfPresent(processType enums.ProcessType, data bytes.Buffer) []domain.Metadata {
	if bd.next != nil {
		return bd.next.Detect(processType, data)
	}
	return []domain.Metadata{}
}

// IsProcessTypeEqual todo
func (bd *BPMNDetector) IsProcessTypeEqual(toCheck enums.ProcessType) bool {
	return enums.BPMN == toCheck
}
