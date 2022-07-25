package config

import (
	"log"

	"github.com/spf13/viper"
)

var (
	App  *appConfig  = &appConfig{}
	Game *gameConfig = &gameConfig{}
	DB   *dbConfig   = &dbConfig{}
)

type configLoader interface {
	load()
}

func init() {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Loading environment variables is failed!")
	}

	load(App, Game, DB)
}

func load(loaders ...configLoader) {
	for _, loader := range loaders {
		loader.load()
	}
}
