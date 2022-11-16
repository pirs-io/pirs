package adapter

import (
	"golang.org/x/net/context"
	"io"
	"pirs.io/commons"
	"pirs.io/process-storage/config"
	pb "pirs.io/process-storage/grpc"
	"pirs.io/process-storage/storage"
)

var log = commons.GetLoggerFor("storage_client")

type IStorageAdapter interface {
	SaveProcess(processMetadata *pb.ProcessMetadata, file []byte) error
	DownloadProcess(downloadRequest *pb.ProcessDownloadRequest, w *io.PipeWriter) (*pb.ProcessMetadata, error)
}

func MakeStorageAdapter(ctx context.Context, provider storage.Provider) (IStorageAdapter, error) {
	switch provider {
	case storage.GitStorageProvider:
		gitAdapter := &GitAdapter{
			ctx: ctx,
			GitClient: storage.GitClient{
				Context:      ctx,
				RepoRootPath: config.GetContext().AppConfig.RepoRootPath,
				Tenant:       config.GetContext().AppConfig.Tenant,
				ChunkSize:    config.GetContext().AppConfig.ChunkSize,
			},
		}
		err := gitAdapter.GitClient.InitializeStorage()
		if err := commons.CheckAndLog(err, log); err != nil {
			return nil, err
		}

		return gitAdapter, nil
	default:
		log.Error().Msgf("Unsupported storage provider: %s", provider)
		return nil, nil
	}
}
