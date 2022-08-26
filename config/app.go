package config

import (
	"uwwolf/util"

	"github.com/spf13/viper"
)

type appConfig struct {
	Debug bool `mapstructure:"APP_DEBUG"`
}

func (c *appConfig) load() {
	util.LoadDefaultConfigValues(map[string]any{
		"APP_DEBUG": false,
	})

	viper.Unmarshal(&c)
}
