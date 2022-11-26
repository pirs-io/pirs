package extractor

import (
	"bufio"
	"bytes"
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
	RESOURCES                  = "../../resources/testing/metadata/"
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
	//basicMapping := parseMetadataMappingFromCsv(BASIC_METADATA_MAPPING)
	petriflowMapping := parseMetadataMappingFromCsv(PETRIFLOW_METADATA_MAPPING)
	bpmnMapping := parseMetadataMappingFromCsv(BPMN_METADATA_MAPPING)

	isWrong, desc := testPetriflow(PF2, petriflowMapping)
	if isWrong {
		t.Error(desc)
	}

	// todo use non-empty process
	isWrong, desc = testBpmn(BPMN1, bpmnMapping)
	if isWrong {
		t.Error(desc)
	}
}

func testPetriflow(filename string, mapping map[string]string) (bool, string) {
	me := &MetadataExtractor{
		PetriflowCustomDataMapping: mapping,
	}
	wrong := false
	wrongDesc := ""
	request := buildRequestData(RESOURCES+filename, filename)
	metadata := me.ExtractMetadata(request)
	if metadata.ID == primitive.NilObjectID {
		return true, "[" + filename + "] Metadata is corrupted, description: ID should not be NilObjectID"
	}
	if isWrong, desc := checkField("Encoding", metadata.Encoding, "UTF-8"); isWrong {
		wrong = true
		wrongDesc = wrongDesc + desc + ";"
	}

	customData := metadata.CustomData.(*domain.PetriflowMetadata)
	if isWrong, desc := checkField("ProcessIdentifier", customData.ProcessIdentifier, "dashboard"); isWrong {
		wrong = true
		wrongDesc = wrongDesc + desc + ";"
	}
	if isWrong, desc := checkField("Initials", customData.Initials, "dashboard"); isWrong {
		wrong = true
		wrongDesc = wrongDesc + desc + ";"
	}
	if isWrong, desc := checkField("Title", customData.Title, "Domov"); isWrong {
		wrong = true
		wrongDesc = wrongDesc + desc + ";"
	}
	if isWrong, desc := checkField("DefaultRole", customData.DefaultRole, "false"); isWrong {
		wrong = true
		wrongDesc = wrongDesc + desc + ";"
	}
	if isWrong, desc := checkField("TransitionRole", customData.TransitionRole, "false"); isWrong {
		wrong = true
		wrongDesc = wrongDesc + desc + ";"
	}
	if isWrong, desc := checkField("Roles", customData.Roles, "system,director,dispatcher,mechanic,driver"); isWrong {
		wrong = true
		wrongDesc = wrongDesc + desc + ";"
	}

	if wrong {
		return true, "[" + filename + "] Metadata is corrupted, description: " + wrongDesc
	} else {
		return false, ""
	}
}

func testBpmn(filename string, mapping map[string]string) (bool, string) {
	me := &MetadataExtractor{
		BPMNCustomDataMapping: mapping,
	}
	wrong := false
	wrongDesc := ""
	request := buildRequestData(RESOURCES+filename, filename)
	metadata := me.ExtractMetadata(request)
	if metadata.ID == primitive.NilObjectID {
		return true, "[" + filename + "] Metadata is corrupted, description: ID should not be NilObjectID"
	}
	if isWrong, desc := checkField("Encoding", metadata.Encoding, "UTF-8"); isWrong {
		wrong = true
		wrongDesc = wrongDesc + desc + ";"
	}

	customData := metadata.CustomData.(*domain.BPMNMetadata)
	if isWrong, desc := checkField("ProcessIdentifier", customData.ProcessIdentifier, "Process_1"); isWrong {
		wrong = true
		wrongDesc = wrongDesc + desc + ";"
	}

	if wrong {
		return true, "[" + filename + "] Metadata is corrupted, description: " + wrongDesc
	} else {
		return false, ""
	}
}

func checkField(fieldName string, fieldValue string, expectedValue string) (bool, string) {
	if fieldValue != expectedValue {
		return true, fieldName + " should " + expectedValue + ", got " + fieldValue
	} else {
		return false, ""
	}
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
				log.Fatal().Msg(err.Error())
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
