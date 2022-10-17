package grpc

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"pirs.io/commons"
	"pirs.io/tracker/config"
)

var (
	log = commons.GetLoggerFor("trackerGrpc")
)

type trackerServer struct {
	UnimplementedTrackerServer
	appContext *config.ApplicationContext
}

func (c *trackerServer) RegisterNewPackage(ctx context.Context, packageInfo *PackageInfo) (*PackageRegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterNewPackage not implemented")
}

func (c *trackerServer) FindPackageLocation(ctx context.Context, in *LocationRequest) (*PackageLocation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindPackageLocation not implemented")
}

func (c *trackerServer) RegisterTrackerInstance(ctx context.Context, in *TrackerInfo) (*InstanceRegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterTrackerInstance not implemented")
}

func (c *trackerServer) GetAllRegisteredInstances(empty *emptypb.Empty, stream Tracker_GetAllRegisteredInstancesServer) error {
	instances, err := c.appContext.InstanceRegistrationService.GetAllRegisteredInstances()
	if err != nil {
		return err
	}
	for _ = range instances {
		stream.Send(nil)
	}
	return nil
}

func StartGrpc(grpcPort int) error {
	flag.Parse()
	lis, networkErr := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", grpcPort))
	if networkErr != nil {
		return networkErr
	}
	grpcServer := grpc.NewServer()

	RegisterTrackerServer(grpcServer, &trackerServer{appContext: config.GetContext()})

	log.Info().Msgf("Running gRPC server on port: %s", grpcPort)
	grpcErr := grpcServer.Serve(lis)
	if grpcErr != nil {
		return grpcErr
	}
	return nil
}
