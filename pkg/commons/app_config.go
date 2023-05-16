package commons

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
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
		log.Warn().Msgf("No .env file on path %s found! Environment variables will be used", configFilePath)
	}

	err = viper.Unmarshal(c)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
		return nil, err
	}
	return c, nil
}

func CreateGrpcOTELInterceptors() (grpc.ServerOption, grpc.ServerOption) {
	return grpc.ChainUnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.ChainStreamInterceptor(otelgrpc.StreamServerInterceptor())
}
