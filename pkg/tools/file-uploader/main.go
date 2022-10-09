package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
	pb "pirs.io/process-storage/grpc"
	"time"
)

func main() {
	client := createClient()
	var processId string
	var processFilePath string
	app := cli.NewApp()
	app.Name = "Process file cli uploader"

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "processId",
			Value:       "",
			Destination: &processId,
		},
		&cli.StringFlag{
			Name:        "processFile",
			Value:       "",
			Destination: &processFilePath,
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:  "upload",
			Usage: "file to be send",
			Flags: flags,
			Action: func(c *cli.Context) error {
				f, err := os.Open(processFilePath)
				err = uploadFile(client, processId, f)
				if err != nil {
					return err
				}
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func uploadFile(client pb.StorageClient, processId string, file *os.File) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Hour)
	defer cancel()
	stream, err := client.UploadProcess(ctx)
	if err != nil {
		return err
	}

	data := &pb.ProcessFileData_Metadata{
		Metadata: &pb.ProcessMetadata{
			ProcessId: processId,
			Filename:  "",
			Encoding:  0,
			Type:      0,
		}}

	if err := stream.Send(&pb.ProcessUploadRequest{
		Data: &pb.ProcessFileData{
			Data: data,
		},
	}); err != nil {
		return err
	}

	if file == nil {
		return nil
	}
	chunkSize := 100
	chunk := make([]byte, chunkSize)
	for {
		bytesread, err := file.Read(chunk)

		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			if err == io.EOF {
				err = stream.CloseSend()
			}
			break
		}

		fileChunk := &pb.ProcessFileData_Chunk{Chunk: chunk}
		err = stream.Send(&pb.ProcessUploadRequest{
			Data: &pb.ProcessFileData{
				Data: fileChunk,
			},
		})
		if bytesread <= chunkSize {
			err = stream.CloseSend()
		}

	}
	return nil
}

func createClient() pb.StorageClient {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	//defer conn.Close()
	return pb.NewStorageClient(conn)
}
