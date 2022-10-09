package main

import (
	"pirs.io/common"
	"pirs.io/tracker/config"
)

var log = common.GetLoggerFor("main")

func main() {

	appConfig := config.InitApp("./tracker-dev.env")
	// start gRpc server
	defer func(grpcPort int) {
		err := StartGrpc(grpcPort)
		if err != nil {
			panic(err)
		}
	}(appConfig.GrpcPort)
}
