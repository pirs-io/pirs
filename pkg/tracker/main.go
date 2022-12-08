package main

import (
	"pirs.io/commons"
	"pirs.io/tracker/config"
	"pirs.io/tracker/grpc"
)

var log = commons.GetLoggerFor("main")

func main() {

	appConfig := config.InitApp("./dev.env")
	// start gRpc server
	grpc.StartGrpc(appConfig.GrpcPort)
}
