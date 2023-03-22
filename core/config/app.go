package config

import "github.com/spf13/viper"

type app struct {
	Debug bool   `mapstructure:"debug"`
	Env   string `mapstructure:"env"`
	Port  uint16 `mapstructure:"port"`
}

var _ configLoader = (*app)(nil)

func App() app {
	return cfg.App
}

func (app) loadDefault() {
	viper.SetDefault("app", map[string]interface{}{
		"debug": true,
		"env":   "development",
		"port":  8080,
	})
}
