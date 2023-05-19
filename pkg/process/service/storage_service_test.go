package service

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"pirs.io/commons"
	"pirs.io/commons/domain"
	"pirs.io/commons/enums"
	mygrpc "pirs.io/process/grpc"
	"testing"
)

const (
	RESOURCES = "../resources/testing/"
	PF2       = "v5_petriflow.xml"
)

var (
	log_test = commons.GetLoggerFor("StorageServiceTest")
)

func TestStorageService_transformMetadataToRequest(t *testing.T) {
	ss := &StorageService{}

	m := domain.NewMetadata()
	m.URI = "com.example.project.my_process:2"
	m.FileName = "awd.xml"
	m.Encoding = "UTF-8"
	m.ProcessType = enums.Petriflow

	reqM := ss.transformMetadataToRequest(m, "")

	assert.Equal(t, m.URI, reqM.Metadata.ProcessId)
	assert.Equal(t, m.FileName, reqM.Metadata.Filename)
	assert.Equal(t, mygrpc.Encoding(int32(0)), reqM.Metadata.Encoding)
	assert.Equal(t, mygrpc.ProcessType(int32(0)), reqM.Metadata.Type)
}

func TestStorageService_createFileChunksAsync(t *testing.T) {
	fileData := readFile(RESOURCES + PF2)
	ss := &StorageService{
		ChunkSize: 1024,
	}

	var receivedData []byte
	sync := make(chan bool)
	c := ss.createFileChunksAsync(fileData, sync)
	for chunk := range c {
		receivedData = append(receivedData, chunk...)
		sync <- true
	}
	close(sync)

	assert.ElementsMatch(t, fileData, receivedData)
}

func readFile(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		log_test.Fatal().Msg(err.Error())
	}
	reader := bufio.NewReader(file)
	buf := make([]byte, 1024)
	var data []byte
	var totalSize = 0
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				log_test.Fatal().Msg(err.Error())
			}
			break
		}
		data = append(data, buf[:n]...)
		totalSize = totalSize + n
	}

	err = file.Close()
	if err != nil {
		log_test.Fatal().Msg(err.Error())
	}

	return data
}
