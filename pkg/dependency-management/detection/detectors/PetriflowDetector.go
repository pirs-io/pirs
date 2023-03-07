package detectors

import (
	"bytes"
	"github.com/antchfx/xmlquery"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"index/suffixarray"
	"pirs.io/commons"
	"pirs.io/commons/db/mongo"
	"pirs.io/commons/domain"
	"pirs.io/commons/enums"
	"pirs.io/dependency-management/detection/models"
	"strings"
	"sync"
)

const (
	actionXpath      = "//action"
	protocolPrefix   = "pirs://"
	protocolBoundary = "\""
)

var (
	log = commons.GetLoggerFor("petriflowDetector")
)

// A PetriflowDetector represents structure for dependency detection of process type enums.Petriflow. It contains field next,
// which is a pointer on the next models.Detector within chain of responsibility pattern and repository for metadata.
type PetriflowDetector struct {
	repository mongo.MetadataRepository
	next       models.Detector
}

// NewPetriflowDetector return pointer of instance of PetriflowDetector. It contains initialized metadata repository.
func NewPetriflowDetector(repo mongo.MetadataRepository) *PetriflowDetector {
	return &PetriflowDetector{
		repository: repo,
	}
}

// An actionContext is a wrapper of action inner text. Contains also suffixarray.Index of action body and parentURI of
// parent metadata, which is going to be handled.
type actionContext struct {
	body        string
	indexedBody *suffixarray.Index
	parentURI   string
	variables   map[string]string
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
	isDone := make(chan bool)

	// receive metadata from actionHandler goroutines
	go func() {
		for received := range responseChan {
			dependencies = append(dependencies, received)
		}
		isDone <- true
	}()

	// handle actions
	for _, actionNode := range nodes {
		action := actionContext{body: actionNode.InnerText()}
		wg.Add(1)
		go pd.handleAction(&wg, action, responseChan)
	}
	wg.Wait()
	close(responseChan)
	// must wait to handle last loop in anonym goroutine
	<-isDone
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

// handleAction handles action in actionContext. todo
func (pd *PetriflowDetector) handleAction(wg *sync.WaitGroup, action actionContext, responseChan chan<- domain.Metadata) {
	defer wg.Done()

	isDone := make(chan bool)
	var wgForSearch sync.WaitGroup
	responseChanForSearch := make(chan string)

	go func() {
		for found := range responseChanForSearch {
			dependency, err := pd.repository.FindByURI(context.Background(), found)
			if err != nil {
				log.Error().Msgf("an error occurred while searching for %s in repository: %v", found, err)
			} else if dependency.ID != primitive.NilObjectID {
				responseChan <- dependency
			}
		}
		isDone <- true
	}()

	action.indexedBody = suffixarray.New([]byte(action.body))
	wgForSearch.Add(2)
	go pd.searchForProtocols(&wgForSearch, action, responseChanForSearch)
	go pd.searchForNonProtocols(&wgForSearch, action, responseChanForSearch)
	wgForSearch.Wait()
	close(responseChanForSearch)
	// must wait to handle last loop in anonym goroutine
	<-isDone
}

func (pd *PetriflowDetector) searchForProtocols(wg *sync.WaitGroup, action actionContext, responseChan chan<- string) {
	defer wg.Done()

	foundIdxs := action.indexedBody.Lookup([]byte(protocolPrefix), -1)
	for _, startingIdx := range foundIdxs {
		endingIdx := strings.Index(action.body[startingIdx:], protocolBoundary)
		uriWithPrefix := action.body[startingIdx : endingIdx+startingIdx]
		trimmed := uriWithPrefix[len(protocolPrefix):]
		responseChan <- trimmed
	}
}

func (pd *PetriflowDetector) searchForNonProtocols(wg *sync.WaitGroup, action actionContext, responseChan chan<- string) {
	defer wg.Done()
	// todo
}
