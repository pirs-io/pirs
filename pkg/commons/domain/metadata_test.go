package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetadata_BuildURI(t *testing.T) {
	m := &Metadata{}
	m.SplitURI = [5]string{"com", "example", "project", "my_process", "1"}

	m.BuildURI()

	assert.Equal(t, "com.example.project.my_process:1", m.URI)
	assert.Equal(t, "com.example.project.my_process", m.URIWithoutVersion)
}

func TestMetadata_UpdateVersion(t *testing.T) {
	m := &Metadata{}
	m.Version = uint32(1)
	m.SplitURI = [5]string{"com", "example", "project", "my_process", "1"}

	m.UpdateVersion(uint32(2))

	assert.Equal(t, uint32(2), m.Version)
	assert.Equal(t, "2", m.SplitURI[4])
}

func TestMetadata_UpdateProcessIdentifier(t *testing.T) {
	m := &Metadata{}
	m.SplitURI = [5]string{"com", "example", "project", "my_process", "1"}

	m.UpdateProcessIdentifier("my_newer_process")

	assert.Equal(t, "my_newer_process", m.SplitURI[3])
}
