package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"pirs.io/commons"
	"pirs.io/process/domain"
	"pirs.io/process/metadata/extractor"
	"pirs.io/process/metadata/repository/mongo"
	"pirs.io/process/service/models"
	"time"
)

var (
	log = commons.GetLoggerFor("MetadataService")
)

// A MetadataService is responsible for metadata operations. It stores extractor and repository instances.\
type MetadataService struct {
	extractor      extractor.MetadataExtractor
	repository     mongo.MetadataRepository
	contextTimeout time.Duration
}

// NewMetadataService creates instance pointer of MetadataService. At the same time it passes mappings to extractor.
func NewMetadataService(r mongo.MetadataRepository, t time.Duration, bMapping map[string]string, pfMapping map[string]string, bpmnMapping map[string]string) *MetadataService {
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
func (ms *MetadataService) CreateMetadata(req models.ImportProcessRequestData) domain.Metadata {
	return ms.extractor.ExtractMetadata(req)
}

// InsertOne passes given domain.Metadata to the repository layer along with added context timeout.
func (ms *MetadataService) InsertOne(c context.Context, m *domain.Metadata) primitive.ObjectID {
	ctx, cancel := context.WithTimeout(c, ms.contextTimeout)
	defer cancel()

	res, err := ms.repository.InsertOne(ctx, m)
	if err != nil {
		log.Error().Msg(status.Errorf(codes.Internal, "could not insert metadata into database: %v", err).Error())
		return primitive.NilObjectID
	}
	return res.(primitive.ObjectID)
}

// FindNewestVersionByURI passes given uri to the repository layer along with added context timeout. If repository returns
// an error, 0 is returned. Otherwise, a version from repository is returned.
func (ms *MetadataService) FindNewestVersionByURI(ctx context.Context, uri string) uint32 {
	newCtx, cancel := context.WithTimeout(ctx, ms.contextTimeout)
	defer cancel()

	ver, err := ms.repository.FindNewestVersionByURI(newCtx, uri)
	if err != nil {
		log.Error().Msg(status.Errorf(codes.Internal, "could not find the newest version in database: %v", err).Error())
		return 0
	}
	return ver
}
