package config

import (
	"pirs.io/commons"
	"pirs.io/process-storage/storage"
)

var (
	log                                           = commons.GetLoggerFor("config")
	trackerApplicationContext *ApplicationContext = nil
)

type TrackerAppConfig struct {
	GrpcPort        int              `mapstructure:"GRPC_PORT"`
	StorageProvider storage.Provider `mapstructure:"STORAGE_PROVIDER"`
	GitUrl          string           `mapstructure:"GIT_URL"`
	GitUsername     string           `mapstructure:"GIT_USERNAME"`
	GitPassword     string           `mapstructure:"GIT_PASSWORD"`
}

func (t TrackerAppConfig) IsConfig() {}

type ApplicationContext struct {
	AppConfig TrackerAppConfig
}

func InitApp(configFilePath string) (conf *TrackerAppConfig) {
	// config loading
	log.Info().Msg("Starting application")
	conf, confErr := commons.GetAppConfig(configFilePath, &TrackerAppConfig{})
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

func createApplicationContext(conf TrackerAppConfig) (appContext *ApplicationContext, err error) {
	return &ApplicationContext{conf}, nil
}
