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
	"time"
)

var (
	log = commons.GetLoggerFor("processGrpc")
)

type processServer struct {
	grpcProto.UnimplementedProcessServer
	appContext *config.ApplicationContext
}

// Import handles request to import process file along with metadata. It authorizes user, validates request,
// extracts metadata, generates response. If success, a success message is sent to the client. Otherwise, a fail message
// is sent.
func (ps *processServer) Import(stream grpcProto.Process_ImportServer) error {
	importProcessResponse := grpcProto.ImportResponse{}
	createFailureResponse := func(response *grpcProto.ImportResponse, filename string) {
		importProcessResponse.Message = "failed to upload file: " + filename
		importProcessResponse.TotalSize = 0
	}
	createSuccessResponse := func(response *grpcProto.ImportResponse, filename string, filesize uint32) {
		importProcessResponse.Message = "successfully uploaded file: " + filename
		importProcessResponse.TotalSize = filesize
	}
	defer func(stream grpcProto.Process_ImportServer, response *grpcProto.ImportResponse) {
		err := stream.SendAndClose(response)
		if err != nil {
			log.Error().Msg(status.Errorf(codes.Unavailable, "could not send response and close stream connection: %v", err).Error())
		}
	}(stream, &importProcessResponse)

	// receive request
	req, err := stream.Recv()
	if err != nil {
		log.Error().Msg(status.Errorf(codes.Unknown, "cannot receive process info").Error())
		createFailureResponse(&importProcessResponse, "unknown")
		return err
	}
	// authorization
	// todo
	// extract req
	filename := req.GetFileInfo().GetFileName()
	partialUri := req.GetPartialUri()
	processData := bytes.Buffer{}
	processSize := 0
	for {
		req, err = stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Error().Msg(status.Errorf(codes.Unknown, "cannot receive chunk data: %v", err).Error())
			createFailureResponse(&importProcessResponse, filename)
			return err
		}
		chunk := req.GetChunkData()
		size := len(chunk)

		processSize += size
		if processSize > ps.appContext.AppConfig.UploadFileMaxSize*1024 {
			err = errors.New(status.Errorf(codes.ResourceExhausted, "file exceeds max size: %d kB",
				ps.appContext.AppConfig.UploadFileMaxSize).Error())
			log.Error().Msg(err.Error())
			createFailureResponse(&importProcessResponse, filename)
			return err
		}
		_, err = processData.Write(chunk)
		if err != nil {
			log.Error().Msg(status.Errorf(codes.Internal, "cannot write chunk data: %v", err).Error())
			createFailureResponse(&importProcessResponse, filename)
			return err
		}
	}
	ctx := stream.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	// todo tmp timeout with cancel
	ctx, cancel := context.WithTimeout(ctx, 10*time.Hour)
	defer cancel()

	reqData := models.ImportProcessRequestData{
		Ctx:             ctx,
		PartialUri:      partialUri,
		ProcessFileName: filename,
		ProcessData:     processData,
		ProcessSize:     processSize,
	}
	// handle request
	responseData := ps.appContext.ImportService.ImportProcess(&reqData)

	// handle response
	if responseData.Status == codes.OK {
		createSuccessResponse(&importProcessResponse, filename, uint32(processSize))
	} else {
		createFailureResponse(&importProcessResponse, filename)
	}
	return nil
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
