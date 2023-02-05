package service

import (
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	mygrpc "pirs.io/process/grpc"
)

type DependencyService struct {
	Host   string
	Port   string
	client mygrpc.DependencyManagementClient
}

func NewDependencyService(hostname string, port string) (*DependencyService, error) {
	service := &DependencyService{
		Host: hostname,
		Port: port,
	}
	client, err := service.createClient()
	if err != nil {
		return service, nil
	}
	service.client = client
	return service, nil
}

func (ds *DependencyService) createClient() (mygrpc.DependencyManagementClient, error) {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(ds.Host+":"+ds.Port, opts...)
	if err != nil {
		log.Error().Msgf("cannot dial dependency-management server: %v", err)
		return nil, err
	}

	return mygrpc.NewDependencyManagementClient(conn), nil
}
