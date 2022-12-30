package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"os"
	"pirs.io/commons"
	"pirs.io/process/domain"
	mygrpc "pirs.io/process/grpc"
	"strings"
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
	UPLOAD_FILENAME3 = "service.xml"
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
	//downloadProcess(processClient)
	//downloadPackage(processClient)
}

func downloadProcess(client mygrpc.ProcessClient) {
	start := time.Now()
	// Call endpoints here
	for i := 0; i < MAX_DOWNLOAD; i++ {
		downloadProcessData(client, PARTIAL_URI+".car:1")
	}
	elapsed := time.Since(start)
	log2.Info().Msgf("downloadProcess elapsed time: ", elapsed)
}

func downloadProcessData(client mygrpc.ProcessClient, uri string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &mygrpc.DownloadRequest{
		TargetUri: uri,
		IsPackage: false,
	}
	stream, err := client.Download(ctx, req)
	if err != nil {
		return
	}

	resp, err := stream.Recv()
	if err != nil {
		log2.Error().Msgf("Cannot receive status response: %v", err)
	}
	if !strings.Contains(resp.Message, codes.OK.String()) {
		log2.Error().Msgf("Downloading metadata failed with code: %s", resp.Message)
		return
	}
	var downloadedMetadata []domain.Metadata
	for {
		resp, err = stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log2.Error().Msgf("Downloading metadata failed with error: %v", err)
			return
		}

		metadataFromResponse := domain.Metadata{}
		jsonString, _ := json.Marshal(resp.Metadata)
		err = json.Unmarshal(jsonString, &metadataFromResponse)
		if err != nil {
			return
		}

		downloadedMetadata = append(downloadedMetadata, metadataFromResponse)
	}
	log2.Info().Msgf("Downloaded metadata of %d processes", len(downloadedMetadata))
}

func downloadPackage(client mygrpc.ProcessClient) {
	start := time.Now()
	// Call endpoints here
	for i := 0; i < MAX_DOWNLOAD; i++ {
		downloadPackageData(client, PARTIAL_URI)
	}
	elapsed := time.Since(start)
	log2.Info().Msgf("downloadPackage elapsed time: ", elapsed)
}

func downloadPackageData(client mygrpc.ProcessClient, packageUri string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &mygrpc.DownloadRequest{
		TargetUri: packageUri,
		IsPackage: true,
	}
	stream, err := client.Download(ctx, req)
	if err != nil {
		return
	}

	resp, err := stream.Recv()
	if err != nil {
		log2.Error().Msgf("Cannot receive status response: %v", err)
	}
	if !strings.Contains(resp.Message, codes.OK.String()) {
		log2.Error().Msgf("Downloading metadata failed with code: %s", resp.Message)
		return
	}
	var downloadedMetadata []domain.Metadata
	for {
		resp, err = stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log2.Error().Msgf("Downloading metadata failed with error: %v", err)
			return
		}

		metadataFromResponse := domain.Metadata{}
		jsonString, _ := json.Marshal(resp.Metadata)
		err = json.Unmarshal(jsonString, &metadataFromResponse)
		if err != nil {
			return
		}

		downloadedMetadata = append(downloadedMetadata, metadataFromResponse)
	}
	log2.Info().Msgf("Downloaded metadata of %d processes", len(downloadedMetadata))
}

func importProcess(client mygrpc.ProcessClient) {
	start := time.Now()
	// Call endpoints here
	for i := 0; i < MAX_IMPORT; i++ {
		uploadFile(client, "./pkg/process/"+UPLOAD_FILENAME2, UPLOAD_FILENAME2)
		uploadFile(client, "./pkg/process/"+UPLOAD_FILENAME3, UPLOAD_FILENAME3)
		uploadFile(client, "./pkg/process/"+UPLOAD_FILENAME3, UPLOAD_FILENAME3)
		uploadFile(client, "./pkg/process/"+UPLOAD_FILENAME3, UPLOAD_FILENAME3)
		uploadFile(client, "./pkg/process/"+UPLOAD_FILENAME3, UPLOAD_FILENAME3)
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

	ctx, cancel := context.WithTimeout(context.Background(), 5000*time.Second)
	defer cancel()

	stream, err := processClient.Import(ctx)
	if err != nil {
		log2.Fatal().Msgf("cannot upload file: ", err)
	}

	req := &mygrpc.ImportRequest{
		Data: &mygrpc.ImportRequest_FileInfo{
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

		req := &mygrpc.ImportRequest{
			Data: &mygrpc.ImportRequest_ChunkData{
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
