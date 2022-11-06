package mongo

import (
	"context"
	"pirs.io/commons/mongo"
	"pirs.io/process/domain"
)

type MetadataRepository struct {
	DB         mongo.Database
	Collection mongo.Collection
}

func NewMetadataRepository(db mongo.Database, collectionName string) *MetadataRepository {
	return &MetadataRepository{db, db.Collection(collectionName)}
}

func (m *MetadataRepository) InsertOne(ctx context.Context, metadata *domain.Metadata) (*domain.Metadata, error) {
	_, err := m.Collection.InsertOne(ctx, metadata)
	if err != nil {
		return metadata, err
	}
	return metadata, nil
}
