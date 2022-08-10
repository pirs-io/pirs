package config

import (
	"pirs.io/pirs/common"
)

var log = common.GetLoggerFor("config")

type TrackerAppConfig struct {
	RedisURl string `mapstructure:"REDIS_URL"`
	RedisPwd string `mapstructure:"REDIS_PWD"`
}

func (t TrackerAppConfig) IsConfig() {}

func InitApp(configFilePath string) {
	log.Info().Msg("Starting application")
	conf, err := common.GetAppConfig(configFilePath, &TrackerAppConfig{})
	if err != nil {
		log.Fatal().Msgf("Unable to load application config! %s", err)
	}
	log.Debug().Msgf("%s", conf)
}
