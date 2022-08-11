package config

import (
	"pirs.io/pirs/common"
	"pirs.io/pirs/tracker/grpc"
)

var (
	log = common.GetLoggerFor("config")
)

type TrackerAppConfig struct {
	RedisURl string `mapstructure:"REDIS_URL"`
	RedisPwd string `mapstructure:"REDIS_PWD"`
	GrpcPort int    `mapstructure:"GRPC_PORT"`
}

func (t TrackerAppConfig) IsConfig() {}

func InitApp(configFilePath string) (conf *TrackerAppConfig) {
	// config loading
	log.Info().Msg("Starting application")
	conf, err := common.GetAppConfig(configFilePath, &TrackerAppConfig{})
	if err != nil {
		log.Fatal().Msgf("Unable to load application config! %s", err)
	}

	// start gRpc server
	grpc.StartGrpc(conf.GrpcPort)

	return conf
}
