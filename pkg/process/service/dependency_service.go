package service

import (
	"flag"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	mygrpc "pirs.io/process/grpc"
	"pirs.io/process/service/models"
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

// Detect todo
func (ds *DependencyService) Detect(reqCtx context.Context, forResource <-chan []byte, forResponse chan<- models.ResponseAdapter) {
	// todo establish connection
	// todo send data
	// todo receive data
}
