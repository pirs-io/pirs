package main

type AppConfig struct {
	Port    int    `mapstructure:"PORT"`
	LogPath string `mapstructure:"LOG_PATH"`
}

func (c AppConfig) IsConfig() {}
