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
	"pirs.io/process/enums"
	mygrpc "pirs.io/process/grpc"
	"strings"
)

var (
	log = commons.GetLoggerFor("StorageService")
)

// A StorageService is used to save process bytes into storage managed by Process-Storage. Host and Port describes
// Process-Storage instance. It's initialized in the config package.
type StorageService struct {
	Port      string
	Host      string
	ChunkSize int
}

// SaveFile takes domain.Metadata with array of bytes of process file and streams this data with grpc.StorageClient.
// On success no error is returned. If it fails, an error is returned.
func (ss *StorageService) SaveFile(reqCtx context.Context, m domain.Metadata, fileData []byte) error {
	client, err := ss.createClient()
	if err != nil {
		return err
	}

	reqMetadata := ss.transformMetadataToRequest(&m)
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

	err = ss.sendMetadataRequest(stream, reqMetadata)
	if err != nil {
		return err
	}

	sync := make(chan bool)
	c := ss.createFileChunksAsync(fileData, sync)
	err = ss.sendFileChunks(stream, c, sync)
	if err != nil {
		return err
	}

	err = stream.CloseSend()
	if err != nil {
		log.Error().Msgf("cannot close stream connection: %v", err)
		return err
	}
	// wait for goroutine to handle all the data and end.
	<-waitc

	return nil
}

// sendMetadataRequest takes stream and metadata. This data get wrapped and sent through the stream.
func (ss *StorageService) sendMetadataRequest(stream mygrpc.Storage_UploadProcessClient, metadata *mygrpc.ProcessFileData_Metadata) error {
	if err := stream.Send(&mygrpc.ProcessUploadRequest{
		Data: &mygrpc.ProcessFileData{
			Data: metadata,
		},
	}); err != nil {
		log.Error().Msg(err.Error())
		return err
	} else {
		return nil
	}
}

// sendFileChunks takes stream, c channel to read bytes to send and sync channel where it signals to createFileChunksAsync,
// that data were sent. This function is the receiver of c and sender of sync.
func (ss *StorageService) sendFileChunks(stream mygrpc.Storage_UploadProcessClient, c <-chan []byte, sync chan<- bool) error {
	for chunk := range c {
		req := &mygrpc.ProcessFileData_Chunk{Chunk: chunk}

		err := stream.Send(&mygrpc.ProcessUploadRequest{
			Data: &mygrpc.ProcessFileData{
				Data: req,
			},
		})
		if err != nil {
			log.Error().Msgf("cannot send chunk to server: %v, %v", err, stream.RecvMsg(nil))
			close(sync)
			return err
		}
		sync <- true
	}
	close(sync)
	return nil
}

// createFileChunksAsync takes array of bytes of process file and sync channel. It firstly creates read-only channel, that
// gets returned. It also executes a goroutine. This goroutine reads byte array, which is divided into chunks and sent
// through the created channel. This function uses sync channel to make sure the chunk is read before writing new one in.
func (ss *StorageService) createFileChunksAsync(fileData []byte, sync <-chan bool) <-chan []byte {
	c := make(chan []byte)

	go func() {
		defer close(c)
		reader := bytes.NewReader(fileData)
		buffer := make([]byte, ss.ChunkSize)
		for {
			n, err := reader.Read(buffer)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Error().Msgf("cannot read chunk to buffer: %v", err)
				return
			}

			c <- buffer[:n]
			<-sync
		}
	}()

	return c
}

func (ss *StorageService) transformMetadataToRequest(m *domain.Metadata) *mygrpc.ProcessFileData_Metadata {
	return &mygrpc.ProcessFileData_Metadata{
		Metadata: &mygrpc.ProcessMetadata{
			ProcessId: m.URI,
			Filename:  m.FileName,
			Encoding:  ss.transformEncoding(m.Encoding),
			Type:      ss.transformProcessType(m.ProcessType),
		},
	}
}

func (ss *StorageService) transformProcessType(pt enums.ProcessType) mygrpc.ProcessType {
	mappedInt := mygrpc.ProcessType_value[pt.String()]
	return mygrpc.ProcessType(mappedInt)
}

func (ss *StorageService) transformEncoding(e string) mygrpc.Encoding {
	trimmed := strings.Replace(e, "-", "", -1)
	return mygrpc.Encoding(mygrpc.Encoding_value[trimmed])
}

func (ss *StorageService) createClient() (mygrpc.StorageClient, error) {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(ss.Host+":"+ss.Port, opts...)
	if err != nil {
		log.Error().Msgf("cannot dial process-storage server: %v", err)
		return nil, err
	}

	return mygrpc.NewStorageClient(conn), nil
}
