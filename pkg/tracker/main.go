package main

import (
	"pirs.io/pirs/common"
	"pirs.io/pirs/tracker/config"
	"pirs.io/pirs/tracker/grpc"
)

var log = common.GetLoggerFor("main")

func main() {

	appConfig := config.InitApp("./tracker.env")
	// start gRpc server
	grpc.StartGrpc(appConfig.GrpcPort)
}
