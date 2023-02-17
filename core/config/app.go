package config

import (
	"github.com/spf13/viper"
)

type appConfig struct {
	Debug bool   `mapstructure:"debug"`
	Env   string `mapstructure:"env"`
	Port  uint16 `mapstructure:"port"`
}

func loadDefaultApp() {
	viper.SetDefault("app", map[string]interface{}{
		"debug": true,
		"env":   "development",
		"port":  8080,
	})
}
