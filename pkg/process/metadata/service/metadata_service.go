package service

import (
	"golang.org/x/net/context"
	"pirs.io/process/domain"
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

func (ms *MetadataService) CreateMetadata(req models.ImportProcessRequestData) (domain.Metadata, error) {
	return ms.extractor.ExtractMetadata(req)
}

func (ms *MetadataService) InsertOne(c context.Context, m *domain.Metadata) (*domain.Metadata, error) {
	ctx, cancel := context.WithTimeout(c, ms.contextTimeout)
	defer cancel()

	res, err := ms.repository.InsertOne(ctx, m)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (ms *MetadataService) FindNewestVersionByURI(ctx context.Context, uri string) (uint32, error) {
	newCtx, cancel := context.WithTimeout(ctx, ms.contextTimeout)
	defer cancel()

	ver, err := ms.repository.FindNewestVersionByURI(newCtx, uri)
	if err != nil {
		return ver, err
	}
	return ver, nil
}
