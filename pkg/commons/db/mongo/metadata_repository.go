package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"pirs.io/commons/domain"
	"strings"
)

// A MetadataRepository holds DB and Collection instances. It's initialized in config package.
type MetadataRepository struct {
	DB         Database
	Collection Collection
}

// NewMetadataRepository creates instance pointer of MetadataRepository
func NewMetadataRepository(db Database, collectionName string) *MetadataRepository {
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

// FindAllInPackage finds all the metadata based on given package URI. Package URI is URI without process identifier and
// process version. Versions of metadata are the newest.
func (m *MetadataRepository) FindAllInPackage(ctx context.Context, packageUri string) ([]domain.Metadata, error) {
	var data []domain.Metadata
	opts := options.MergeFindOptions(
		options.Find().SetSkip(0),
		options.Find().SetSort(bson.D{{"created_at", -1}}),
	)

	//cursor, err := m.Collection.Find(ctx, bson.M{"uri_without_version": primitive.Regex{Pattern: packageUri, Options: ""}}, opts)
	splitPackageUri := strings.Split(packageUri, ".")
	cursor, err := m.Collection.Find(ctx, bson.M{"split_uri": bson.M{"$all": splitPackageUri}}, opts)
	if err != nil {
		return []domain.Metadata{}, err
	}
	if cursor == nil {
		return []domain.Metadata{}, fmt.Errorf("nil cursor value")
	}

	err = cursor.All(ctx, &data)
	if err != nil {
		return []domain.Metadata{}, err
	}

	if len(data) == 0 {
		return []domain.Metadata{}, nil
	} else {
		return filterNewVersionsOnSortedList(data), nil
	}
}

func filterNewVersionsOnSortedList(metadataList []domain.Metadata) []domain.Metadata {
	var visitedIdentifiers []string
	var filteredMetadata []domain.Metadata

	for _, metadata := range metadataList {
		if contains(visitedIdentifiers, metadata.SplitURI[3]) {
			continue
		} else {
			visitedIdentifiers = append(visitedIdentifiers, metadata.SplitURI[3])
			filteredMetadata = append(filteredMetadata, metadata)
		}
	}
	return filteredMetadata
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
