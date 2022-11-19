package service

import (
	"bytes"
	"flag"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"pirs.io/commons"
	"pirs.io/process/domain"
	mygrpc "pirs.io/process/grpc"
)

var (
	log = commons.GetLoggerFor("StorageService")
)

type StorageService struct {
	Port      string
	Host      string
	ChunkSize int
}

func (ss *StorageService) SaveFile(reqCtx context.Context, m domain.Metadata, fileData []byte) error {
	client, err := ss.createClient()
	if err != nil {
		return err
	}

	reqData := &mygrpc.ProcessFileData_Metadata{
		Metadata: &mygrpc.ProcessMetadata{
			ProcessId: m.URI,
			Filename:  m.FileName,
			Encoding:  0,
			Type:      0,
		},
	}

	stream, err := client.UploadProcess(reqCtx)
	if err != nil {
		log.Error().Msgf("could not establish stream connection: %v", err)
		return err
	}

	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Error().Msgf("client.RouteChat failed: %v", err)
				return
			}
			log.Debug().Msg(string(in.Status))
		}
	}()

	if err := stream.Send(&mygrpc.ProcessUploadRequest{
		Data: &mygrpc.ProcessFileData{
			Data: reqData,
		},
	}); err != nil {
		log.Error().Msg(err.Error())
		return err
	}

	reader := bytes.NewReader(fileData)
	buffer := make([]byte, ss.ChunkSize)
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error().Msgf("cannot read chunk to buffer: ", err)
			return err
		}

		req := &mygrpc.ProcessFileData_Chunk{Chunk: buffer[:n]}

		err = stream.Send(&mygrpc.ProcessUploadRequest{
			Data: &mygrpc.ProcessFileData{
				Data: req,
			},
		})
		if err != nil {
			log.Error().Msgf("cannot send chunk to server: ", err, stream.RecvMsg(nil))
			return err
		}
	}

	err = stream.CloseSend()
	if err != nil {
		log.Error().Msgf("cannot close stream connection: ", err)
		return err
	}
	<-waitc

	return nil
}

func (ss *StorageService) createClient() (mygrpc.StorageClient, error) {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(ss.Host+":"+ss.Port, opts...)
	if err != nil {
		log.Error().Msgf("cannot dial process-storage server: ", err)
		return nil, err
	}

	return mygrpc.NewStorageClient(conn), nil
}
