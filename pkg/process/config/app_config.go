package config

//import (
//	"pirs.io/common"
//)

var (
	processApplicationContext *ApplicationContext = nil
)

type ProcessAppConfig struct {
	GrpcPort int `mapstructure:"GRPC_PORT"`
}

type ApplicationContext struct {
	AppConfig *ProcessAppConfig
}

func GetContext() *ApplicationContext {
	return processApplicationContext
}

func InitApp(configFilePath string) (conf *ProcessAppConfig) {
	// config loading
	//conf, confErr := common.GetAppConfig(configFilePath, &ProcessAppConfig{})
	//if confErr != nil {
	//	panic("use logging")
	//}

	// create app context
	appCtx, contextErr := createApplicationContext(*conf)
	if contextErr != nil {
		panic("use logging")
	}
	processApplicationContext = appCtx

	processApplicationContext.AppConfig = conf
	return conf
}

func createApplicationContext(conf ProcessAppConfig) (appContext *ApplicationContext, err error) {
	return &ApplicationContext{}, nil
}
