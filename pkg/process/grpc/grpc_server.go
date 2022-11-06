package grpc

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
	"pirs.io/process/service/models"
)

var (
	log = commons.GetLoggerFor("processGrpc")
)

type processServer struct {
	UnimplementedProcessServer
	appContext *config.ApplicationContext
}

func (ps *processServer) ImportProcess(stream Process_ImportProcessServer) error {
	// receive request
	req, err := stream.Recv()
	if err != nil {
		log.Error().Msg(status.Errorf(codes.Unknown, "cannot receive process info").Error())
		return err
	}
	filename := req.GetFileInfo().GetFileName()
	// authorization
	// todo mocks
	//userRoles := ctx.Value("ROLES").(string)
	//authorize := mocks.CheckAuthorization(userRoles, []string{service.IMPORT_PROCESS_ROLE})
	//if !authorize {
	//	return nil, errors.New("could not authorize with roles: " + userRoles)
	//}
	// extract req
	processData := bytes.Buffer{}
	processSize := 0
	for {
		req, err = stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Error().Msg(status.Errorf(codes.Unknown, "cannot receive chunk data: %v", err).Error())
			return err
		}
		chunk := req.GetChunkData()
		size := len(chunk)

		processSize += size
		if processSize > ps.appContext.AppConfig.UploadFileMaxSize*1024 {
			err = errors.New(status.Errorf(codes.ResourceExhausted, "file exceeds max size: %d kB",
				ps.appContext.AppConfig.UploadFileMaxSize).Error())
			log.Error().Msg(err.Error())
			return err
		}
		_, err = processData.Write(chunk)
		if err != nil {
			log.Error().Msg(status.Errorf(codes.Internal, "cannot write chunk data: %v", err).Error())
			return err
		}
	}
	ctx := stream.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	reqData := models.ImportProcessRequestData{
		Ctx:             ctx,
		ProcessFileName: filename,
		ProcessData:     processData,
		ProcessSize:     processSize,
	}
	// handle request
	responseData, err := ps.appContext.ImportService.ImportProcess(&reqData)
	if err != nil {
		return err
	}
	// handle response
	importProcessResponse := ImportProcessResponse{}
	if responseData.Status == codes.OK {
		importProcessResponse.Message = "successfully uploaded file: " + filename
		importProcessResponse.TotalSize = int32(processSize)
	} else {
		importProcessResponse.Message = "failed to upload file: " + filename
		importProcessResponse.TotalSize = 0
	}
	err = stream.SendAndClose(&importProcessResponse)
	if err != nil {
		return err
	}
	return nil
}

func (c *processServer) ImportPackage(ctx context.Context, req *ImportPackageRequest) (*ImportPackageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ImportPackage not implemented")
}

func (c *processServer) RemoveProcess(ctx context.Context, req *RemoveProcessRequest) (*RemoveProcessResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveProcess not implemented")
}

func (c *processServer) DownloadProcess(ctx context.Context, req *DownloadProcessRequest) (*DownloadProcessResponse, error) {
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

	RegisterProcessServer(grpcServer, &processServer{appContext: config.GetContext()})

	grpcErr := grpcServer.Serve(lis)
	if grpcErr != nil {
		return grpcErr
	}
	return nil
}
