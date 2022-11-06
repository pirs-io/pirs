package extractor

import (
	"bufio"
	"bytes"
	"github.com/antchfx/xmlquery"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"pirs.io/process/domain"
	"pirs.io/process/enums"
	"pirs.io/process/service/models"
	"regexp"
	"strings"
)

type MetadataExtractor struct {
	PetriflowCustomDataXpaths []string
	BPMNCustomDataXpaths      []string
}

func (me *MetadataExtractor) ExtractMetadata(pt enums.ProcessType, v int64, req models.ImportProcessRequestData) (domain.Metadata, error) {
	data := req.ProcessData.Bytes()
	doc, err := xmlquery.Parse(bytes.NewReader(data))
	if err != nil {
		return domain.Metadata{}, err
	}
	// basic data
	id := xmlquery.FindOne(doc, "//id")
	// dependency data
	// todo
	// custom data
	customData := extractCustomData(data, me.PetriflowCustomDataXpaths, doc)

	// todo specify some data
	return domain.Metadata{
		URI:               "not implemented",
		ProcessIdentifier: id.InnerText(),
		FileName:          req.ProcessFileName,
		FileSize:          req.ProcessSize,
		Encoding:          extractEncodingFromHeaders(&data),
		Version:           v + 1,
		Publisher:         "not implemented",
		DependencyData: domain.DependencyMetadata{
			ParentID: primitive.ObjectID{},
			ChildIDs: []primitive.ObjectID{primitive.ObjectID{}, primitive.ObjectID{}},
		},
		CustomData: customData,
	}, nil
}

func extractHeaders(data *[]byte) string {
	bytesReader := bytes.NewReader(*data)
	bufReader := bufio.NewReader(bytesReader)
	headerLine, err := bufReader.ReadString('\n')
	if err != nil {
		return ""
	}
	return headerLine
}

func extractEncodingFromHeaders(data *[]byte) string {
	headers := extractHeaders(data)
	regex := regexp.MustCompile("(encoding=\")(.*)")
	match := regex.FindString(headers)
	if len(match) == 0 {
		return match
	} else {
		return strings.Split(match, "\"")[1]
	}
}

func extractCustomData(data []byte, dataExprs []string, doc *xmlquery.Node) map[string]string {
	result := map[string]string{}
	if doc == nil {
		var err error
		doc, err = xmlquery.Parse(bytes.NewReader(data))
		if err != nil {
			return result
		}
	}
	for _, expr := range dataExprs {
		nodes := xmlquery.Find(doc, expr)
		if nodes == nil || len(nodes) == 0 {
			continue
		}
		key := transformXpathToAttrs(expr)
		if len(nodes) == 1 {
			result[key] = nodes[0].InnerText()
		} else {
			values := []string{}
			for _, node := range nodes {
				values = append(values, node.InnerText())
			}
			value := strings.Join(values, ",")
			result[key] = value
		}
	}
	return result
}

func transformXpathToAttrs(expr string) string {
	transformed := strings.ReplaceAll(expr, "//", ".")
	transformed = strings.TrimPrefix(transformed, ".")
	regex := regexp.MustCompile("\\[(.*)\\]")
	substrsToRemove := regex.FindAllString(transformed, -1)
	for _, substr := range substrsToRemove {
		transformed = strings.ReplaceAll(transformed, substr, "")
	}
	return transformed
}
