package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"pirs.io/process/enums"
	"strconv"
	"strings"
	"time"
)

type Metadata struct {
	ID                primitive.ObjectID `bson:"_id" json:"id"`
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

type DependencyMetadata struct {
	ParentID primitive.ObjectID   `bson:"parent_id" json:"parent_id"`
	ChildIDs []primitive.ObjectID `bson:"child_ids" json:"child_ids"`
}

type PetriflowMetadata struct {
	ProcessIdentifier string `bson:"process_identifier" json:"process_identifier"`
	Initials          string `bson:"initials" json:"initials"`
	Title             string `bson:"title" json:"title"`
	DefaultRole       string `bson:"default_role" json:"default_role"`
	TransitionRole    string `bson:"transition_role" json:"transition_role"`
	Roles             string `bson:"roles" json:"roles"`
}

type BPMNMetadata struct {
	ProcessIdentifier string `bson:"process_identifier" json:"process_identifier"`
	Todo              string
}

func NewMetadata() *Metadata {
	return &Metadata{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
	}
}

func (m *Metadata) BuildURI() {
	uriWithoutVersion := strings.Join(m.SplitURI[:4], ".")
	uri := uriWithoutVersion + ":" + m.SplitURI[4]

	m.URI = uri
	m.URIWithoutVersion = uriWithoutVersion
}

func (m *Metadata) UpdateVersion(newVersion uint32) {
	m.Version = newVersion
	m.SplitURI[4] = strconv.FormatUint(uint64(newVersion), 10)
	m.BuildURI()
}

func (m *Metadata) UpdateProcessIdentifier(newIdentifier string) {
	m.SplitURI[3] = newIdentifier
	m.BuildURI()
}
