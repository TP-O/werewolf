package config

import (
	"github.com/spf13/viper"

	"uwwolf/util"
)

type appConfig struct {
	Debug    bool `mapstructure:"APP_DEBUG"`
	HttpPort int  `mapstructure:"APP_HTTP_PORT"`
	WsPort   int  `mapstructure:"APP_WS_PORT"`
}

func (c *appConfig) load() {
	util.LoadDefaultConfigValues(map[string]any{
		"APP_DEBUG":     false,
		"APP_HTTP_PORT": 80,
		"APP_WS_PORT":   81,
	})

	viper.Unmarshal(&c)
}
