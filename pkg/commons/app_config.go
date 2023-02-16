package commons

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type BaseConfig interface {
	IsConfig()
}

// GetAppConfig reads .env file specified but system ENV variables override .env file values
func GetAppConfig[T BaseConfig](configFilePath string, c *T) (res *T, err error) {
	viper.SetConfigType("env")
	viper.SetConfigFile(configFilePath)
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Warn().Msg("No .env file found! Environment variables will be used")
	}

	err = viper.Unmarshal(c)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
		return nil, err
	}
	return c, nil
}
