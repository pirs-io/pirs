package main

import (
	"path/filepath"
	"pirs.io/pirs/common"
	"pirs.io/pirs/tracker/config"
)

var log = common.GetLoggerFor("main")

func main() {

	bs, err := filepath.Abs("./tracker.env")
	if err != nil {
		panic(err)
	}
	config.InitApp(bs)

}
