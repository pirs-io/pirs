package main

import (
	"pirs.io/commons"
	"pirs.io/process/config"
	"pirs.io/process/grpc"
)

var log = commons.GetLoggerFor("main")
var configFilePath = "./pkg/process-dev.env"

func main() {
	appConfig := config.InitApp(configFilePath)
	err := grpc.StartGrpc(appConfig.GrpcIp, appConfig.GrpcPort, appConfig.UseGrpcReflection)
	if err != nil {
		log.Error().Msg("Failed to start Process service.")
		return
	}
}
