package determiner

import (
	"github.com/antchfx/xmlquery"
	"pirs.io/process/enums"
)

func DetermineProcessType(doc *xmlquery.Node) enums.ProcessType {
	var (
		finalType = enums.UNKNOWN
		counter   int
	)
	determiners := []func(*xmlquery.Node) enums.ProcessType{
		isPetriflowType,
		isBPMNType,
	}

	for _, f := range determiners {
		maybeType := f(doc)
		if maybeType != enums.UNKNOWN {
			finalType = maybeType
			counter = counter + 1
		}
	}

	if counter > 1 {
		return enums.UNKNOWN
	} else {
		return finalType
	}
}

func isPetriflowType(doc *xmlquery.Node) enums.ProcessType {
	defaultRoleLen := len(xmlquery.Find(doc, "//defaultRole[parent::document]"))
	transitionsLen := len(xmlquery.Find(doc, "//transition[parent::document]"))
	placesLen := len(xmlquery.Find(doc, "//place[parent::document]"))
	arcsLen := len(xmlquery.Find(doc, "//arc[parent::document]"))

	if defaultRoleLen > 0 ||
		transitionsLen > 0 ||
		placesLen > 0 ||
		arcsLen > 0 {
		return enums.Petriflow
	} else {
		return enums.UNKNOWN
	}
}

func isBPMNType(doc *xmlquery.Node) enums.ProcessType {
	processLen := len(xmlquery.Find(doc, "//process"))
	definitionsLen := len(xmlquery.Find(doc, "//definitions"))
	bpmnDiagramLen := len(xmlquery.Find(doc, "//bpmndi:BPMNDiagram|//BPMNDiagram"))

	if processLen > 0 ||
		definitionsLen > 0 ||
		bpmnDiagramLen > 0 {
		return enums.BPMN
	} else {
		return enums.UNKNOWN
	}
}
