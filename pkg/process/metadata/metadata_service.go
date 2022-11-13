package metadata

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"pirs.io/process/domain"
	"pirs.io/process/enums"
	"pirs.io/process/metadata/extractor"
	"pirs.io/process/metadata/repository/mongo"
	"pirs.io/process/service/models"
	"time"
)

type MetadataService struct {
	extractor      extractor.MetadataExtractor
	repository     mongo.MetadataRepository
	contextTimeout time.Duration
}

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

func (ms *MetadataService) CreateMetadata(pt enums.ProcessType, v uint32, req models.ImportProcessRequestData) (domain.Metadata, error) {
	return ms.extractor.ExtractMetadata(pt, v, req)
}

func (ms *MetadataService) InsertOne(c context.Context, m *domain.Metadata) (*domain.Metadata, error) {
	ctx, cancel := context.WithTimeout(c, ms.contextTimeout)
	defer cancel()

	m.ID = primitive.NewObjectID()
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	res, err := ms.repository.InsertOne(ctx, m)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (ms *MetadataService) FindNewestByURI(ctx context.Context, uri string) (*domain.Metadata, error) {
	newCtx, cancel := context.WithTimeout(ctx, ms.contextTimeout)
	defer cancel()

	res, err := ms.repository.FindNewestByUri(newCtx, uri)
	if err != nil {
		return res, err
	}
	return res, nil
}
