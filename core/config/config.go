package config

import (
	"log"

	"github.com/spf13/viper"
)

type configLoader interface {
	loadDefault()
}

type config struct {
	App      `mapstructure:"app"`
	Postgres `mapstructure:"postgres"`
	Game     `mapstructure:"game"`
	Redis    `mapstructure:"redis"`
}

var cfg *config

// loadDefaultConfig loads the default config values.
func loadDefaultConfig(cfg *config) {
	cfg.App.loadDefault()
	cfg.Postgres.loadDefault()
	cfg.Game.loadDefault()
	cfg.Redis.loadDefault()
}

// Load loads config values from the given path and
// uses the default values if not set.
func Load(path string) *config {
	if cfg == nil {
		cfg = &config{}
		viper.AddConfigPath(path)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")

		if err := viper.ReadInConfig(); err != nil {
			log.Println("Unable to load config:", err)
			log.Println("Use default config!")
		}

		loadDefaultConfig(cfg)
		if err := viper.Unmarshal(cfg); err != nil {
			panic(err)
		}
	}

	return cfg
}
