package main

import "network.pirs.io/tracker/config"

var log = config.GetLoggerFor("main")

func main() {

	config.InitApp()

}
