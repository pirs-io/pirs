package determiner

import (
	"github.com/antchfx/xmlquery"
	"os"
	"pirs.io/commons"
	"pirs.io/process/enums"
	"testing"
)

const (
	RESOURCES = "../../resources/testing/metadata-determiner/"
	PF1       = "empty_petriflow.xml"
	PF2       = "v5_petriflow.xml"
	PF3       = "v6_petriflow.xml"
	BPMN1     = "empty.bpmn"
	WRONG1    = "wrong_pf_and_bpmn.xml"
	WRONG2    = "wrong_random.xml"
)

var (
	log = commons.GetLoggerFor("MetadataDeterminerTest")
)

func TestDetermineProcessType(t *testing.T) {
	// PF1
	result := getProcessType(RESOURCES + PF1)
	if result != enums.Petriflow {
		t.Error(PF1 + "should be: " + enums.Petriflow.String() + ", got: " + result.String())
	}

	// PF2
	result = getProcessType(RESOURCES + PF2)
	if result != enums.Petriflow {
		t.Error(PF2 + "should be: " + enums.Petriflow.String() + ", got: " + result.String())
	}

	// PF3
	result = getProcessType(RESOURCES + PF3)
	if result != enums.Petriflow {
		t.Error(PF3 + "should be: " + enums.Petriflow.String() + ", got: " + result.String())
	}

	// BPMN1
	result = getProcessType(RESOURCES + BPMN1)
	if result != enums.BPMN {
		t.Error(BPMN1 + "should be: " + enums.BPMN.String() + ", got: " + result.String())
	}

	// WRONG1
	result = getProcessType(RESOURCES + WRONG1)
	if result != enums.UNKNOWN {
		t.Error(WRONG1 + "should be: " + enums.UNKNOWN.String() + ", got: " + result.String())
	}

	// WRONG2
	result = getProcessType(RESOURCES + WRONG2)
	if result != enums.UNKNOWN {
		t.Error(WRONG2 + "should be: " + enums.UNKNOWN.String() + ", got: " + result.String())
	}
}

func getProcessType(path string) enums.ProcessType {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	doc, err := xmlquery.Parse(file)
	err = file.Close()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	return DetermineProcessType(doc)
}
