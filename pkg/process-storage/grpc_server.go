package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"net"
	"pirs.io/common"
	"pirs.io/process-storage/adapter"
	"pirs.io/process-storage/config"
	pb "pirs.io/process-storage/grpc"
)

var (
	log = common.GetLoggerFor("processStorageGrpc")
)

type processStorage struct {
	pb.UnimplementedStorageServer
}

func (p *processStorage) UploadProcess(stream pb.Storage_UploadProcessServer) error {
	_, err := adapter.MakeStorageAdapter(stream.Context(), config.GetContext().AppConfig.StorageProvider)
	if err != nil {
		return err
	}
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			stream.Send(&pb.ProcessUploadResponse{Status: pb.UploadStatus_SUCCESS})
			return nil
		}
		if err != nil {
			stream.Send(&pb.ProcessUploadResponse{Status: pb.UploadStatus_FAILED})
			return nil
		}
		log.Info().Msg(in.String())
		err = stream.Send(&pb.ProcessUploadResponse{Status: pb.UploadStatus_IN_PROGRESS})
		if err != nil {
			return err
		}
	}
}

func (p *processStorage) DownloadProcess(req *pb.ProcessDownloadRequest, stream pb.Storage_DownloadProcessServer) error {
	return nil
}

func StartGrpc(grpcPort int) error {
	flag.Parse()
	lis, networkErr := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", grpcPort))
	if networkErr != nil {
		return networkErr
	}
	grpcServer := grpc.NewServer()

	pb.RegisterStorageServer(grpcServer, &processStorage{})

	log.Info().Msgf("Running gRPC server on port: %s", grpcPort)
	grpcErr := grpcServer.Serve(lis)
	if grpcErr != nil {
		return grpcErr
	}
	return nil
}
