package main

import (
	"pirs.io/process/grpc"
)

func main() {
	//appConfig := config.InitApp("./tracker-dev.env")
	err := grpc.StartGrpc(8090)
	if err != nil {
		return
	}
}
