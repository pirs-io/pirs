package detectors

import (
	"bytes"
	"pirs.io/process/domain"
)

type PetriflowDetector struct{}

func (pd *PetriflowDetector) Detect(data bytes.Buffer) []domain.Metadata {
	return []domain.Metadata{}
}
