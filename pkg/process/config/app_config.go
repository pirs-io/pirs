// Package config is responsible for config loading and creating application context.
package config

import (
	"golang.org/x/net/context"
	"pirs.io/commons"
	"pirs.io/commons/db/mongo"
	"pirs.io/commons/parsers"
	metadata "pirs.io/process/metadata/service"
	"pirs.io/process/service"
	validation "pirs.io/process/validation/service"
	"time"
)

var (
	log                                           = commons.GetLoggerFor("config")
	processApplicationContext *ApplicationContext = nil
)

// A ProcessAppConfig contains loaded data from ENV config file
type ProcessAppConfig struct {
	GrpcIp                   string `mapstructure:"GRPC_IP"`
	GrpcPort                 int    `mapstructure:"GRPC_PORT"`
	UseGrpcReflection        bool   `mapstructure:"USE_GRPC_REFLECTION"`
	UploadFileMaxSize        int    `mapstructure:"UPLOAD_FILE_MAX_SIZE"`
	ChunkSize                int    `mapstructure:"CHUNK_SIZE"`
	StreamSeparator          string `mapstructure:"STREAM_SEPARATOR"`
	AllowedFileExtensions    string `mapstructure:"ALLOWED_FILE_EXTENSIONS"`
	MongoUri                 string `mapstructure:"MONGO_URI"`
	MongoName                string `mapstructure:"MONGO_NAME"`
	MongoDrop                bool   `mapstructure:"MONGO_DROP"`
	ContextTimeout           int    `mapstructure:"CONTEXT_TIMEOUT"`
	MetadataCollection       string `mapstructure:"METADATA_COLLECTION"`
	BasicMetadataCsv         string `mapstructure:"BASIC_METADATA_CSV"`
	PetriflowMetadataCsv     string `mapstructure:"PETRIFLOW_METADATA_CSV"`
	BPMNMetadataCsv          string `mapstructure:"BPMN_METADATA_CSV"`
	IgnoreWrongExtension     bool   `mapstructure:"IGNORE_WRONG_EXTENSION"`
	ProcessStoragePort       string `mapstructure:"PROCESS_STORAGE_PORT"`
	ProcessStorageHost       string `mapstructure:"PROCESS_STORAGE_HOST"`
	DependencyManagementPort string `mapstructure:"DEPENDENCY_MANAGEMENT_PORT"`
	DependencyManagementHost string `mapstructure:"DEPENDENCY_MANAGEMENT_HOST"`
}

func (p ProcessAppConfig) IsConfig() {}

// An ApplicationContext contains initialized config struct and all the main services
type ApplicationContext struct {
	AppConfig       *ProcessAppConfig
	ImportService   *service.ImportService
	DownloadService *service.DownloadService
}

// GetContext returns ApplicationContext instance, that is stored in a variable processApplicationContext
func GetContext() *ApplicationContext {
	return processApplicationContext
}

// InitApp initializes ProcessAppConfig from given configFilePath and initializes services by createApplicationContext().
// If success, ProcessAppConfig instance is returned, otherwise, it panics.
func InitApp(configFilePath string) (conf *ProcessAppConfig) {
	log.Info().Msg("Initializing Process service...")
	conf, confErr := commons.GetAppConfig(configFilePath, &ProcessAppConfig{})
	if confErr != nil {
		log.Fatal().Msgf("Unable to load application config for Process service! %s", confErr)
	}

	appCtx, contextErr := createApplicationContext(*conf)
	if contextErr != nil {
		log.Fatal().Msgf("Error intializing app context for Process service")
	}
	processApplicationContext = appCtx
	processApplicationContext.AppConfig = conf

	log.Info().Msg("Process service initialized!")
	return conf
}

func initMongoDatabase(conf ProcessAppConfig) mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Info().Msgf("connecting to %s", conf.MongoUri)

	client, err := mongo.NewClient(conf.MongoUri)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	err = client.Ping(ctx)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	log.Info().Msgf("successfully connected to %s", conf.MongoUri)
	return client
}

func parseCustomMetadataMappingFromCsv(csvPath string) map[string]string {
	csv := parsers.ReadCsvFile(csvPath, false)
	mapping := map[string]string{}
	for _, row := range csv {
		mapping[row[0]] = row[1]
	}
	return mapping
}

func createApplicationContext(conf ProcessAppConfig) (appContext *ApplicationContext, err error) {
	mongoClient := initMongoDatabase(conf)
	metadataRepo := mongo.NewMetadataRepository(mongoClient.Database(conf.MongoName, conf.MongoDrop), conf.MetadataCollection)
	validationService := validation.NewValidationService(conf.AllowedFileExtensions, conf.IgnoreWrongExtension)
	metadataService := metadata.NewMetadataService(
		metadataRepo,
		time.Duration(conf.ContextTimeout)*time.Second,
		parseCustomMetadataMappingFromCsv(conf.BasicMetadataCsv),
		parseCustomMetadataMappingFromCsv(conf.PetriflowMetadataCsv),
		parseCustomMetadataMappingFromCsv(conf.BPMNMetadataCsv),
	)
	storageService, err := service.NewStorageService(conf.ProcessStorageHost, conf.ProcessStoragePort, conf.ChunkSize)
	if err != nil {
		log.Error().Msgf("Process-Storage service was not correctly initialized: %v", err)
	}

	dependencyService, err := service.NewDependencyService(
		conf.DependencyManagementHost,
		conf.DependencyManagementPort,
		conf.StreamSeparator,
		conf.ChunkSize,
	)
	if err != nil {
		log.Error().Msgf("Dependency service was not correctly initialized: %v", err)
	}

	return &ApplicationContext{
		ImportService: &service.ImportService{
			ProcessStorageClient: storageService,
			ValidationService:    validationService,
			MetadataService:      metadataService,
			DependencyService:    dependencyService,
			MongoClient:          mongoClient,
		},
		DownloadService: &service.DownloadService{
			ValidationService: validationService,
			MetadataService:   metadataService,
			DependencyService: dependencyService,
		},
	}, nil
}
