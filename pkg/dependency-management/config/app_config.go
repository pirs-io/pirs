package config

import (
	"golang.org/x/net/context"
	"pirs.io/commons"
	"pirs.io/commons/db/mongo"
	"pirs.io/commons/parsers"
	"pirs.io/dependency-management/detection"
	"strings"
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
	MongoUri           string `mapstructure:"MONGO_URI"`
	MongoName          string `mapstructure:"MONGO_NAME"`
	MongoDrop          bool   `mapstructure:"MONGO_DROP"`
	MetadataCollection string `mapstructure:"METADATA_COLLECTION"`
	PetriflowApi       string `mapstructure:"PETRIFLOW_API_CSV"`
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

func parseApiForDetectionFromCsv(csvPath string) map[string][]string {
	csv := parsers.ReadCsvFile(csvPath, true)
	if csv == nil {
		log.Warn().Msg("CSV " + csvPath + " was not found.")
	}
	result := map[string][]string{}
	for idx, row := range csv {
		// skip header
		if idx == 0 {
			continue
		}
		fromDelim := strings.Replace(row[1], "\\", "", 1)
		untilDelim := strings.Replace(row[2], "\\", "", 1)
		result[row[0]] = []string{fromDelim, untilDelim}
	}
	return result
}

func createApplicationContext(conf DependencyAppConfig) (appContext *ApplicationContext, err error) {
	mongoClient := initMongoDatabase(conf)
	metadataRepo := mongo.NewMetadataRepository(mongoClient.Database(conf.MongoName, conf.MongoDrop), conf.MetadataCollection)
	petriflowApi := parseApiForDetectionFromCsv(conf.PetriflowApi)
	return &ApplicationContext{
		DetectionService: detection.NewDetectionService(metadataRepo, petriflowApi),
	}, nil
}
