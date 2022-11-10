package main

import (
	"context"
	"flag"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
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
				if err != nil {
					return err
				}
				uploadFile(client, processId, f)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func uploadFile(client pb.StorageClient, processId string, file *os.File) {
	data := &pb.ProcessFileData_Metadata{
		Metadata: &pb.ProcessMetadata{
			ProcessId: processId,
			Filename:  filepath.Base(file.Name()),
			Encoding:  0,
			Type:      0,
			Version:   "0.2.0",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Hour)
	defer cancel()
	stream, err := client.UploadProcess(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("client.RouteChat failed: %v", err)
			}
			log.Println(in.Status)
		}
	}()

	if err := stream.Send(&pb.ProcessUploadRequest{
		Data: &pb.ProcessFileData{
			Data: data,
		},
	}); err != nil {
		log.Fatal(err.Error())
	}

	stat, err := os.Stat(file.Name())

	var fileSize = stat.Size()
	const fileChunk = 1 * (1 << 20)
	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))
	log.Printf("Splitting to %d pieces.\n", totalPartsNum)
	for i := uint64(0); i < totalPartsNum; i++ {
		partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
		partBuffer := make([]byte, partSize)
		_, err := file.Read(partBuffer)
		if err != nil {
			if err == io.EOF {
				err = stream.CloseSend()
			} else {
				log.Fatal(err.Error())
			}
			break
		}
		fileChunk := &pb.ProcessFileData_Chunk{Chunk: partBuffer}
		err = stream.Send(&pb.ProcessUploadRequest{
			Data: &pb.ProcessFileData{
				Data: fileChunk,
			},
		})
	}

	stream.CloseSend()
	<-waitc
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
