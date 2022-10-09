package adapter

import (
	"golang.org/x/net/context"
	"io/ioutil"
	"os"
	"pirs.io/common"
	"pirs.io/process-storage/config"
	pb "pirs.io/process-storage/grpc"
	"pirs.io/process-storage/storage"
)

var log = common.GetLoggerFor("storage_client")

type IStorageAdapter interface {
	SaveProcess(processMetadata *pb.ProcessMetadata, file *os.File) error
	DownloadProcess(processId string) (*os.File, error)
}

func MakeStorageAdapter(ctx context.Context, provider storage.Provider) (IStorageAdapter, error) {
	switch provider {
	case storage.GitStorageProvider:
		dir, err := ioutil.TempDir(".", "repo")
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
		defer func(GitClient *storage.GitClient) {
			err := GitClient.Close()
			if err != nil {
				common.CheckAndLog(err, log)
			}
		}(&gitAdapter.GitClient)
		return gitAdapter, nil
	default:
		log.Error().Msgf("Unsupported storage provider: %s", provider)
		return nil, nil
	}
}
