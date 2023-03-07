package mongo

import (
	"github.com/stretchr/testify/assert"
	"pirs.io/commons/domain"
	"testing"
)

func TestFilterNewVersionsOnSortedList(t *testing.T) {
	m1 := domain.NewMetadata()
	m1.SplitURI = [5]string{"awd", "awd", "awd", "awd", "1"}

	m2 := domain.NewMetadata()
	m2.SplitURI = [5]string{"awd", "awd", "awd", "awd", "2"}

	m3 := domain.NewMetadata()
	m3.SplitURI = [5]string{"awd", "awd", "awd", "awd", "3"}

	m4 := domain.NewMetadata()
	m4.SplitURI = [5]string{"awd", "awd", "awd", "xxx", "4"}

	m5 := domain.NewMetadata()
	m5.SplitURI = [5]string{"awd", "awd", "awd", "xxx", "5"}

	sortedMetadataList := []domain.Metadata{*m3, *m2, *m1, *m5, *m4}

	filtered := filterNewVersionsOnSortedList(sortedMetadataList)

	assert.Len(t, filtered, 2)
	assert.ElementsMatch(t, filtered[0].SplitURI, [5]string{"awd", "awd", "awd", "awd", "3"})
	assert.ElementsMatch(t, filtered[1].SplitURI, [5]string{"awd", "awd", "awd", "xxx", "5"})

}
