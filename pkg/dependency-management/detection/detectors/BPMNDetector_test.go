package detectors

import (
	"github.com/stretchr/testify/assert"
	"pirs.io/commons/enums"
	"pirs.io/dependency-management/mocks"
	"testing"
)

func TestBPMNDetector_IsProcessTypeEqual(t *testing.T) {
	repo := mocks.Repository{}
	detector := NewBPMNDetector(&repo)

	assert.Equal(t, false, detector.IsProcessTypeEqual(enums.Petriflow))
	assert.Equal(t, true, detector.IsProcessTypeEqual(enums.BPMN))
	assert.Equal(t, false, detector.IsProcessTypeEqual(enums.UNKNOWN))
}
