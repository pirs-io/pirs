package detectors

import (
	"bytes"
	"github.com/antchfx/xmlquery"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"pirs.io/commons"
	"pirs.io/commons/enums"
	"pirs.io/dependency-management/detection/models"
	"pirs.io/process/domain"
	"sync"
)

const (
	actionXpath = "//action"
)

var (
	log = commons.GetLoggerFor("petriflowDetector")
)

// A PetriflowDetector represents structure for dependency detection of process type enums.Petriflow. It contains field next,
// which is a pointer on the next models.Detector within chain of responsibility pattern.
type PetriflowDetector struct {
	next models.Detector
}

// Detect is a handler for dependency detection out of process data. The input is enums.ProcessType and bytes.Buffer.
// If enums.ProcessType is not the type of this handler, function ExecuteNextIfPresent is called. It searches for action
// tags and sends tags one by one into handleAction goroutine. These goroutines send dependencies one by one. Every
// dependency is collected and returned in the end.
func (pd *PetriflowDetector) Detect(processType enums.ProcessType, data bytes.Buffer) []domain.Metadata {
	if !pd.IsProcessTypeEqual(processType) {
		return pd.ExecuteNextIfPresent(processType, data)
	}

	// find action nodes by xpath
	doc, err := xmlquery.Parse(bytes.NewReader(data.Bytes()))
	if err != nil {
		log.Error().Msg(status.Errorf(codes.Internal, "could not prepare xml data to parse: %v", err).Error())
		return []domain.Metadata{}
	}
	nodes := xmlquery.Find(doc, actionXpath)
	if nodes == nil || len(nodes) == 0 {
		return []domain.Metadata{}
	}

	// create channels
	var wg sync.WaitGroup
	var dependencies []domain.Metadata
	responseChan := make(chan domain.Metadata)

	// receive metadata from actionHandler goroutines
	go func() {
		for received := range responseChan {
			dependencies = append(dependencies, received)
		}
	}()

	// handle actions
	for _, actionNode := range nodes {
		wg.Add(1)
		go pd.handleAction(&wg, actionNode.InnerText(), responseChan)
	}
	wg.Wait()
	close(responseChan)

	return dependencies
}

// SetNext sets next models.Detector within chain of responsibility.
func (pd *PetriflowDetector) SetNext(detector models.Detector) {
	pd.next = detector
}

// ExecuteNextIfPresent executes next models.Detector if exists.
func (pd *PetriflowDetector) ExecuteNextIfPresent(processType enums.ProcessType, data bytes.Buffer) []domain.Metadata {
	if pd.next != nil {
		return pd.next.Detect(processType, data)
	}
	return []domain.Metadata{}
}

// IsProcessTypeEqual checks if enums.ProcessType is equal to handler type
func (pd *PetriflowDetector) IsProcessTypeEqual(toCheck enums.ProcessType) bool {
	return enums.Petriflow == toCheck
}

// handleAction todo
func (pd *PetriflowDetector) handleAction(wg *sync.WaitGroup, body string, responseChan chan<- domain.Metadata) {
	defer wg.Done()

	// todo implement detection
	responseChan <- *domain.NewMetadata()
}
