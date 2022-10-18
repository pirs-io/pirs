package grpc

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"net"
	"pirs.io/commons"
	"pirs.io/process/config"
)

var (
	log = commons.GetLoggerFor("trackerGrpc")
)

type processServer struct {
	UnimplementedProcessServer
	appContext *config.ApplicationContext
}

func (c *processServer) ImportProcess(ctx context.Context, req *ImportProcessRequest) (*ImportProcessResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ImportProcess not implemented")
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
	log.Info().Msgf("Starting Process service on %s:%d...", grpcIp, grpcPort)
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
