package main

import (
	"pirs.io/process-storage/config"
)

func main() {
	appConfig := config.InitApp("./process-storage--dev.env")
	// start gRpc server
	defer func(grpcPort int) {
		err := StartGrpc(grpcPort)
		if err != nil {
			panic(err)
		}
	}(appConfig.GrpcPort)
}
