package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"pirs.io/process/enums"
	"strconv"
	"strings"
	"time"
)

// basic metadata
type DependencyMetadata struct {
	ParentID primitive.ObjectID   `bson:"parent_id" json:"parent_id"`
	ChildIDs []primitive.ObjectID `bson:"child_ids" json:"child_ids"`
}

type Metadata struct {
	ID                primitive.ObjectID `bson:"_id" json:"id"`
	URI               string             `bson:"uri" json:"uri"`
	URIWithoutVersion string             `bson:"uri_without_version" json:"uri_without_version"`
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
	ProcessIdentifier string             `bson:"process_identifier" json:"process_identifier"`
	FileName          string             `bson:"file_name" json:"file_name"`
	FileSize          int                `bson:"file_size" json:"file_size"`
	Encoding          string             `bson:"encoding" json:"encoding"`
	Version           uint32             `bson:"version" json:"version"`
	Publisher         string             `bson:"publisher" json:"publisher"`
	DependencyData    DependencyMetadata `bson:"dependency_data" json:"dependency_data"`
	ProcessType       enums.ProcessType  `bson:"process_type" json:"process_type"`
	CustomData        interface{}        `bson:"custom_data" json:"custom_data"`
}

type PetriflowMetadata struct {
	Initials       string `bson:"initials" json:"initials"`
	Title          string `bson:"title" json:"title"`
	DefaultRole    string `bson:"default_role" json:"default_role"`
	TransitionRole string `bson:"transition_role" json:"transition_role"`
	Roles          string `bson:"roles" json:"roles"`
}

type BPMNMetadata struct {
	Todo string `bson:"todo" json:"todo"`
}

func NewMetadata() *Metadata {
	return &Metadata{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func NewDependencyMetadata() *DependencyMetadata {
	return &DependencyMetadata{}
}

func NewPetriflowMetadata() *PetriflowMetadata {
	return &PetriflowMetadata{}
}

func NewBPMNMetadata() *BPMNMetadata {
	return &BPMNMetadata{}
}

func (m *Metadata) UpdateVersion(newVersion uint32) {
	splitUri := strings.Split(m.URI, ":")
	m.URI = splitUri[0] + ":" + strconv.FormatUint(uint64(newVersion), 10)
	m.Version = newVersion
}
