package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"pirs.io/commons"
)

func main() {
	path := "./.env"
	if fileExists("./dev.env") {
		path = "./dev.env"
	}

	config, err := commons.GetAppConfig(path, &AppConfig{})
	if err != nil {
		fmt.Printf("Cannot load config! %s", err)
		panic(err)
	}
	commons.PrintConfigFields(*config)
	commons.InitLogger(config.LogPath, 10, 3, 30)
	log.Info().Msg("User service application has started")
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}
