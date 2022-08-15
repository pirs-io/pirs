package grpc

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"pirs.io/pirs/common"
	"pirs.io/pirs/common/trackerProto"
)

var (
	log = common.GetLoggerFor("trackerGrpc")
)

type TrackerServer struct {
	trackerProto.UnimplementedTrackerServer
}

func (c *TrackerServer) RegisterNewPackage(ctx context.Context, packageInfo *trackerProto.PackageInfo) (*trackerProto.RegisterResponse, error) {
	log.Info().Msg("Registering new package")
	return &trackerProto.RegisterResponse{}, nil
}

func StartGrpc(grpcPort int) error {
	flag.Parse()
	lis, networkErr := net.Listen("tcp", fmt.Sprintf("localhost:%d", grpcPort))
	if networkErr != nil {
		return networkErr
	}
	grpcServer := grpc.NewServer()

	trackerProto.RegisterTrackerServer(grpcServer, &TrackerServer{})

	grpcErr := grpcServer.Serve(lis)
	if grpcErr != nil {
		return grpcErr
	}
	return nil
}
