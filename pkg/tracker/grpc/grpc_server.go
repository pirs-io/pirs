package grpc

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"pirs.io/pirs/common"
	"pirs.io/pirs/common/trackerProto"
	"pirs.io/pirs/tracker/config"
)

var (
	log = common.GetLoggerFor("trackerGrpc")
)

type TrackerServer struct {
	trackerProto.UnimplementedTrackerServer
	appContext *config.ApplicationContext
}

func (c *TrackerServer) RegisterNewPackage(ctx context.Context, packageInfo *trackerProto.PackageInfo) (*trackerProto.PackageRegisterResponse, error) {
	c.appContext.LocationService.RegisterPackage(packageInfo)
	return nil, status.Errorf(codes.Unimplemented, "method FindPackageLocation not implemented")
}

func (c *TrackerServer) FindPackageLocation(ctx context.Context, in *trackerProto.LocationRequest) (*trackerProto.PackageLocation, error) {
	log.Info().Msgf("Finding package")
	return &trackerProto.PackageLocation{}, nil
}

func (c *TrackerServer) RegisterTrackerInstance(ctx context.Context, in *trackerProto.TrackerInfo) (*trackerProto.InstanceRegisterResponse, error) {
	return c.appContext.InstanceRegistrationService.RegisterInstance(in)
}

func StartGrpc(grpcPort int) error {
	flag.Parse()
	lis, networkErr := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", grpcPort))
	if networkErr != nil {
		return networkErr
	}
	grpcServer := grpc.NewServer()

	trackerProto.RegisterTrackerServer(grpcServer, &TrackerServer{appContext: config.GetContext()})

	log.Info().Msgf("Running gRPC server on port: %s", grpcPort)
	grpcErr := grpcServer.Serve(lis)
	if grpcErr != nil {
		return grpcErr
	}
	return nil
}
