package config

import (
	"pirs.io/commons"
	"pirs.io/process/mock"
	"pirs.io/process/service"
	"pirs.io/process/validation"
)

var (
	log                                           = commons.GetLoggerFor("config")
	processApplicationContext *ApplicationContext = nil
)

type ProcessAppConfig struct {
	GrpcIp                string `mapstructure:"GRPC_IP"`
	GrpcPort              int    `mapstructure:"GRPC_PORT"`
	UseGrpcReflection     bool   `mapstructure:"USE_GRPC_REFLECTION"`
	UploadFileMaxSize     int    `mapstructure:"UPLOAD_FILE_MAX_SIZE"`
	AllowedFileExtensions string `mapstructure:"ALLOWED_FILE_EXTENSIONS"`
}

func (p ProcessAppConfig) IsConfig() {}

type ApplicationContext struct {
	AppConfig         *ProcessAppConfig
	ImportService     *service.ImportService
	ValidationService *validation.ValidationService
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
	return &ApplicationContext{
		ImportService: &service.ImportService{
			// todo mockup
			ProcessStorageClient: mock.NewDiskProcessStore("./pkg/process/imported-files"),
			ValidationService:    validation.NewValidationService(conf.AllowedFileExtensions),
		},
	}, nil
}
