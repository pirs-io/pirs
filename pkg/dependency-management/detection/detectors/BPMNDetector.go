package detectors

import (
	"bytes"
	"pirs.io/dependency-management/detection/models"
	"pirs.io/process/domain"
)

// A BPMNDetector todo
type BPMNDetector struct {
	next models.Detector
}

// Detect todo
func (bd *BPMNDetector) Detect(data bytes.Buffer) []domain.Metadata {
	return []domain.Metadata{}
}

// SetNext todo
func (bd *BPMNDetector) SetNext(detector models.Detector) {
	bd.next = detector
}

// ExecuteNextIfPresent todo
func (bd *BPMNDetector) ExecuteNextIfPresent(data bytes.Buffer) []domain.Metadata {
	if bd.next != nil {
		return bd.next.Detect(data)
	}
	return []domain.Metadata{}
}
