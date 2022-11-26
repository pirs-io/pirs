package determiner

import (
	"github.com/antchfx/xmlquery"
	"github.com/magiconair/properties/assert"
	"os"
	"pirs.io/commons"
	"pirs.io/process/enums"
	"testing"
)

const (
	RESOURCES = "../../resources/testing/metadata/"
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
	assert.Equal(t, getProcessType(RESOURCES+PF1), enums.Petriflow)
	assert.Equal(t, getProcessType(RESOURCES+PF2), enums.Petriflow)
	assert.Equal(t, getProcessType(RESOURCES+PF3), enums.Petriflow)
	assert.Equal(t, getProcessType(RESOURCES+BPMN1), enums.BPMN)
	assert.Equal(t, getProcessType(RESOURCES+WRONG1), enums.UNKNOWN)
	assert.Equal(t, getProcessType(RESOURCES+WRONG2), enums.UNKNOWN)
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
