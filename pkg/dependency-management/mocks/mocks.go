package mocks

import (
	"golang.org/x/net/context"
	"pirs.io/commons/domain"
	"pirs.io/commons/enums"
	"pirs.io/dependency-management/detection/models"
	"strconv"
	"strings"
)

// A Detector todo
type Detector struct {
	SizeOfReturnedArray int
}

func (md *Detector) Detect(req models.DetectRequestData) []domain.Metadata {
	var result []domain.Metadata
	for i := 0; i < md.SizeOfReturnedArray; i++ {
		result = append(result, *domain.NewMetadata())
	}
	return result
}

func (md *Detector) SetNext(detector models.Detector) {}

func (md *Detector) ExecuteNextIfPresent(req models.DetectRequestData) []domain.Metadata {
	return []domain.Metadata{}
}

func (md *Detector) IsProcessTypeEqual(toCheck enums.ProcessType) bool {
	return true
}

// A Repository todo
type Repository struct {
	Data []domain.Metadata
}

func (r *Repository) InsertOne(ctx context.Context, metadata *domain.Metadata) (interface{}, error) {
	return nil, nil
}

func (r *Repository) FindNewestVersionByURI(ctx context.Context, uri string) (uint32, error) {
	maxVer := 0
	for _, foundDoc := range r.Data {
		if uri == foundDoc.URIWithoutVersion {
			foundVer, _ := strconv.Atoi(foundDoc.SplitURI[4])
			if foundVer > maxVer {
				maxVer = foundVer
			}
		}
	}
	return uint32(maxVer), nil
}

func (r *Repository) FindNewestByURI(ctx context.Context, uri string) (domain.Metadata, error) {
	maxVer := 0
	maxDoc := domain.Metadata{}
	for _, foundDoc := range r.Data {
		if uri == foundDoc.URIWithoutVersion {
			foundVer, _ := strconv.Atoi(foundDoc.SplitURI[4])
			if foundVer > maxVer {
				maxVer = foundVer
				maxDoc = foundDoc
			}
		}
	}
	return maxDoc, nil
}

func (r *Repository) FindByURI(ctx context.Context, uri string) (domain.Metadata, error) {
	for _, foundDoc := range r.Data {
		if uri == foundDoc.URI {
			return foundDoc, nil
		}
	}
	return domain.Metadata{}, nil
}

func (r *Repository) FindAllInPackage(ctx context.Context, packageUri string) ([]domain.Metadata, error) {
	splitProjectUri := strings.Split(packageUri, ".")
	var result []domain.Metadata
	for _, foundDoc := range r.Data {
		if areSlicesEqual(foundDoc.SplitURI[:4], splitProjectUri) {
			result = append(result, foundDoc)
		}
	}
	return result, nil
}

func areSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
