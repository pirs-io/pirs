package commons

import (
	"fmt"
	"github.com/spf13/viper"
	"reflect"
)

type BaseConfig interface {
	IsConfig()
}

func GetAppConfig[T BaseConfig](configFilePath string, c *T) (res *T, err error) {
	viper.SetConfigType("env")
	viper.SetConfigFile(configFilePath)
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(c)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
		return nil, err
	}
	return c, nil
}

func PrintConfigFields(config BaseConfig) {
	reflected := reflect.ValueOf(config)
	for i := 0; i < reflected.NumField(); i++ {
		fmt.Printf("%s: %v\n", reflected.Type().Field(i).Name, reflected.Field(i).Interface())
	}
}
