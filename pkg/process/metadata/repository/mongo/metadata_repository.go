package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (m *MetadataRepository) FindNewestByUri(ctx context.Context, uri string) (*domain.Metadata, error) {
	var metadata []domain.Metadata
	opts := options.MergeFindOptions(
		options.Find().SetLimit(1),
		options.Find().SetSkip(0),
		options.Find().SetSort(bson.D{{"created_at", -1}}),
	)

	cursor, err := m.Collection.Find(ctx, bson.M{"uri_without_version": uri}, opts)
	if err != nil {
		return nil, err
	}
	if cursor == nil {
		return nil, fmt.Errorf("nil cursor value")
	}

	err = cursor.All(ctx, &metadata)
	if err != nil {
		return nil, err
	}

	if len(metadata) == 0 {
		return nil, nil
	} else {
		return &metadata[0], nil
	}
}
