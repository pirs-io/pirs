package adapter

import (
	"golang.org/x/net/context"
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

func (a *GitAdapter) DownloadProcess(downloadRequest *pb.ProcessDownloadRequest) (*pb.ProcessMetadata, []byte, error) {
	return a.GitClient.DownloadProcess(downloadRequest)
}
