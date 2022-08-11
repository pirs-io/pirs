package main

import (
	"pirs.io/pirs/common"
	"pirs.io/pirs/tracker/config"
)

var log = common.GetLoggerFor("main")

func main() {

	appConfig := config.InitApp("./tracker.env")
	log.Info().Msg(appConfig.RedisPwd)

}
