package service

import (
	"bytes"
	"errors"
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
	client    mygrpc.StorageClient
}

// A ResourceAdapter is wrapper for metadata and file data. It is intended for StorageService.SaveFiles.
type ResourceAdapter struct {
	Metadata domain.Metadata
	FileData []byte
}

func NewStorageService(host string, port string, chunkSize int) (*StorageService, error) {
	service := &StorageService{
		Host:      host,
		Port:      port,
		ChunkSize: chunkSize,
	}
	client, err := service.createClient()
	if err != nil {
		return service, err
	}
	service.client = client
	return service, nil
}

// SaveFiles runs in a separate goroutine. It waits for resources coming through the forResource channel. Resource is
// ResourceAdapter, which is a wrapper for metadata and []bytes. On success no error is sent to the caller (more likely
// ImportService). If it fails, an error is sent through the forResponse channel.
func (ss *StorageService) SaveFiles(reqCtx context.Context, forResource <-chan ResourceAdapter, forResponse chan<- error) {
	defer close(forResponse)
	stream, err := ss.establishConnection(reqCtx)
	if err != nil {
		forResponse <- err
		return
	} else {
		forResponse <- nil
	}

	for resource := range forResource {
		// todo checksum
		reqMetadata := ss.transformMetadataToRequest(&resource.Metadata)
		err = ss.sendMetadataRequest(stream, reqMetadata)
		if err != nil {
			forResponse <- err
			return
		}

		sync := make(chan bool)
		c := ss.createFileChunksAsync(resource.FileData, sync)
		err = ss.sendFileChunks(stream, c, sync)
		if err != nil {
			forResponse <- err
		} else {
			_ = ss.destroyConnection(stream)
			forResponse <- ss.checkResponse(stream)
		}
	}
}

func (ss *StorageService) establishConnection(ctx context.Context) (mygrpc.Storage_UploadProcessClient, error) {
	var err error

	if ss.client == nil {
		err = errors.New("Process-Storage client is not initialized")
		log.Error().Msg(err.Error())
		return nil, err
	}

	stream, err := ss.client.UploadProcess(ctx)
	if err != nil {
		log.Error().Msgf("could not establish stream connection: %v", err)
		return nil, err
	}
	return stream, nil
}

func (ss *StorageService) destroyConnection(stream mygrpc.Storage_UploadProcessClient) error {
	err := stream.CloseSend()
	if err != nil {
		log.Error().Msgf("cannot close stream connection: %v", err)
		return err
	}
	return ss.checkResponse(stream)
}

// sendMetadataRequest takes stream and metadata. This data get wrapped and sent through the stream. todo
func (ss *StorageService) sendMetadataRequest(stream mygrpc.Storage_UploadProcessClient, metadata *mygrpc.ProcessFileData_Metadata) error {
	if err := stream.Send(&mygrpc.ProcessUploadRequest{
		Data: &mygrpc.ProcessFileData{
			Data: metadata,
		},
	}); err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	return ss.checkResponse(stream)
}

func (ss *StorageService) checkResponse(stream mygrpc.Storage_UploadProcessClient) error {
	_, err := stream.Recv()
	if err != nil && err != io.EOF {
		log.Error().Msgf("client.RouteChat failed: %v", err)
		return err
	}
	return nil
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
		err = ss.checkResponse(stream)
		if err != nil {
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
