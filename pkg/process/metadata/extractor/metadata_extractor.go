package extractor

import (
	"bufio"
	"bytes"
	"github.com/antchfx/xmlquery"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"pirs.io/commons"
	"pirs.io/commons/enums"
	"pirs.io/process/domain"
	"pirs.io/process/metadata/determiner"
	"pirs.io/process/service/models"
	"reflect"
	"regexp"
	"strings"
)

var (
	log = commons.GetLoggerFor("MetadataExtractor")
)

// A MetadataExtractor is created in config package. It holds mappings, that contain XPATHs for process files.
type MetadataExtractor struct {
	BasicDataMapping           map[string]string
	PetriflowCustomDataMapping map[string]string
	BPMNCustomDataMapping      map[string]string
}

// ExtractMetadata extracts values from models.ImportProcessRequestData and file content, which is stored in request.
// If enums.ProcessType is determined successfully, a CustomData is initialized by mapping. Otherwise, CustomData is empty.
// If an error occurs, empty domain.Metadata is returned. Otherwise, initialized domain.Metadata is returned.
func (me *MetadataExtractor) ExtractMetadata(req models.ImportRequestData) domain.Metadata {
	data := req.ProcessData.Bytes()
	doc, err := xmlquery.Parse(bytes.NewReader(data))
	if err != nil {
		log.Error().Msg(status.Errorf(codes.Internal, "could not prepare xml data to parse: %v", err).Error())
		return domain.Metadata{}
	}
	m := domain.NewMetadata()

	// extracted basic data
	m.Encoding = me.extractEncodingFromHeaders(&data)
	me.setValuesUsingReflection(reflect.ValueOf(m), me.BasicDataMapping, doc)

	// not extracted basic data
	splitPartialUri := strings.Split(req.PartialUri, ".")
	for idx, part := range splitPartialUri {
		m.SplitURI[idx] = part
	}
	m.FileName = req.ProcessFileName
	m.FileSize = req.ProcessSize
	m.Publisher = "not implemented"
	m.ProcessType = determiner.DetermineProcessType(doc)

	// dependency data
	dd := &domain.DependencyMetadata{}
	m.DependencyData = *dd

	// custom data
	if m.ProcessType == enums.Petriflow {
		customData := me.extractCustomData(m.ProcessType, data, me.PetriflowCustomDataMapping, doc).(*domain.PetriflowMetadata)
		m.CustomData = customData
		m.UpdateProcessIdentifier(customData.ProcessIdentifier)
	} else if m.ProcessType == enums.BPMN {
		customData := me.extractCustomData(m.ProcessType, data, me.BPMNCustomDataMapping, doc).(*domain.BPMNMetadata)
		m.CustomData = customData
		m.UpdateProcessIdentifier(customData.ProcessIdentifier)
	}

	return *m
}

func (me *MetadataExtractor) extractHeaders(data *[]byte) string {
	bytesReader := bytes.NewReader(*data)
	bufReader := bufio.NewReader(bytesReader)
	headerLine, err := bufReader.ReadString('\n')
	if err != nil {
		return ""
	}
	return headerLine
}

func (me *MetadataExtractor) extractEncodingFromHeaders(data *[]byte) string {
	headers := me.extractHeaders(data)
	regex := regexp.MustCompile("(encoding=\")(.*)")
	match := regex.FindString(headers)
	if len(match) == 0 {
		return match
	} else {
		return strings.Split(match, "\"")[1]
	}
}

func (me *MetadataExtractor) extractCustomData(pt enums.ProcessType, data []byte, mapping map[string]string, doc *xmlquery.Node) interface{} {
	if doc == nil {
		var err error
		doc, err = xmlquery.Parse(bytes.NewReader(data))
		if err != nil {
			return nil
		}
	}
	if pt == enums.Petriflow {
		customMetadata := &domain.PetriflowMetadata{}
		me.setValuesUsingReflection(reflect.ValueOf(customMetadata), mapping, doc)
		return customMetadata
	} else if pt == enums.BPMN {
		customMetadata := &domain.BPMNMetadata{}
		me.setValuesUsingReflection(reflect.ValueOf(customMetadata), mapping, doc)
		return customMetadata
	} else {
		return nil
	}
}

func (me *MetadataExtractor) setValuesUsingReflection(myReflectedStruct reflect.Value, mapping map[string]string, doc *xmlquery.Node) {
	ms := myReflectedStruct.Elem()
	if ms.Kind() == reflect.Struct {
		for field, xpath := range mapping {
			f := ms.FieldByName(field)
			value := me.getValueByXPath(doc, xpath)
			if f.IsValid() && f.CanSet() && f.Kind() == reflect.String {
				f.SetString(value)
			}
		}
	}
}

func (me *MetadataExtractor) getValueByXPath(doc *xmlquery.Node, expr string) string {
	nodes := xmlquery.Find(doc, expr)
	if nodes == nil || len(nodes) == 0 {
		return ""
	}
	return me.extractValueFromNodes(nodes)
}

func (me *MetadataExtractor) countByXpath(doc *xmlquery.Node, expr string) int {
	nodes := xmlquery.Find(doc, expr)
	return len(nodes)
}

func (me *MetadataExtractor) extractValueFromNodes(n []*xmlquery.Node) string {
	if len(n) == 1 {
		return n[0].InnerText()
	} else {
		var values []string
		for _, node := range n {
			values = append(values, node.InnerText())
		}
		return strings.Join(values, ",")
	}
}
