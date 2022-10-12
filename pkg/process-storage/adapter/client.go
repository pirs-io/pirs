package adapter

import (
	"golang.org/x/net/context"
	"os"
	"pirs.io/common"
	"pirs.io/process-storage/config"
	pb "pirs.io/process-storage/grpc"
	"pirs.io/process-storage/storage"
)

var log = common.GetLoggerFor("storage_client")

type IStorageAdapter interface {
	SaveProcess(processMetadata *pb.ProcessMetadata, file []byte) error
	DownloadProcess(processId string) (*os.File, error)
	Close() error
}

func MakeStorageAdapter(ctx context.Context, provider storage.Provider) (IStorageAdapter, error) {
	switch provider {
	case storage.GitStorageProvider:
		dir, err := os.MkdirTemp(".", "repo")
		gitAdapter := &GitAdapter{
			ctx: ctx,
			GitClient: storage.GitClient{
				Context:         ctx,
				Url:             config.GetContext().AppConfig.GitUrl,
				Username:        config.GetContext().AppConfig.GitUsername,
				Password:        config.GetContext().AppConfig.GitPassword,
				TempRepoDirPath: dir,
			},
		}
		err = gitAdapter.GitClient.InitializeStorage()
		if err := common.CheckAndLog(err, log); err != nil {
			return nil, err
		}

		return gitAdapter, nil
	default:
		log.Error().Msgf("Unsupported storage provider: %s", provider)
		return nil, nil
	}
}
