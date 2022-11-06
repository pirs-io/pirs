package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	CreatedAt         time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time          `bson:"updated_at" json:"updated_at"`
	ProcessIdentifier string             `bson:"process_identifier" json:"process_identifier"`
	FileName          string             `bson:"file_name" json:"file_name"`
	FileSize          int                `bson:"file_size" json:"file_size"`
	Encoding          string             `bson:"encoding" json:"encoding"`
	Version           int64              `bson:"version" json:"version"`
	Publisher         string             `bson:"publisher" json:"publisher"`
	DependencyData    DependencyMetadata `bson:"dependency_data" json:"dependency_data"`
	CustomData        map[string]string  `bson:"custom_data" json:"custom_data"`
}
