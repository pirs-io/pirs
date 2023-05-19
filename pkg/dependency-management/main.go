package main

import (
	"pirs.io/commons"
	"pirs.io/dependency-management/config"
	"pirs.io/dependency-management/server"
)

var log = commons.GetLoggerFor("main")
var configFilePath = "./pkg/dependency-management/.env"

func main() {
	appConfig := config.InitApp(configFilePath)
	err := server.StartGrpc(appConfig.GrpcIp, appConfig.GrpcPort, appConfig.UseGrpcReflection)
	if err != nil {
		log.Error().Msgf("Failed to start Dependency Management service: %v", err)
		return
	}
}
