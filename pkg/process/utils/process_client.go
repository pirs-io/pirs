package utils

import (
	"bufio"
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"os"
	"pirs.io/commons"
	mygrpc "pirs.io/process/grpc"
	"time"
)

var (
	log2 = commons.GetLoggerFor("process-client")
)

const (
	IP               = "localhost"
	PORT             = "8080"
	UPLOAD_FILENAME1 = "uvod.pdf"
	UPLOAD_FILENAME2 = "car.xml"
	CHUNK_SIZE       = 1024
	PARTIAL_URI      = "stu.fei.myproject2"
	MAX              = 3
)

func main() {
	serverAddress := flag.String("address", IP+":"+PORT, "the server address")
	flag.Parse()
	conn, err := grpc.Dial(*serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log2.Fatal().Msgf("cannot dial process server: ", err)
		return
	}
	processClient := mygrpc.NewProcessClient(conn)

	start := time.Now()
	// Call endpoints here
	for i := 0; i < MAX; i++ {
		uploadFile(processClient, "./pkg/process/"+UPLOAD_FILENAME2, UPLOAD_FILENAME2)
	}
	elapsed := time.Since(start)
	log2.Info().Msgf("Elapsed time: ", elapsed)
}

func uploadFile(processClient mygrpc.ProcessClient, processPath string, filename string) {
	file, err := os.Open(processPath)
	if err != nil {
		log2.Fatal().Msgf("cannot open file: ", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log2.Fatal().Msgf("cannot close file: ", err)
		}
	}(file)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := processClient.ImportProcess(ctx)
	if err != nil {
		log2.Fatal().Msgf("cannot upload file: ", err)
	}

	req := &mygrpc.ImportProcessRequest{
		Data: &mygrpc.ImportProcessRequest_FileInfo{
			FileInfo: &mygrpc.FileInfo{
				FileName: filename,
			},
		},
		PartialUri: PARTIAL_URI,
	}

	err = stream.Send(req)
	if err != nil {
		log2.Fatal().Msgf("cannot send file_info to server: ", err, stream.RecvMsg(nil))
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, CHUNK_SIZE)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log2.Fatal().Msgf("cannot read chunk to buffer: ", err)
		}

		req := &mygrpc.ImportProcessRequest{
			Data: &mygrpc.ImportProcessRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			log2.Fatal().Msgf("cannot send chunk to server: ", err, stream.RecvMsg(nil))
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log2.Fatal().Msgf("cannot receive response: ", err)
	}

	log2.Debug().Msgf("response with msg: %s, size: %dB", res.GetMessage(), res.GetTotalSize())
}
