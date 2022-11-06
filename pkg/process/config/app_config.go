package config

import (
	"fmt"
	"golang.org/x/net/context"
	"pirs.io/commons"
	"pirs.io/commons/mongo"
	"pirs.io/commons/parsers"
	"pirs.io/process/metadata"
	metadataMongo "pirs.io/process/metadata/repository/mongo"
	"pirs.io/process/mocks"
	"pirs.io/process/service"
	"pirs.io/process/validation"
	"time"
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
	MongoHost             string `mapstructure:"MONGO_HOST"`
	MongoPort             string `mapstructure:"MONGO_PORT"`
	MongoUser             string `mapstructure:"MONGO_USER"`
	MongoPass             string `mapstructure:"MONGO_PASS"`
	MongoName             string `mapstructure:"MONGO_NAME"`
	MongoDrop             bool   `mapstructure:"MONGO_DROP"`
	ContextTimeout        int    `mapstructure:"CONTEXT_TIMEOUT"`
	MetadataCollection    string `mapstructure:"METADATA_COLLECTION"`
	CustomMetadataCsv     string `mapstructure:"CUSTOM_METADATA_CSV"`
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

func initMongoDatabase(conf ProcessAppConfig) mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbHost := conf.MongoHost
	dbPort := conf.MongoPort
	dbUser := conf.MongoUser
	dbPass := conf.MongoPass
	dbName := conf.MongoName
	mongodbURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	if dbUser == "" || dbPass == "" {
		mongodbURI = fmt.Sprintf("mongodb://%s:%s/%s", dbHost, dbPort, dbName)
	}
	log.Info().Msgf("connecting to %s", mongodbURI)

	client, err := mongo.NewClient(mongodbURI)
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
	log.Info().Msgf("successfully connected to %s", mongodbURI)
	return client
}

func parseCustomMetadataXpathsFromCsv(csvPath string) [][]string {
	return parsers.ReadCsvFile(csvPath)
}

func createApplicationContext(conf ProcessAppConfig) (appContext *ApplicationContext, err error) {
	mongoClient := initMongoDatabase(conf)
	metadataRepo := metadataMongo.NewMetadataRepository(mongoClient.Database(conf.MongoName, conf.MongoDrop), conf.MetadataCollection)
	return &ApplicationContext{
		ImportService: &service.ImportService{
			// todo mockup
			ProcessStorageClient: mocks.NewDiskProcessStore("./pkg/process/imported-files"),
			MongoClient:          &mongoClient,
			ValidationService:    validation.NewValidationService(conf.AllowedFileExtensions),
			MetadataService: metadata.NewMetadataService(
				*metadataRepo,
				time.Duration(conf.ContextTimeout)*time.Second,
				parseCustomMetadataXpathsFromCsv(conf.CustomMetadataCsv),
			),
		},
	}, nil
}
