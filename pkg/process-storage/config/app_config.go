package config

import (
	"pirs.io/commons"
	"pirs.io/process-storage/storage"
)

var (
	log                                           = commons.GetLoggerFor("config")
	trackerApplicationContext *ApplicationContext = nil
)

type ProcessStorageAppConfig struct {
	GrpcPort        int              `mapstructure:"GRPC_PORT"`
	StorageProvider storage.Provider `mapstructure:"STORAGE_PROVIDER"`
	RepoRootPath    string           `mapstructure:"GIT_ROOT"`
	Tenant          string           `mapstructure:"TENANT"`
	ChunkSize       int64            `mapstructure:"CHUNK_SIZE"`
}

func (t ProcessStorageAppConfig) IsConfig() {}

type ApplicationContext struct {
	AppConfig ProcessStorageAppConfig
}

func InitApp(configFilePath string) (conf *ProcessStorageAppConfig) {
	// config loading
	log.Info().Msg("Starting application")
	conf, confErr := commons.GetAppConfig(configFilePath, &ProcessStorageAppConfig{})
	if confErr != nil {
		log.Fatal().Msgf("Unable to load application config! %s", confErr)
	}

	// create app context
	appCtx, contextErr := createApplicationContext(*conf)
	if contextErr != nil {
		log.Fatal().Msgf("Error intializing app context")
	}
	trackerApplicationContext = appCtx

	log.Info().Msg("Application started!")
	trackerApplicationContext.AppConfig = *conf
	return conf
}

func GetContext() *ApplicationContext {
	return trackerApplicationContext
}

func createApplicationContext(conf ProcessStorageAppConfig) (appContext *ApplicationContext, err error) {
	return &ApplicationContext{conf}, nil
}
