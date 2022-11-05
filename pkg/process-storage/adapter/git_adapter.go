package adapter

import (
	"golang.org/x/net/context"
	"os"
	pb "pirs.io/process-storage/grpc"
	"pirs.io/process-storage/storage"
)

type GitAdapter struct {
	ctx       context.Context
	GitClient storage.GitClient
}

func (a *GitAdapter) SaveProcess(processMetadata *pb.ProcessMetadata, file []byte) error {
	return a.GitClient.SaveFile(processMetadata, file)
}

func (a *GitAdapter) DownloadProcess(processId string) (*os.File, error) {
	return nil, nil
}
