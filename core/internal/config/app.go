package config

import "github.com/spf13/viper"

type App struct {
	Debug     bool   `mapstructure:"debug"`
	Env       string `mapstructure:"env"`
	SecretKey string `mapstructure:"secretKey"`
	Port      uint16 `mapstructure:"port"`
}

var _ configLoader = (*App)(nil)

func (App) loadDefault() {
	viper.SetDefault("app", map[string]interface{}{
		"debug": true,
		"env":   "development",
		"port":  8080,
	})
}
