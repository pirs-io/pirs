package detectors

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"pirs.io/commons/domain"
	"pirs.io/commons/enums"
	"pirs.io/commons/parsers"
	"pirs.io/dependency-management/detection/models"
	"pirs.io/dependency-management/mocks"
	"strings"
	"testing"
)

const (
	CHUNK_SIZE                       = 1024
	RESOURCES                        = "../../resources/"
	API_CSV_PATH                     = RESOURCES + "csv/api_functions_for_detection.csv"
	DATA_WITH_PROTOCOLS_PATH         = RESOURCES + "testing/with_protocols.xml"
	DATA_WITH_API_PATH               = RESOURCES + "testing/with_api.xml"
	DATA_WITH_PROTOCOLS_AND_API_PATH = RESOURCES + "testing/with_protocols_and_api.xml"
)

func TestPetriflowDetector_Detect(t *testing.T) {
	apiForDetection := parseApiForDetectionFromCsv(API_CSV_PATH)
	repo := mocks.Repository{}
	metadata1 := domain.NewMetadata()
	metadata1.SplitURI = [5]string{"myorg", "mytenant", "myproject", "myid", "1"}
	metadata1.BuildURI()
	metadata2 := domain.NewMetadata()
	metadata2.SplitURI = [5]string{"myorg", "mytenant", "myproject", "myid", "2"}
	metadata2.BuildURI()
	metadata3 := domain.NewMetadata()
	metadata3.SplitURI = [5]string{"myorg", "mytenant", "myproject2", "myid", "1"}
	metadata3.BuildURI()
	repo.Data = []domain.Metadata{*metadata1, *metadata2, *metadata3}
	detector := NewPetriflowDetector(&repo, apiForDetection)

	req1 := models.DetectRequestData{
		ProjectUri:  "myorg.mytenant.myproject",
		ProcessType: enums.Petriflow,
		ProcessData: *bytes.NewBuffer(readXmlToBytes(DATA_WITH_PROTOCOLS_PATH)),
	}
	deps1 := detector.Detect(req1)
	assert.Len(t, deps1, 1)
	assert.Equal(t, "myorg.mytenant.myproject.myid:1", deps1[0].URI)

	req2 := models.DetectRequestData{
		ProjectUri:  "myorg.mytenant.myproject",
		ProcessType: enums.Petriflow,
		ProcessData: *bytes.NewBuffer(readXmlToBytes(DATA_WITH_API_PATH)),
	}
	deps2 := detector.Detect(req2)
	assert.Len(t, deps2, 1)
	assert.Equal(t, "myorg.mytenant.myproject.myid:2", deps2[0].URI)

	req3 := models.DetectRequestData{
		ProjectUri:  "myorg.mytenant.myproject",
		ProcessType: enums.Petriflow,
		ProcessData: *bytes.NewBuffer(readXmlToBytes(DATA_WITH_PROTOCOLS_AND_API_PATH)),
	}
	deps3 := detector.Detect(req3)
	assert.Len(t, deps3, 3)
	collectedURIs := collectURIs(deps3)
	assert.Contains(t, collectedURIs, "myorg.mytenant.myproject.myid:1")
	assert.Contains(t, collectedURIs, "myorg.mytenant.myproject.myid:2")
	assert.Contains(t, collectedURIs, "myorg.mytenant.myproject2.myid:1")
}

func collectURIs(m []domain.Metadata) []string {
	var result []string
	for _, elem := range m {
		result = append(result, elem.URI)
	}
	return result
}

func parseApiForDetectionFromCsv(csvPath string) map[string][]string {
	csv := parsers.ReadCsvFile(csvPath, true)
	if csv == nil {
		log.Warn().Msg("CSV " + csvPath + " was not found.")
	}
	result := map[string][]string{}
	for idx, row := range csv {
		// skip header
		if idx == 0 {
			continue
		}
		fromDelim := strings.Replace(row[1], "\\", "", 1)
		untilDelim := strings.Replace(row[2], "\\", "", 1)
		result[row[0]] = []string{fromDelim, untilDelim}
	}
	return result
}

func readXmlToBytes(path string) []byte {
	file, _ := os.Open(path)
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := make([]byte, CHUNK_SIZE)
	var result []byte

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		result = append(result, buffer[:n]...)
	}
	return result
}

func TestPetriflowDetector_IsProcessTypeEqual(t *testing.T) {
	repo := mocks.Repository{}
	detector := NewPetriflowDetector(&repo, map[string][]string{})

	assert.Equal(t, true, detector.IsProcessTypeEqual(enums.Petriflow))
	assert.Equal(t, false, detector.IsProcessTypeEqual(enums.BPMN))
	assert.Equal(t, false, detector.IsProcessTypeEqual(enums.UNKNOWN))
}
