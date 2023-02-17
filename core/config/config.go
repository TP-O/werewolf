package config

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

type config struct {
	App  appConfig  `mapstructure:"app"`
	Game gameConfig `mapstructure:"game"`
	DB   dbConfig   `mapstructure:"database"`
}

var cfg config

func init() {
	_, filePath, _, _ := runtime.Caller(0)
	currentPath := filepath.Dir(filePath)
	viper.AddConfigPath(currentPath + "/../")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("Unable to load config:", err)
	}

	loadDefaultApp()
	loadDefaultDB()
	loadDefaultGame()

	viper.Unmarshal(&cfg)
}

func App() appConfig {
	return cfg.App
}

func DB() dbConfig {
	return cfg.DB
}

func Game() gameConfig {
	return cfg.Game
}
