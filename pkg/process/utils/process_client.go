package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"os"
	"pirs.io/commons"
	"pirs.io/process/domain"
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
	URI1             = "awd.awd.awd.awd:11"
	URI2             = ""
	CHUNK_SIZE       = 1024
	PARTIAL_URI      = "stu.fei.myproject2"
	MAX_IMPORT       = 1
	MAX_DOWNLOAD     = 1
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

	importProcess(processClient)
	downloadProcess(processClient)
}

func downloadProcess(client mygrpc.ProcessClient) {
	start := time.Now()
	// Call endpoints here
	for i := 0; i < MAX_DOWNLOAD; i++ {
		downloadData(client, PARTIAL_URI+".car:1")
	}
	elapsed := time.Since(start)
	log2.Info().Msgf("downloadProcess elapsed time: ", elapsed)
}

func downloadData(client mygrpc.ProcessClient, uri string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &mygrpc.DownloadProcessRequest{
		Uri: uri,
	}

	response, err := client.DownloadProcess(ctx, req)
	if err != nil {
		return
	}
	if response.Metadata == nil {
		return
	}

	// example handling response
	metadataFromResponse := domain.Metadata{}
	jsonString, _ := json.Marshal(response.Metadata)
	err = json.Unmarshal(jsonString, &metadataFromResponse)
	if err != nil {
		return
	}
}

func importProcess(client mygrpc.ProcessClient) {
	start := time.Now()
	// Call endpoints here
	for i := 0; i < MAX_IMPORT; i++ {
		uploadFile(client, "./pkg/process/"+UPLOAD_FILENAME2, UPLOAD_FILENAME2)
	}
	elapsed := time.Since(start)
	log2.Info().Msgf("importProcess elapsed time: ", elapsed)
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
