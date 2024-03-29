package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"pirs.io/commons/enums"
	"strconv"
	"strings"
	"time"
)

// A Metadata represents process metadata. Some field values are generated and some read from the delivered file. This
// struct is also used as MongoDB model
type Metadata struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`
	// SplitURI is array whose elements represent URI parts. URI and URIWithoutVersion are generated by this field.
	SplitURI          [5]string          `bson:"split_uri" json:"split_uri"`
	URI               string             `bson:"uri" json:"uri"`
	URIWithoutVersion string             `bson:"uri_without_version" json:"uri_without_version"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
	FileName          string             `bson:"file_name" json:"file_name"`
	FileSize          int                `bson:"file_size" json:"file_size"`
	Encoding          string             `bson:"encoding" json:"encoding"`
	Version           uint32             `bson:"version" json:"version"`
	Publisher         string             `bson:"publisher" json:"publisher"`
	DependencyData    DependencyMetadata `bson:"dependency_data" json:"dependency_data"`
	ProcessType       enums.ProcessType  `bson:"process_type" json:"process_type"`
	CustomData        interface{}        `bson:"custom_data" json:"custom_data"`
}

// A NestedMetadata represents dependency of Metadata. It contains reduced fields of Metadata for simplicity.
type NestedMetadata struct {
	URI         string            `bson:"uri" json:"uri"`
	CreatedAt   time.Time         `bson:"created_at" json:"created_at"`
	FileName    string            `bson:"file_name" json:"file_name"`
	FileSize    int               `bson:"file_size" json:"file_size"`
	Publisher   string            `bson:"publisher" json:"publisher"`
	ProcessType enums.ProcessType `bson:"process_type" json:"process_type"`
}

// A DependencyMetadata is a wrapper for an array of NestedMetadata.
type DependencyMetadata struct {
	Dependencies []NestedMetadata `bson:"dependecies" json:"dependencies"`
}

// PetriflowMetadata is Metadata.CustomData field.
type PetriflowMetadata struct {
	ProcessIdentifier string `bson:"process_identifier" json:"process_identifier"`
	Initials          string `bson:"initials" json:"initials"`
	Title             string `bson:"title" json:"title"`
	DefaultRole       string `bson:"default_role" json:"default_role"`
	TransitionRole    string `bson:"transition_role" json:"transition_role"`
	Roles             string `bson:"roles" json:"roles"`
}

// BPMNMetadata is CustomData field in Metadata struct.
type BPMNMetadata struct {
	ProcessIdentifier string `bson:"process_identifier" json:"process_identifier"`
	Todo              string
}

// NewMetadata creates Metadata instance pointer, which contains initialized Metadata.ID and Metadata.CreatedAt fields.
func NewMetadata() *Metadata {
	m := &Metadata{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
	}
	m.UpdateVersion(1)

	return m
}

// BuildURI joins elements of SplitURI in URI and URIWithoutVersion.
func (m *Metadata) BuildURI() {
	uriWithoutVersion := strings.Join(m.SplitURI[:4], ".")
	uri := uriWithoutVersion + ":" + m.SplitURI[4]

	m.URI = uri
	m.URIWithoutVersion = uriWithoutVersion
}

// GetProjectURI returns project URI. For example "myorg.mytenant.myproject".
func (m *Metadata) GetProjectURI() string {
	return strings.Join(m.SplitURI[:3], ".")
}

// UpdateVersion updates field Version and SplitURI. Then calls BuildURI.
func (m *Metadata) UpdateVersion(newVersion uint32) {
	m.Version = newVersion
	m.SplitURI[4] = strconv.FormatUint(uint64(newVersion), 10)
	m.BuildURI()
}

// UpdateProcessIdentifier updates field SplitURI. Then calls BuildURI.
func (m *Metadata) UpdateProcessIdentifier(newIdentifier string) {
	m.SplitURI[3] = newIdentifier
	m.BuildURI()
}

// TransformToNestedMetadata transforms Metadata to NestedMetadata. Returns pointer to NestedMetadata instance.
func (m *Metadata) TransformToNestedMetadata() *NestedMetadata {
	return &NestedMetadata{
		URI:         m.URI,
		CreatedAt:   m.CreatedAt,
		FileName:    m.FileName,
		FileSize:    m.FileSize,
		Publisher:   m.Publisher,
		ProcessType: m.ProcessType,
	}
}
