package extractor

import (
	"bufio"
	"bytes"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"io"
	"os"
	"pirs.io/commons"
	"pirs.io/commons/parsers"
	"pirs.io/process/domain"
	"pirs.io/process/service/models"
	"testing"
)

const (
	RESOURCES                  = "../../resources/testing/"
	PF1                        = "empty_petriflow.xml"
	PF2                        = "v5_petriflow.xml"
	PF3                        = "v6_petriflow.xml"
	BPMN1                      = "empty.bpmn"
	WRONG1                     = "wrong_pf_and_bpmn.xml"
	WRONG2                     = "wrong_random.xml"
	PARTIAL_URI                = "io.pirs.testing"
	BASIC_METADATA_MAPPING     = "../../resources/csv/basic_metadata_mapping.csv"
	PETRIFLOW_METADATA_MAPPING = "../../resources/csv/petriflow_metadata_mapping.csv"
	BPMN_METADATA_MAPPING      = "../../resources/csv/bpmn_metadata_mapping.csv"
)

var (
	log_test = commons.GetLoggerFor("MetadataExtractorTest")
)

func TestMetadataExtractor_ExtractMetadata(t *testing.T) {
	petriflowMapping := parseMetadataMappingFromCsv(PETRIFLOW_METADATA_MAPPING)
	bpmnMapping := parseMetadataMappingFromCsv(BPMN_METADATA_MAPPING)

	testPetriflow1(t, PF2, petriflowMapping)
	testPetriflow2(t, PF3, petriflowMapping)
	// todo use non-empty process
	testBpmn(t, BPMN1, bpmnMapping)
}

func testPetriflow1(t *testing.T, filename string, mapping map[string]string) {
	me := &MetadataExtractor{
		PetriflowCustomDataMapping: mapping,
	}
	request := buildRequestData(RESOURCES+filename, filename)
	metadata := me.ExtractMetadata(request)

	assert.NotEqual(t, primitive.NilObjectID, metadata.ID)
	assert.Equal(t, "UTF-8", metadata.Encoding)

	customData := metadata.CustomData.(*domain.PetriflowMetadata)
	assert.Equal(t, "dashboard", customData.ProcessIdentifier)
	assert.Equal(t, "dashboard", customData.Initials)
	assert.Equal(t, "Domov", customData.Title)
	assert.Equal(t, "false", customData.DefaultRole)
	assert.Equal(t, "false", customData.TransitionRole)
	assert.Equal(t, "system,director,dispatcher,mechanic,driver", customData.Roles)
}

func testPetriflow2(t *testing.T, filename string, mapping map[string]string) {
	me := &MetadataExtractor{
		PetriflowCustomDataMapping: mapping,
	}
	request := buildRequestData(RESOURCES+filename, filename)
	metadata := me.ExtractMetadata(request)

	assert.NotEqual(t, primitive.NilObjectID, metadata.ID)
	assert.Equal(t, "", metadata.Encoding)

	customData := metadata.CustomData.(*domain.PetriflowMetadata)
	assert.Equal(t, "vehicle", customData.ProcessIdentifier)
	assert.Equal(t, "VEH", customData.Initials)
	assert.Equal(t, "Vehicle", customData.Title)
	assert.Equal(t, "true", customData.DefaultRole)
	assert.Equal(t, "false", customData.TransitionRole)
	assert.Equal(t, "system,admin,mechanic", customData.Roles)
}

func testBpmn(t *testing.T, filename string, mapping map[string]string) {
	me := &MetadataExtractor{
		BPMNCustomDataMapping: mapping,
	}
	request := buildRequestData(RESOURCES+filename, filename)
	metadata := me.ExtractMetadata(request)

	assert.NotEqual(t, primitive.NilObjectID, metadata.ID)
	assert.Equal(t, "UTF-8", metadata.Encoding)

	customData := metadata.CustomData.(*domain.BPMNMetadata)
	assert.Equal(t, "Process_1", customData.ProcessIdentifier)
}

func parseMetadataMappingFromCsv(csvPath string) map[string]string {
	csv := parsers.ReadCsvFile(csvPath)
	mapping := map[string]string{}
	for _, row := range csv {
		mapping[row[0]] = row[1]
	}
	return mapping
}

func buildRequestData(path string, filename string) models.ImportProcessRequestData {
	file, err := os.Open(path)
	if err != nil {
		log_test.Fatal().Msg(err.Error())
	}
	reader := bufio.NewReader(file)
	buf := make([]byte, 1024)
	var data []byte
	var totalSize = 0
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				log_test.Fatal().Msg(err.Error())
			}
			break
		}
		data = append(data, buf[:n]...)
		totalSize = totalSize + n
	}

	err = file.Close()
	if err != nil {
		log_test.Fatal().Msg(err.Error())
	}

	return models.ImportProcessRequestData{
		Ctx:             context.Background(),
		PartialUri:      PARTIAL_URI,
		ProcessFileName: filename,
		ProcessData:     *bytes.NewBuffer(data),
		ProcessSize:     totalSize,
	}
}
