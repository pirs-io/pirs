package config

import (
	"pirs.io/commons"
	"pirs.io/dependency-management/detection"
)

var (
	log                           = commons.GetLoggerFor("config")
	depAppCtx *ApplicationContext = nil
)

// A DependencyAppConfig contains loaded data from ENV config file
type DependencyAppConfig struct {
	GrpcIp            string `mapstructure:"GRPC_IP"`
	GrpcPort          int    `mapstructure:"GRPC_PORT"`
	UseGrpcReflection bool   `mapstructure:"USE_GRPC_REFLECTION"`
	StreamSeparator   string `mapstructure:"STREAM_SEPARATOR"`
}

func (p DependencyAppConfig) IsConfig() {}

// An ApplicationContext contains initialized config struct and all the main services
type ApplicationContext struct {
	AppConfig        *DependencyAppConfig
	DetectionService *detection.DetectionService
}

// GetContext returns ApplicationContext instance, that is stored in a variable depAppCtx
func GetContext() *ApplicationContext {
	return depAppCtx
}

// InitApp initializes DependencyAppConfig from given configFilePath and initializes services by createApplicationContext().
// If success, DependencyAppConfig instance is returned, otherwise, it panics.
func InitApp(configFilePath string) (conf *DependencyAppConfig) {
	log.Info().Msg("Initializing Dependency Management service...")
	conf, confErr := commons.GetAppConfig(configFilePath, &DependencyAppConfig{})
	if confErr != nil {
		log.Fatal().Msgf("Unable to load application config for Dependency Management service! %s", confErr)
	}

	appCtx, contextErr := createApplicationContext(*conf)
	if contextErr != nil {
		log.Fatal().Msgf("Error intializing app context for Dependency Management service")
	}
	depAppCtx = appCtx
	depAppCtx.AppConfig = conf

	log.Info().Msg("Dependency Management service initialized!")
	return conf
}

func createApplicationContext(conf DependencyAppConfig) (appContext *ApplicationContext, err error) {
	return &ApplicationContext{
		DetectionService: &detection.DetectionService{},
	}, nil
}
