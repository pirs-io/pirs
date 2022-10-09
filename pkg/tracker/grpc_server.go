package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
	"pirs.io/common"
	"pirs.io/tracker/config"
	pb "pirs.io/tracker/grpc"
)

var (
	grpcLog = common.GetLoggerFor("trackerGrpc")
)

type trackerServer struct {
	pb.UnimplementedTrackerServer
	appContext *config.ApplicationContext
}

func (c *trackerServer) RegisterNewPackage(ctx context.Context, packageInfo *pb.PackageInfo) (*pb.PackageRegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterNewPackage not implemented")
}

func (c *trackerServer) FindPackageLocation(ctx context.Context, in *pb.LocationRequest) (*pb.PackageLocation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindPackageLocation not implemented")
}

func (c *trackerServer) RegisterTrackerInstance(ctx context.Context, in *pb.TrackerInfo) (*pb.InstanceRegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterTrackerInstance not implemented")
}

func (c *trackerServer) GetAllRegisteredInstances(empty *emptypb.Empty, stream pb.Tracker_GetAllRegisteredInstancesServer) error {
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

	pb.RegisterTrackerServer(grpcServer, &trackerServer{appContext: config.GetContext()})

	grpcLog.Info().Msgf("Running gRPC server on port: %s", grpcPort)
	grpcErr := grpcServer.Serve(lis)
	if grpcErr != nil {
		return grpcErr
	}
	return nil
}
