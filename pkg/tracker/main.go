package main

import (
	"pirs.io/common"
	"pirs.io/tracker/config"
	"pirs.io/tracker/grpc"
)

var log = common.GetLoggerFor("main")

func main() {

	appConfig := config.InitApp("./tracker-dev.env")
	// start gRpc server
	grpc.StartGrpc(appConfig.GrpcPort)
}
