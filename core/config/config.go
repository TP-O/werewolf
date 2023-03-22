package config

import (
	"log"

	"github.com/spf13/viper"
)

type config struct {
	App      app      `mapstructure:"app"`
	Postgres postgres `mapstructure:"postgres"`
	Game     game     `mapstructure:"game"`
	Redis    redis    `mapstructure:"redis"`
}

type configLoader interface {
	loadDefault()
}

var cfg *config

// loadDefaultConfig loads the default config values.
func loadDefaultConfig() {
	cfg = &config{}
	cfg.App.loadDefault()
	cfg.Postgres.loadDefault()
	cfg.Game.loadDefault()
	cfg.Redis.loadDefault()
}

// Load loads config values from the given path and
// uses the default values if not set.
func Load(path string) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Unable to load config:", err)
		log.Println("Use default config!")
	}

	loadDefaultConfig()
	if err := viper.Unmarshal(cfg); err != nil {
		panic(err)
	}
}
