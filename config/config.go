package config

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

var (
	App   *appConfig   = &appConfig{}
	Game  *gameConfig  = &gameConfig{}
	DB    *dbConfig    = &dbConfig{}
	Cache *cacheConfig = &cacheConfig{}
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
		log.Fatalln("Loading environment variables is failed!", err)
	}

	load(App, Game, DB, Cache)
}

func load(loaders ...configLoader) {
	for _, loader := range loaders {
		loader.load()
	}
}
