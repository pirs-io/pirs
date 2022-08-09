package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var log = GetLoggerFor("config")

type AppConfig struct {
	RedisURl string `mapstructure:"REDIS_URL"`
	RedisPwd string `mapstructure:"REDIS_PWD"`
}

func InitApp() {
	log.Info().Msg("Starting application")
	config, err := GetAppConfig()
	if err != nil {
		log.Fatal().Msgf("Unable to load application config! %s", err)
	}
	log.Debug().Msgf("%s", config)
}

func GetAppConfig() (config *AppConfig, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
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
