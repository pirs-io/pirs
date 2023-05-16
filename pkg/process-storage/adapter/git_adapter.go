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

func (r *GitAdapter) SaveProcess(processMetadata *pb.ProcessMetadata, file []byte) error {
	return r.GitClient.SaveFile(processMetadata, file)
}

func (r *GitAdapter) DownloadProcess(downloadRequest *pb.ProcessDownloadRequest, w *io.PipeWriter) (*pb.ProcessMetadata, error) {
	return r.GitClient.DownloadProcess(downloadRequest, w)
}
