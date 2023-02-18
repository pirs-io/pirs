package detectors

import (
	"bytes"
	"pirs.io/process/domain"
)

type BPMNDetector struct{}

func (bd *BPMNDetector) Detect(data bytes.Buffer) []domain.Metadata {
	return []domain.Metadata{}
}
