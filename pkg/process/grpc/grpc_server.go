package grpc

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"net"
	"pirs.io/commons"
	"pirs.io/process/config"
	"pirs.io/process/mock"
	"pirs.io/process/service"
)

var (
	log = commons.GetLoggerFor("processGrpc")
)

type processServer struct {
	UnimplementedProcessServer
	appContext *config.ApplicationContext
}

func (ps *processServer) ImportProcess(ctx context.Context, req *ImportProcessRequest) (*ImportProcessResponse, error) {
	// authorization
	// todo mock
	userRoles := ctx.Value("ROLES").(string)
	authorize := mock.CheckAuthorization(userRoles, []string{service.IMPORT_PROCESS_ROLE})
	if !authorize {
		return nil, errors.New("could not authorize with roles: " + userRoles)
	}
	// extract req
	// todo mock
	myMockFilePath := "awd.txt"
	reqFilePtr, err := mock.FindOrCreateFile(myMockFilePath)
	if err != nil {
		return nil, errors.New("could not find or create file: " + myMockFilePath)
	}
	reqData := service.ImportProcessRequestData{
		ProcessFile: reqFilePtr,
	}
	// handle request
	response, err := ps.appContext.ImportService.ImportProcess(&reqData)
	if err != nil {
		return nil, err
	}
	// handle response
	importProcessResponse := ImportProcessResponse{}
	if response.Status == 0 {
		log.Info().Msg("Process was successfully imported.")
		// initialize response based on state
	}
	return &importProcessResponse, nil
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
