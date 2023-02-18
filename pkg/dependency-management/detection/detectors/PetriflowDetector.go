package detectors

import (
	"bytes"
	"pirs.io/dependency-management/detection/models"
	"pirs.io/process/domain"
)

// A PetriflowDetector todo
type PetriflowDetector struct {
	next models.Detector
}

// Detect todo
func (pd *PetriflowDetector) Detect(data bytes.Buffer) []domain.Metadata {
	// cez Xpath detekuj v data actions elementy (len v tych mozu byt dependencies). Pre kazdu najdenu akciu vytvor
	// gorutinu, ktora prezre content akcie a vrati result (bud najdene URI alebo ziadne URI - kedze budeme predpokladat protokol)
	// tieto zistenia vrati cez kanal do tejto funkcie

	// na zaklade URIs vyhladaj v mongu metadata a zoskup ich

	return []domain.Metadata{}
}

// SetNext todo
func (pd *PetriflowDetector) SetNext(detector models.Detector) {
	pd.next = detector
}

// ExecuteNextIfPresent todo
func (pd *PetriflowDetector) ExecuteNextIfPresent(data bytes.Buffer) []domain.Metadata {
	if pd.next != nil {
		return pd.next.Detect(data)
	}
	return []domain.Metadata{}
}
