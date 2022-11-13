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

func (m *MetadataRepository) FindNewestVersionByURI(ctx context.Context, uri string) (uint32, error) {
	var data []domain.Metadata
	opts := options.MergeFindOptions(
		options.Find().SetLimit(1),
		options.Find().SetSkip(0),
		options.Find().SetSort(bson.D{{"created_at", -1}}),
	)

	cursor, err := m.Collection.Find(ctx, bson.M{"uri_without_version": uri}, opts)
	if err != nil {
		return uint32(0), err
	}
	if cursor == nil {
		return uint32(0), fmt.Errorf("nil cursor value")
	}

	err = cursor.All(ctx, &data)
	if err != nil {
		return uint32(0), err
	}

	if len(data) == 0 {
		return uint32(0), nil
	} else {
		return data[0].Version, nil
	}
}
