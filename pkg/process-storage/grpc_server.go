package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"net"
	"pirs.io/commons"
	"pirs.io/process-storage/adapter"
	"pirs.io/process-storage/config"
	pb "pirs.io/process-storage/grpc"
)

var (
	log = commons.GetLoggerFor("processStorageGrpc")
)

type processStorage struct {
	pb.UnimplementedStorageServer
}

func (p *processStorage) UploadProcess(stream pb.Storage_UploadProcessServer) error {
	storageAdapter, err := adapter.MakeStorageAdapter(stream.Context(), config.GetContext().AppConfig.StorageProvider)
	if err != nil {
		return err
	}
	fileBuffer := make([]byte, 0)
	fileMetadata := &pb.ProcessMetadata{}
	chunkNum := 0
	for {
		in, err := stream.Recv()
		chunkNum += 1
		if err == io.EOF {
			log.Info().Msg("EOF! File upload done")

			stream.Send(&pb.ProcessUploadResponse{Status: pb.UploadStatus_SUCCESS})
			break
		}
		if err != nil {
			log.Error().Msg(err.Error())
			stream.Send(&pb.ProcessUploadResponse{Status: pb.UploadStatus_FAILED})
			return err
		}
		if chunkNum == 1 {
			fileMetadata = in.Data.GetMetadata()
			log.Info().Msgf("Got first chunk containing file metadata: %s", fileMetadata.String())
			err = stream.Send(&pb.ProcessUploadResponse{Status: pb.UploadStatus_IN_PROGRESS})
			if err != nil {
				return err
			}
		} else {
			log.Info().Msgf("Got file chunk %s", in.GetData().GetChunk())
			err = stream.Send(&pb.ProcessUploadResponse{Status: pb.UploadStatus_IN_PROGRESS})
			fileBuffer = append(fileBuffer, in.GetData().GetChunk()...)
		}
	}
	err = storageAdapter.SaveProcess(fileMetadata, fileBuffer)
	if err != nil {
		return err
	}
	return nil
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
