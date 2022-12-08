package adapter

import (
	"golang.org/x/net/context"
	"io"
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

func (a *GitAdapter) DownloadProcess(downloadRequest *pb.ProcessDownloadRequest, w *io.PipeWriter) (*pb.ProcessMetadata, error) {
	return a.GitClient.DownloadProcess(downloadRequest, w)
}
