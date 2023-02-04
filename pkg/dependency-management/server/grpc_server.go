package server

import (
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/structpb"
	"io"
	"net"
	"pirs.io/commons"
	"pirs.io/commons/structs"
	"pirs.io/dependency-management/config"
	grpcProto "pirs.io/dependency-management/grpc"
	"pirs.io/dependency-management/service/models"
	"strconv"
	"strings"
)

var (
	log = commons.GetLoggerFor("dependencyGrpc")
)

type dependencyServer struct {
	grpcProto.UnimplementedDependencyManagementServer
	appContext *config.ApplicationContext
}

// Detect todo
func (ds *dependencyServer) Detect(stream grpcProto.DependencyManagement_DetectServer) error {
	ctx := stream.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var currentCheckSum string
	var currentChunks []byte
	var currentTotalChunks int
	var chunkCounter int
	separator := ds.appContext.AppConfig.StreamSeparator

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Error().Msgf("cannot receive from the stream: %v", err)
			return err
		}

		// todo authorization (just once I think)

		if req.GetCountAndChecksum() != "" {
			currentChunks = []byte{}
			chunkCounter = 0

			splitInput := strings.Split(req.GetCountAndChecksum(), separator)
			currentTotalChunks, err = strconv.Atoi(splitInput[0])
			if err != nil {
				log.Error().Msgf("cannot parse number to integer: %v", err)
				return err
			}
			currentCheckSum = splitInput[1]
		} else {
			chunkCounter += 1
			if chunkCounter <= currentTotalChunks {
				chunk := req.GetChunkData()
				currentChunks = append(currentChunks, chunk...)
				if chunkCounter == currentTotalChunks {
					data := models.DetectRequestData{
						CheckSum:    currentCheckSum,
						ProcessData: *bytes.NewBuffer(currentChunks),
					}
					response := ds.appContext.DetectionService.Detect(data)
					err = streamDetectResponse(&response, stream)
					if err != nil {
						return err
					}
				}
			} else {
				err = streamDetectResponse(&models.DetectResponseData{
					Status: codes.OutOfRange,
				}, stream)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (ds *dependencyServer) Resolve(stream grpcProto.DependencyManagement_ResolveServer) error {
	panic("not implemented")
}

func streamDetectResponse(response *models.DetectResponseData, stream grpc.ServerStream) error {
	if response.Status == codes.OK {
		err := stream.SendMsg(&grpcProto.DetectResponse{
			Message: "success: " + response.Status.String(),
		})
		if err != nil {
			log.Error().Msgf("cannot send message: %v", err)
			return err
		}

		for _, m := range response.Metadata {
			grpcM, err := structpb.NewStruct(structs.ToMap(m))
			if err != nil {
				return err
			}

			err = stream.SendMsg(&grpcProto.DetectResponse{
				Metadata: grpcM,
			})
			if err != nil {
				return err
			}
		}
	} else {
		err := stream.SendMsg(&grpcProto.DetectResponse{
			Message: "fail: " + response.Status.String(),
		})
		if err != nil {
			log.Error().Msgf("cannot send message: %v", err)
			return err
		}
	}

	return nil
}

// StartGrpc serves GRPC server on given host and port. If it cannot serve, an error is returned.
func StartGrpc(host string, port int, isReflection bool) error {
	flag.Parse()
	log.Info().Msgf("Starting Dependency Management service listening on %s:%d...", host, port)
	lis, networkErr := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if networkErr != nil {
		return networkErr
	}
	grpcServer := grpc.NewServer()

	if isReflection {
		log.Info().Msg("Using GRPC reflection for Dependency Management service.")
		reflection.Register(grpcServer)
	}

	grpcProto.RegisterDependencyManagementServer(grpcServer, &dependencyServer{appContext: config.GetContext()})

	grpcErr := grpcServer.Serve(lis)
	if grpcErr != nil {
		return grpcErr
	}
	return nil
}
