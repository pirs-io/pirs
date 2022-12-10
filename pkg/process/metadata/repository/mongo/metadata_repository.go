package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"pirs.io/process/db/mongo"
	"pirs.io/process/domain"
)

// A MetadataRepository holds DB and Collection instances. It's initialized in config package.
type MetadataRepository struct {
	DB         mongo.Database
	Collection mongo.Collection
}

// NewMetadataRepository creates instance pointer of MetadataRepository
func NewMetadataRepository(db mongo.Database, collectionName string) *MetadataRepository {
	return &MetadataRepository{db, db.Collection(collectionName)}
}

// InsertOne inserts given domain.Metadata into database. It returns ID if success. Otherwise, an error is returned.
func (m *MetadataRepository) InsertOne(ctx context.Context, metadata *domain.Metadata) (interface{}, error) {
	id, err := m.Collection.InsertOne(ctx, metadata)
	if err != nil {
		return id, err
	}
	return id, nil
}

// FindNewestVersionByURI finds the newest version number of process based on given URI. Parameter uri corresponds
// to uri_without_version in database. If success, the version from database is returned. Otherwise, 0 as version along
// with error is returned.
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

// FindByURI finds the metadata of process based on given URI. If success, found metadata is returned. Otherwise, empty
// metadata is returned.
func (m *MetadataRepository) FindByURI(ctx context.Context, uri string) (domain.Metadata, error) {
	var data []domain.Metadata
	opts := options.MergeFindOptions(
		options.Find().SetLimit(1),
		options.Find().SetSkip(0),
	)

	cursor, err := m.Collection.Find(ctx, bson.M{"uri": uri}, opts)
	if err != nil {
		return domain.Metadata{}, err
	}
	if cursor == nil {
		return domain.Metadata{}, fmt.Errorf("nil cursor value")
	}

	err = cursor.All(ctx, &data)
	if err != nil {
		return domain.Metadata{}, err
	}

	if len(data) == 0 {
		return domain.Metadata{}, nil
	} else {
		return data[0], nil
	}
}
