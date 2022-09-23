package config

import (
	"github.com/spf13/viper"

	"uwwolf/app/util"
)

type appConfig struct {
	Debug bool `mapstructure:"APP_DEBUG"`
	Port  int  `mapstructure:"APP_PORT"`
}

func (c *appConfig) load() {
	util.LoadDefaultConfigValues(map[string]any{
		"APP_DEBUG": false,
		"APP_PORT":  80,
	})

	viper.Unmarshal(&c)
}
