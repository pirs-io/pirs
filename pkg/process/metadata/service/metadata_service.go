package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"pirs.io/commons"
	"pirs.io/commons/db/mongo"
	"pirs.io/commons/domain"
	"pirs.io/process/metadata/extractor"
	"pirs.io/process/service/models"
	"time"
)

var (
	log = commons.GetLoggerFor("MetadataService")
)

// A MetadataService is responsible for metadata operations. It stores extractor and repository instances.\
type MetadataService struct {
	extractor      extractor.MetadataExtractor
	repository     mongo.Repository
	contextTimeout time.Duration
}

// NewMetadataService creates instance pointer of MetadataService. At the same time it passes mappings to extractor.
func NewMetadataService(r mongo.Repository, t time.Duration, bMapping map[string]string, pfMapping map[string]string, bpmnMapping map[string]string) *MetadataService {
	return &MetadataService{
		extractor: extractor.MetadataExtractor{
			BasicDataMapping:           bMapping,
			PetriflowCustomDataMapping: pfMapping,
			BPMNCustomDataMapping:      bpmnMapping,
		},
		repository:     r,
		contextTimeout: t,
	}
}

// CreateMetadata is wrapper to extract metadata using extractor.
func (ms *MetadataService) CreateMetadata(req models.ImportRequestData) domain.Metadata {
	return ms.extractor.ExtractMetadata(req)
}

// InsertOne passes given domain.Metadata to the repository layer along with added context timeout.
func (ms *MetadataService) InsertOne(c context.Context, m *domain.Metadata) primitive.ObjectID {
	res, err := ms.repository.InsertOne(c, m)
	if err != nil {
		log.Error().Msg(status.Errorf(codes.Internal, "could not insert metadata into database: %v", err).Error())
		return primitive.NilObjectID
	}
	return res.(primitive.ObjectID)
}

// FindNewestVersionByURI passes given uri to the repository layer along with added context timeout. If repository returns
// an error, 0 is returned. Otherwise, a version from repository is returned.
func (ms *MetadataService) FindNewestVersionByURI(ctx context.Context, uri string) uint32 {
	ver, err := ms.repository.FindNewestVersionByURI(ctx, uri)
	if err != nil {
		log.Error().Msg(status.Errorf(codes.Internal, "could not find the newest version in database: %v", err).Error())
		return 0
	}
	return ver
}

// FindByURI passes given uri to the repository layer along with added context timeout. If repository returns an error,
// empty metadata is returned. Otherwise, found metadata is returned.
func (ms *MetadataService) FindByURI(ctx context.Context, uri string) domain.Metadata {
	found, err := ms.repository.FindByURI(ctx, uri)
	if err != nil {
		log.Error().Msg(status.Errorf(codes.Internal, "could not find the metadata in database: %v", err).Error())
		return domain.Metadata{}
	}
	return found
}

// FindAllInPackage passes given package URI to the repository layer along with added context timeout. If repository
// returns an error, empty metadata is returned. Otherwise, found metadata is returned.
func (ms *MetadataService) FindAllInPackage(ctx context.Context, packageUri string) []domain.Metadata {
	foundList, err := ms.repository.FindAllInPackage(ctx, packageUri)
	if err != nil {
		log.Error().Msg(status.Errorf(codes.Internal, "could not find the metadata in database: %v", err).Error())
		return []domain.Metadata{}
	}
	return foundList
}
