package config

import (
	"github.com/spf13/viper"
)

type appConfig struct {
	Debug    bool   `mapstructure:"debug"`
	Env      string `mapstructure:"env"`
	GrpcPort uint16 `mapstructure:"api_port"`
	WsPort   uint16 `mapstructure:"ws_port"`
}

func loadDefaultApp() {
	viper.SetDefault("app", map[string]interface{}{
		"debug":    true,
		"env":      "development",
		"api_port": 8080,
		"ws_port":  8081,
	})
}
