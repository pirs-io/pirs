package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"math"
	"pirs.io/commons"
	"pirs.io/commons/enums"
	"pirs.io/process/domain"
	mygrpc "pirs.io/process/grpc"
	"pirs.io/process/service/models"
	"strconv"
	"strings"
)

// A DependencyService detects and resolves dependencies. Host and Port describes dependency-management service instance.
// Separator is used as separator in streaming (separating total chunks and checksum). ChunkSize is used to calculate number
// of chunks and to devide array of bytes into chunks
type DependencyService struct {
	Host      string
	Port      string
	Separator string
	ChunkSize int
	client    mygrpc.DependencyManagementClient
}

var (
	logDs = commons.GetLoggerFor("DependencyService")
)

// NewDependencyService creates instance with params. It also creates client for dependency-management service instance.
func NewDependencyService(hostname string, port string, sep string, chunkSize int) (*DependencyService, error) {
	service := &DependencyService{
		Host:      hostname,
		Port:      port,
		Separator: sep,
		ChunkSize: chunkSize,
	}
	client, err := service.createClient()
	if err != nil {
		return service, nil
	}
	service.client = client
	return service, nil
}

func (ds *DependencyService) createClient() (mygrpc.DependencyManagementClient, error) {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(ds.Host+":"+ds.Port, opts...)
	if err != nil {
		logDs.Error().Msgf("cannot dial dependency-management server: %v", err)
		return nil, err
	}

	return mygrpc.NewDependencyManagementClient(conn), nil
}

// Detect is a client service for client.Detect endpoint. It is run as goroutine by ImportService. First it establishes
// bi-directional stream connection with dependency-management service. It accepts file data as an array of bytes via
// forResource channel. Then it handles that data and sends to dependency-management service. Then metadata are received
// and sent to parent ImportService via forResponse channel. If any error occurs during opened stream connection an error
// is sent to ImportService and stream connection is shutdown.
func (ds *DependencyService) Detect(reqCtx context.Context, forResource <-chan models.DetectResourceAdapter, forResponse chan<- models.ResponseAdapter) {
	defer close(forResponse)
	// open stream
	stream, err := ds.establishDetectConnection(reqCtx)
	if err != nil {
		forResponse <- models.ResponseAdapter{Err: err}
		return
	} else {
		forResponse <- models.ResponseAdapter{}
	}

	for resource := range forResource {
		// send request metadata
		rawHash := commons.HashBytesToSHA256(resource.FileData)
		checksum := commons.ConvertBytesToString(rawHash)
		totalChunks := uint64(math.Ceil(float64(len(resource.FileData)) / float64(ds.ChunkSize)))
		err = ds.sendDetectRequest(stream,
			strconv.FormatUint(totalChunks, 10)+ds.Separator+checksum,
			resource.ProcessType,
			nil,
		)
		if err != nil {
			forResponse <- models.ResponseAdapter{Err: err}
			return
		}

		// send chunks
		reader := bytes.NewReader(resource.FileData)
		buffer := make([]byte, ds.ChunkSize)
		for {
			n, err := reader.Read(buffer)
			if err == io.EOF {
				break
			}
			if err != nil {
				logDs.Error().Msgf("cannot read chunk to buffer: ", err)
				forResponse <- models.ResponseAdapter{Err: err}
				return
			}

			err = ds.sendDetectRequest(stream, "", -1, buffer[:n])
			if err != nil {
				forResponse <- models.ResponseAdapter{Err: err}
				return
			}
		}
		// receive status response
		response, err := stream.Recv()
		if err != nil {
			logDs.Error().Msgf("cannot receive from the stream: ", err, stream.RecvMsg(nil))
			forResponse <- models.ResponseAdapter{Err: err}
			return
		}
		if strings.Contains(response.Message, "fail") {
			logDs.Error().Msgf("dependencies failed to detect: ", err, stream.RecvMsg(nil))
			forResponse <- models.ResponseAdapter{Err: errors.New(response.Message)}
			continue
		}
		// receive metadata
		for {
			response, err = stream.Recv()
			if err != nil {
				logDs.Error().Msgf("cannot receive from the stream: ", err, stream.RecvMsg(nil))
				forResponse <- models.ResponseAdapter{Err: err}
				return
			}

			metadataFromResponse := domain.Metadata{}
			jsonString, _ := json.Marshal(response.Metadata)
			err = json.Unmarshal(jsonString, &metadataFromResponse)
			if err != nil {
				logDs.Error().Msgf("cannot handle received metadata: %v", err)
				forResponse <- models.ResponseAdapter{Err: err}
				return
			}
			forResponse <- models.ResponseAdapter{Metadata: metadataFromResponse}
			if metadataFromResponse.ID == primitive.NilObjectID {
				break
			}
		}
	}
	// shutdown
	_ = ds.destroyConnection(stream)
}

func (ds *DependencyService) establishDetectConnection(ctx context.Context) (mygrpc.DependencyManagement_DetectClient, error) {
	var err error

	if ds.client == nil {
		err = errors.New("Dependency-Management client is not initialized")
		logDs.Error().Msg(err.Error())
		return nil, err
	}

	stream, err := ds.client.Detect(ctx)
	if err != nil {
		logDs.Error().Msgf("could not establish stream connection: %v", err)
		return nil, err
	}
	return stream, nil
}

func (ds *DependencyService) sendDetectRequest(stream mygrpc.DependencyManagement_DetectClient, countAndChecksum string, processType enums.ProcessType, chunk []byte) error {
	var err error
	if countAndChecksum != "" {
		err = stream.Send(&mygrpc.DetectRequest{
			ProcessType: int32(processType.Int()),
			Data: &mygrpc.DetectRequest_CountAndChecksum{
				CountAndChecksum: countAndChecksum,
			},
		})
	} else {
		err = stream.Send(&mygrpc.DetectRequest{
			Data: &mygrpc.DetectRequest_ChunkData{
				ChunkData: chunk,
			},
		})
	}
	if err != nil {
		logDs.Error().Msg(err.Error())
		return err
	}
	return nil
}

func (ds *DependencyService) destroyConnection(stream mygrpc.DependencyManagement_DetectClient) error {
	err := stream.CloseSend()
	if err != nil {
		log.Error().Msgf("cannot close stream connection: %v", err)
		return err
	}
	return nil
}
