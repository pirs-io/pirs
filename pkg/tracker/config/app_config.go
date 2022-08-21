package config

import (
	"pirs.io/pirs/common"
	"pirs.io/pirs/tracker/db"
	"pirs.io/pirs/tracker/service"
)

var (
	log                                           = common.GetLoggerFor("config")
	trackerApplicationContext *ApplicationContext = nil
)

type TrackerAppConfig struct {
	RedisURl  string `mapstructure:"REDIS_URL"`
	RedisPort string `mapstructure:"REDIS_PORT"`
	RedisPwd  string `mapstructure:"REDIS_PWD"`
	GrpcPort  int    `mapstructure:"GRPC_PORT"`
}

func (t TrackerAppConfig) IsConfig() {}

type ApplicationContext struct {
	LocationService *service.LocationService
}

func InitApp(configFilePath string) (conf *TrackerAppConfig) {
	// config loading
	log.Info().Msg("Starting application")
	conf, confErr := common.GetAppConfig(configFilePath, &TrackerAppConfig{})
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
	return conf
}

func GetContext() *ApplicationContext {
	return trackerApplicationContext
}

func createApplicationContext(conf TrackerAppConfig) (appContext *ApplicationContext, err error) {
	return &ApplicationContext{
		LocationService: &service.LocationService{LocationRepository: &db.RedisRepo{
			Client: common.NewRedisClient(conf.RedisURl, conf.RedisPort, conf.RedisPwd, 0)},
		},
	}, nil
}
