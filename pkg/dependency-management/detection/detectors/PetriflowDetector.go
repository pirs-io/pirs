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
	"regexp"
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
	apiToDetect map[string][]string
	repository  mongo.Repository
	next        models.Detector
}

// NewPetriflowDetector return pointer of instance of PetriflowDetector. It contains initialized metadata repository.
func NewPetriflowDetector(repo mongo.Repository, api map[string][]string) *PetriflowDetector {
	return &PetriflowDetector{
		apiToDetect: api,
		repository:  repo,
	}
}

// An actionContext is a wrapper of action inner text. Contains also suffixarray.Index of action body and parentURI of
// parent metadata, which is going to be handled.
type actionContext struct {
	body             string
	indexedBody      *suffixarray.Index
	parentProjectURI string
	variables        map[string]string
}

// Detect is a handler for dependency detection out of process data. The input is enums.ProcessType and bytes.Buffer.
// If enums.ProcessType is not the type of this handler, function ExecuteNextIfPresent is called. It searches for action
// tags and sends tags one by one into handleAction goroutine. These goroutines send dependencies one by one. Every
// dependency is collected and returned in the end. todo
func (pd *PetriflowDetector) Detect(req models.DetectRequestData) []domain.Metadata {
	if !pd.IsProcessTypeEqual(req.ProcessType) {
		return pd.ExecuteNextIfPresent(req)
	}

	// find action nodes by xpath
	doc, err := xmlquery.Parse(bytes.NewReader(req.ProcessData.Bytes()))
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
	receivedURIs := map[string]struct{}{}

	// receive metadata from actionHandler goroutines
	go func() {
		for received := range responseChan {
			if _, ok := receivedURIs[received.URI]; !ok {
				receivedURIs[received.URI] = struct{}{}
				dependencies = append(dependencies, received)
			}
		}
		isDone <- true
	}()

	// handle actions
	for _, actionNode := range nodes {
		action := actionContext{body: actionNode.InnerText()}
		action.parentProjectURI = req.ProjectUri
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
func (pd *PetriflowDetector) ExecuteNextIfPresent(req models.DetectRequestData) []domain.Metadata {
	if pd.next != nil {
		return pd.next.Detect(req)
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
			if !strings.Contains(found, ":") {
				latestM, err := pd.repository.FindNewestByURI(context.Background(), found)
				if err == nil && latestM.ID != primitive.NilObjectID {
					responseChan <- latestM
				}
			} else {
				dependency, err := pd.repository.FindByURI(context.Background(), found)
				if err != nil {
					log.Error().Msgf("an error occurred while searching for %s in repository: %v", found, err)
				} else if dependency.ID != primitive.NilObjectID {
					responseChan <- dependency
				}
			}
		}
		isDone <- true
	}()

	action.indexedBody = suffixarray.New([]byte(action.body))
	// if 1 searching goroutine is added, delta for wgForSearch must be incremented by 1
	wgForSearch.Add(2)
	go pd.searchForProtocols(&wgForSearch, action, responseChanForSearch)
	go pd.searchForApiFunctions(&wgForSearch, action, responseChanForSearch)
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

// searchForApiFunctions todo
func (pd *PetriflowDetector) searchForApiFunctions(wg *sync.WaitGroup, action actionContext, responseChan chan<- string) {
	defer wg.Done()

	keys := collectMapKeys(pd.apiToDetect)
	r := regexp.MustCompile(strings.Join(keys, "|"))
	// find any matching api by defined input from CSV
	res := action.indexedBody.FindAllIndex(r, -1)
	sentURIs := map[string]struct{}{}
	for _, pairIdx := range res {
		// determine what is found
		apiFunc := action.body[pairIdx[0]:pairIdx[1]]
		fromDelimiter := pd.apiToDetect[apiFunc][0]
		untilDelimiter := pd.apiToDetect[apiFunc][1]

		// use from and until to remove unnecessary substrings
		fromIdx := strings.Index(action.body[pairIdx[1]:], fromDelimiter)
		if fromIdx == -1 {
			continue
		}
		halfTrimmed := action.body[pairIdx[1]+fromIdx+len(fromDelimiter):]
		untilIdx := strings.Index(halfTrimmed, untilDelimiter)
		if untilIdx == -1 {
			continue
		}

		// use found indexes to resolve identifier
		identifier := strings.TrimSpace(halfTrimmed[:untilIdx])
		uriToSend := action.parentProjectURI + "." + identifier

		// send built uri
		if _, ok := sentURIs[uriToSend]; !ok {
			sentURIs[uriToSend] = struct{}{}
			responseChan <- uriToSend
		}
	}
}

func collectMapKeys(aMap map[string][]string) []string {
	keys := make([]string, len(aMap))
	i := 0
	for k := range aMap {
		keys[i] = k
		i++
	}
	return keys
}
