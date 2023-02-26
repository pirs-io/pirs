package config

import (
	"fmt"
	"golang.org/x/net/context"
	"pirs.io/commons"
	"pirs.io/commons/db/mongo"
	"pirs.io/dependency-management/detection"
	"time"
)

var (
	log                           = commons.GetLoggerFor("config")
	depAppCtx *ApplicationContext = nil
)

// A DependencyAppConfig contains loaded data from ENV config file
type DependencyAppConfig struct {
	GrpcIp             string `mapstructure:"GRPC_IP"`
	GrpcPort           int    `mapstructure:"GRPC_PORT"`
	UseGrpcReflection  bool   `mapstructure:"USE_GRPC_REFLECTION"`
	StreamSeparator    string `mapstructure:"STREAM_SEPARATOR"`
	MongoHost          string `mapstructure:"MONGO_HOST"`
	MongoPort          string `mapstructure:"MONGO_PORT"`
	MongoUser          string `mapstructure:"MONGO_USER"`
	MongoPass          string `mapstructure:"MONGO_PASS"`
	MongoName          string `mapstructure:"MONGO_NAME"`
	MongoDrop          bool   `mapstructure:"MONGO_DROP"`
	MetadataCollection string `mapstructure:"METADATA_COLLECTION"`
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

func initMongoDatabase(conf DependencyAppConfig) mongo.Client {
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

func createApplicationContext(conf DependencyAppConfig) (appContext *ApplicationContext, err error) {
	mongoClient := initMongoDatabase(conf)
	metadataRepo := mongo.NewMetadataRepository(mongoClient.Database(conf.MongoName, conf.MongoDrop), conf.MetadataCollection)
	return &ApplicationContext{
		DetectionService: detection.NewDetectionService(*metadataRepo),
	}, nil
}
