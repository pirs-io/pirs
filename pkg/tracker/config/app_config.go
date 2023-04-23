package config

import (
	"context"
	"pirs.io/commons"
	"pirs.io/tracker/db"
	"pirs.io/tracker/redis"
	"pirs.io/tracker/service"
)

var (
	log                                           = commons.GetLoggerFor("config")
	trackerApplicationContext *ApplicationContext = nil
)

type TrackerAppConfig struct {
	RedisURl     string `mapstructure:"REDIS_URL"`
	RedisPort    string `mapstructure:"REDIS_PORT"`
	RedisPwd     string `mapstructure:"REDIS_PWD"`
	GrpcPort     int    `mapstructure:"GRPC_PORT"`
	Instance0Url string `mapstructure:"INSTANCE0_URL"`
}

func (t TrackerAppConfig) IsConfig() { return }

type ApplicationContext struct {
	AppConfig                   *TrackerAppConfig
	InstanceRegistrationService *service.InstanceRegistrationService
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
	trackerApplicationContext.AppConfig = conf
	return conf
}

func GetContext() *ApplicationContext {
	return trackerApplicationContext
}

func createApplicationContext(conf TrackerAppConfig) (appContext *ApplicationContext, err error) {
	ctx := context.Background()
	redisClient := redis.NewRedisClient(ctx, conf.RedisURl, conf.RedisPort, conf.RedisPwd, 0)
	return &ApplicationContext{
		InstanceRegistrationService: &service.InstanceRegistrationService{RegisterRepo: &db.RegisterRepo{
			Context: &ctx,
			Client:  redisClient},
		},
	}, nil
}
