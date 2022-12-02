package config

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

var (
	App  *appConfig  = &appConfig{}
	DB   *dbConfig   = &dbConfig{}
	Game *gameConfig = &gameConfig{}
)

type configLoader interface {
	load()
}

func init() {
	_, filePath, _, _ := runtime.Caller(0)

	// Current directory
	rootPath := filepath.Dir(filePath)

	viper.SetConfigFile(rootPath + "/../.env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("Unable to load .env:", err)
	}

	load(App, DB, Game)
}

func load(loaders ...configLoader) {
	for _, loader := range loaders {
		loader.load()
	}
}
