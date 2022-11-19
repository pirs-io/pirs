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
	"io"
	"net"
	"pirs.io/commons"
	"pirs.io/process/config"
	grpcProto "pirs.io/process/grpc"
	"pirs.io/process/service/models"
)

var (
	log = commons.GetLoggerFor("processGrpc")
)

type processServer struct {
	grpcProto.UnimplementedProcessServer
	appContext *config.ApplicationContext
}

func (ps *processServer) ImportProcess(stream grpcProto.Process_ImportProcessServer) error {
	importProcessResponse := grpcProto.ImportProcessResponse{}
	createFailureResponse := func(response *grpcProto.ImportProcessResponse, filename string) {
		importProcessResponse.Message = "failed to upload file: " + filename
		importProcessResponse.TotalSize = 0
	}
	createSuccessResponse := func(response *grpcProto.ImportProcessResponse, filename string, filesize uint32) {
		importProcessResponse.Message = "successfully uploaded file: " + filename
		importProcessResponse.TotalSize = filesize
	}
	defer func(stream grpcProto.Process_ImportProcessServer, response *grpcProto.ImportProcessResponse) {
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

func (c *processServer) ImportPackage(ctx context.Context, req *grpcProto.ImportPackageRequest) (*grpcProto.ImportPackageResponse, error) {
	// todo
	return nil, status.Errorf(codes.Unimplemented, "method ImportPackage not implemented")
}

func (c *processServer) RemoveProcess(ctx context.Context, req *grpcProto.RemoveProcessRequest) (*grpcProto.RemoveProcessResponse, error) {
	// todo
	return nil, status.Errorf(codes.Unimplemented, "method RemoveProcess not implemented")
}

func (c *processServer) DownloadProcess(ctx context.Context, req *grpcProto.DownloadProcessRequest) (*grpcProto.DownloadProcessResponse, error) {
	// todo
	return nil, status.Errorf(codes.Unimplemented, "method DownloadProcess not implemented")
}

func StartGrpc(grpcIp string, grpcPort int, isReflection bool) error {
	flag.Parse()
	log.Info().Msgf("Starting Process service listening on %s:%d...", grpcIp, grpcPort)
	lis, networkErr := net.Listen("tcp", fmt.Sprintf("%s:%d", grpcIp, grpcPort))
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
