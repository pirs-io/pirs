package server

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"io"
	"net"
	"pirs.io/commons"
	"pirs.io/commons/structs"
	"pirs.io/process/config"
	grpcProto "pirs.io/process/grpc"
	"pirs.io/process/service/models"
	"strconv"
)

var (
	log = commons.GetLoggerFor("processGrpc")
)

type processServer struct {
	grpcProto.UnimplementedProcessServer
	appContext *config.ApplicationContext
}

// Import handles input, that is sent through the opened stream. The input can be metadata or chunks of content. Before
// chunks must arrive metadata of that chunks. It creates goroutine of service.ImportService's method. Channels are used
// to communicate with this goroutine. On success, success message is send to the stream. On failure, error with GRPC
// status code is returned.
func (ps *processServer) Import(stream grpcProto.Process_ImportServer) error {
	response := grpcProto.ImportResponse{}
	createFailureResponse := func(code codes.Code, filename string) error {
		return status.Error(code, "failed to upload file: "+filename)
	}
	createSuccessResponse := func(response *grpcProto.ImportResponse, totalSize uint32, totalFiles int) {
		if totalFiles == 1 {
			response.Message = "successfully uploaded " + strconv.FormatInt(int64(totalFiles), 10) + " file."
		} else {
			response.Message = "successfully uploaded " + strconv.FormatInt(int64(totalFiles), 10) + " files."
		}
		response.TotalSize = totalSize
	}
	defer func(stream grpcProto.Process_ImportServer, response *grpcProto.ImportResponse) {
		err := stream.SendAndClose(response)
		if err != nil {
			log.Error().Msg(status.Errorf(codes.Unavailable, "could not send response and close stream connection: %v", err).Error())
		}
	}(stream, &response)

	ctx := stream.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	requestChan := make(chan models.ImportRequestData)
	responseChan := make(chan models.ImportResponseData)
	defer func() {
		close(requestChan)
		<-responseChan
	}()
	go ps.appContext.ImportService.ImportProcesses(requestChan, responseChan)

	currentFileName := ""
	currentPartialUri := ""
	currentData := bytes.Buffer{}
	currentProcessSize := 0
	totalSize := 0
	totalFiles := 0
	// Receive metadata and bytes. Iteration is one process file.
	for {
		// receive request
		req, err := stream.Recv()
		if err == io.EOF {
			// handle previous received file if needed
			if currentFileName != "" && currentData.Len() != 0 {
				totalSize = totalSize + currentProcessSize
				reqData := models.ImportRequestData{
					Ctx:             ctx,
					PartialUri:      currentPartialUri,
					ProcessFileName: currentFileName,
					ProcessData:     currentData,
					ProcessSize:     currentProcessSize,
					IsLast:          true,
				}
				// handle request
				requestChan <- reqData

				// handle response
				state := (<-responseChan).Status
				if state != codes.OK {
					return createFailureResponse(state, currentFileName)
				}
			}
			createSuccessResponse(&response, uint32(totalSize), totalFiles)
			return nil
		} else if err != nil {
			log.Error().Msgf("cannot receive from the stream: %v", err)
			if req.GetFileInfo() != nil {
				return createFailureResponse(codes.Unavailable, req.GetFileInfo().FileName)
			} else {
				return createFailureResponse(codes.Unavailable, currentFileName)
			}
		}

		// authorization (just once I think)
		// todo

		if req.GetFileInfo() != nil {
			// handle previous received file if needed
			if currentFileName != "" && currentData.Len() != 0 {
				reqData := models.ImportRequestData{
					Ctx:             ctx,
					PartialUri:      currentPartialUri,
					ProcessFileName: currentFileName,
					ProcessData:     currentData,
					ProcessSize:     currentProcessSize,
					IsLast:          false,
				}
				// handle request
				requestChan <- reqData

				// handle response
				state := (<-responseChan).Status
				if state != codes.OK {
					return createFailureResponse(state, currentFileName)
				}
			}

			// new file - reinitialize variables
			totalFiles = totalFiles + 1
			totalSize = totalSize + currentProcessSize
			currentData = bytes.Buffer{}
			currentProcessSize = 0

			// get request metadata
			currentFileName = req.GetFileInfo().GetFileName()
			currentPartialUri = req.GetPartialUri()
		} else {
			chunk := req.GetChunkData()
			size := len(chunk)

			currentProcessSize += size
			if currentProcessSize > ps.appContext.AppConfig.UploadFileMaxSize*1024 {
				err = errors.New(status.Errorf(codes.ResourceExhausted, "file exceeds max size: %d kB",
					ps.appContext.AppConfig.UploadFileMaxSize).Error())
				log.Error().Msg(err.Error())
				return createFailureResponse(codes.ResourceExhausted, currentFileName)
			}
			_, err = currentData.Write(chunk)
			if err != nil {
				log.Error().Msg(status.Errorf(codes.Internal, "cannot write chunk data: %v", err).Error())
				return createFailureResponse(codes.Internal, currentFileName)
			}
		}
	}
}

// Download handles request to download process or package metadata. It streams the response. First it sends success or fail
// message and then metadata one by one.
func (ps *processServer) Download(req *grpcProto.DownloadRequest, stream grpcProto.Process_DownloadServer) error {
	// authorization
	// todo

	// main logic
	ctx := stream.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	reqData := extractDownloadRequest(req, ctx)
	response := ps.appContext.DownloadService.DownloadProcesses(reqData)

	return streamDownloadResponse(response, stream)
}

func extractDownloadRequest(req *grpcProto.DownloadRequest, ctx context.Context) *models.DownloadRequestData {
	return &models.DownloadRequestData{
		Ctx:       ctx,
		TargetUri: req.TargetUri,
		IsPackage: req.IsPackage,
	}
}

func streamDownloadResponse(response *models.DownloadResponseData, stream grpc.ServerStream) error {
	// handle response
	if response.Status == codes.OK {
		err := stream.SendMsg(&grpcProto.DownloadResponse{
			Message: "success: " + response.Status.String(),
		})
		if err != nil {
			return err
		}
	} else {
		err := stream.SendMsg(&grpcProto.DownloadResponse{
			Message: "fail: " + response.Status.String(),
		})
		if err != nil {
			return err
		}
		return nil
	}

	// send all metadata
	for _, m := range response.Metadata {
		grpcM, err := structpb.NewStruct(structs.ToMap(m))
		if err != nil {
			return err
		}

		err = stream.SendMsg(&grpcProto.DownloadResponse{
			Metadata: grpcM,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// StartGrpc serves GRPC server on given host and port. If it cannot serve, an error is returned.
func StartGrpc(host string, port int, isReflection bool) error {
	flag.Parse()
	log.Info().Msgf("Starting Process service listening on %s:%d...", host, port)
	lis, networkErr := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if networkErr != nil {
		return networkErr
	}
	grpcServer := grpc.NewServer()

	if isReflection {
		log.Info().Msg("Using GRPC reflection for Process service.")
		reflection.Register(grpcServer)
	}

	grpcProto.RegisterProcessServer(grpcServer, &processServer{appContext: config.GetContext()})

	grpcErr := grpcServer.Serve(lis)
	if grpcErr != nil {
		return grpcErr
	}
	return nil
}
