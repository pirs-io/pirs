package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"os"
	"pirs.io/commons"
	"pirs.io/commons/domain"
	mygrpc "pirs.io/process/grpc"
	"strings"
	"time"
)

var (
	log2 = commons.GetLoggerFor("process-client")
)

const (
	IP           = "localhost"
	PORT         = "8080"
	CHUNK_SIZE   = 1024
	MAX_IMPORT   = 1
	MAX_DOWNLOAD = 1
)

func main() {
	serverAddress := flag.String("address", IP+":"+PORT, "the server address")
	flag.Parse()
	conn, err := grpc.Dial(*serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log2.Fatal().Msgf("cannot dial process server: %v", err)
		return
	}
	processClient := mygrpc.NewProcessClient(conn)
	reader := bufio.NewReader(os.Stdin)

	printMainHeadline()
	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		switch text {
		case "import":
			importProcess(processClient, reader)
			printMainHeadline()
		case "download":
			downloadDialog(reader, processClient)
			printMainHeadline()
		case "quit":
			return
		default:
			fmt.Println("Bad command, try again")
		}
	}
}

func printMainHeadline() {
	fmt.Println("\nSimple client for process-service")
	fmt.Println("Commands:")
	fmt.Println("\timport - Calls import endpoint")
	fmt.Println("\tdownload - Calls download endpoint")
	fmt.Println("\tquit - Stops the client")
}

func downloadDialog(reader *bufio.Reader, processClient mygrpc.ProcessClient) {
	fmt.Println("\nPlease select what you want to download")
	fmt.Println("Commands:")
	fmt.Println("\tprocess - Downloads specified process metadata")
	fmt.Println("\tpackage - Downloads specified package of process metadata")
	fmt.Println("\tback - goes back to the main menu")
	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		switch text {
		case "process":
			downloadProcess(processClient, reader)
			return
		case "package":
			downloadPackage(processClient, reader)
			return
		case "back":
			return
		default:
			fmt.Println("Bad command, try again")
		}
	}
}

func downloadProcess(client mygrpc.ProcessClient, reader *bufio.Reader) {
	fmt.Println("\nSpecify process URI (f.e. \"myorg.mytenant.myproject.myprocess:1\":")
	fmt.Print("-> ")
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)

	for i := 0; i < MAX_DOWNLOAD; i++ {
		downloadProcessData(client, text)
	}
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
		log2.Error().Msgf("%v", err)
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
			log2.Error().Msgf("%v", err)
			return
		}

		downloadedMetadata = append(downloadedMetadata, metadataFromResponse)
	}
	log2.Info().Msgf("Downloaded metadata of %d processes", len(downloadedMetadata))
}

func downloadPackage(client mygrpc.ProcessClient, reader *bufio.Reader) {
	fmt.Println("\nSpecify package URI (f.e. \"myorg.mytenant.myproject\":")
	fmt.Print("-> ")
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)

	for i := 0; i < MAX_DOWNLOAD; i++ {
		downloadPackageData(client, text)
	}
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
		log2.Error().Msgf("%v", err)
		return
	}

	resp, err := stream.Recv()
	if err != nil {
		log2.Error().Msgf("Cannot receive status response: %v", err)
		return
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
			log2.Error().Msgf("%v", err)
			return
		}

		downloadedMetadata = append(downloadedMetadata, metadataFromResponse)
	}
	log2.Info().Msgf("Downloaded metadata of %d processes", len(downloadedMetadata))
}

func importProcess(client mygrpc.ProcessClient, reader *bufio.Reader) {
	fmt.Println("\nSpecify path for process file (f.e. \"./pkg/process/myprocess.xml\":")
	fmt.Print("-> ")
	path, _ := reader.ReadString('\n')
	path = strings.Replace(path, "\n", "", -1)
	name := parseNameFromPath(path)

	fmt.Println("\nSpecify package URI (f.e. \"myorg.mytenant.myproject\":")
	fmt.Print("-> ")
	uri, _ := reader.ReadString('\n')
	uri = strings.Replace(uri, "\n", "", -1)

	for i := 0; i < MAX_IMPORT; i++ {
		uploadFile(client, path, name, uri)
	}
}

func parseNameFromPath(path string) string {
	splitPath := strings.Split(path, "/")
	return splitPath[len(splitPath)-1]
}

func uploadFile(processClient mygrpc.ProcessClient, processPath string, filename string, partialUri string) {
	file, err := os.Open(processPath)
	if err != nil {
		log2.Error().Msgf("cannot open file: %v", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log2.Error().Msgf("cannot close file: %v", err)
			return
		}
	}(file)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := processClient.Import(ctx)
	if err != nil {
		log2.Error().Msgf("cannot upload file: %v", err)
		return
	}

	req := &mygrpc.ImportRequest{
		Data: &mygrpc.ImportRequest_FileInfo{
			FileInfo: &mygrpc.FileInfo{
				FileName: filename,
			},
		},
		PartialUri: partialUri,
	}

	err = stream.Send(req)
	if err != nil {
		log2.Error().Msgf("cannot send file_info to server: %v %v", err, stream.RecvMsg(nil))
		return
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, CHUNK_SIZE)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log2.Error().Msgf("cannot read chunk to buffer: %v", err)
			return
		}

		req := &mygrpc.ImportRequest{
			Data: &mygrpc.ImportRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			log2.Error().Msgf("cannot send chunk to server: %v %v", err, stream.RecvMsg(nil))
			return
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log2.Error().Msgf("cannot receive response: %v", err)
		return
	}

	log2.Info().Msgf("response with msg: %s, size: %dB", res.GetMessage(), res.GetTotalSize())
}
