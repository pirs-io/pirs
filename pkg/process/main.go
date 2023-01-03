package main

import (
	"pirs.io/commons"
	"pirs.io/process/config"
	"pirs.io/process/server"
)

var log = commons.GetLoggerFor("main")
var configFilePath = "./pkg/process/dev.env"

func main() {
	appConfig := config.InitApp(configFilePath)
	err := server.StartGrpc(appConfig.GrpcIp, appConfig.GrpcPort, appConfig.UseGrpcReflection)
	if err != nil {
		log.Error().Msgf("Failed to start Process service: %v", err)
		return
	}
}
