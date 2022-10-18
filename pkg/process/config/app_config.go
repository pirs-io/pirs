package config

import (
	"pirs.io/commons"
)

var (
	log                                           = commons.GetLoggerFor("config")
	processApplicationContext *ApplicationContext = nil
)

type ProcessAppConfig struct {
	GrpcIp            string `mapstructure:"GRPC_IP"`
	GrpcPort          int    `mapstructure:"GRPC_PORT"`
	UseGrpcReflection bool   `mapstructure:"USE_GRPC_REFLECTION"`
}

func (p ProcessAppConfig) IsConfig() {}

type ApplicationContext struct {
	AppConfig *ProcessAppConfig
}

func GetContext() *ApplicationContext {
	return processApplicationContext
}

func InitApp(configFilePath string) (conf *ProcessAppConfig) {
	// config loading
	log.Info().Msg("Initializing Process service...")
	conf, confErr := commons.GetAppConfig(configFilePath, &ProcessAppConfig{})
	if confErr != nil {
		log.Fatal().Msgf("Unable to load application config for Process service! %s", confErr)
	}

	// create app context
	appCtx, contextErr := createApplicationContext(*conf)
	if contextErr != nil {
		log.Fatal().Msgf("Error intializing app context for Process service")
	}
	processApplicationContext = appCtx
	processApplicationContext.AppConfig = conf

	log.Info().Msg("Process service initialized!")
	return conf
}

func createApplicationContext(conf ProcessAppConfig) (appContext *ApplicationContext, err error) {
	return &ApplicationContext{}, nil
}
