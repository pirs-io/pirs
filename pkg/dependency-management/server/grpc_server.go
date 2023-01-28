package server

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"pirs.io/commons"
	"pirs.io/dependency-management/config"
	grpcProto "pirs.io/dependency-management/grpc"
)

var (
	log = commons.GetLoggerFor("dependencyGrpc")
)

type dependencyServer struct {
	grpcProto.UnimplementedDependencyManagementServer
	appContext *config.ApplicationContext
}

// Detect todo
func (ds *dependencyServer) Detect(stream grpcProto.DependencyManagement_DetectServer) error {
	panic("not implemented")
}

// StartGrpc serves GRPC server on given host and port. If it cannot serve, an error is returned.
func StartGrpc(host string, port int, isReflection bool) error {
	flag.Parse()
	log.Info().Msgf("Starting Dependency Management service listening on %s:%d...", host, port)
	lis, networkErr := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if networkErr != nil {
		return networkErr
	}
	grpcServer := grpc.NewServer()

	if isReflection {
		log.Info().Msg("Using GRPC reflection for Dependency Management service.")
		reflection.Register(grpcServer)
	}

	grpcProto.RegisterDependencyManagementServer(grpcServer, &dependencyServer{appContext: config.GetContext()})

	grpcErr := grpcServer.Serve(lis)
	if grpcErr != nil {
		return grpcErr
	}
	return nil
}
