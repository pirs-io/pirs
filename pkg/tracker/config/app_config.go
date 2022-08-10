package config

import (
	"fmt"
	"github.com/spf13/viper"
	"pirs.io/pirs/common"
)

var log = common.GetLoggerFor("config")

type AppConfig struct {
	RedisURl string `mapstructure:"REDIS_URL"`
	RedisPwd string `mapstructure:"REDIS_PWD"`
}

func InitApp(configFilePath string) {
	log.Info().Msg("Starting application")
	config, err := GetAppConfig(configFilePath)
	if err != nil {
		log.Fatal().Msgf("Unable to load application config! %s", err)
	}
	log.Debug().Msgf("%s", config)
}

// TODO: move to common
func GetAppConfig(configFilePath string) (config *AppConfig, err error) {
	viper.SetConfigType("env")
	viper.SetConfigFile(configFilePath)
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return &AppConfig{}, err
	}

	conf := &AppConfig{}
	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
		return &AppConfig{}, err
	}
	return conf, nil
}
