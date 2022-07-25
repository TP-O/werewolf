package config

import (
	"uwwolf/util"

	"github.com/spf13/viper"
)

type appConfig struct {
	Debug bool `mapstructure:"DEBUG"`
}

func (c *appConfig) load() {
	util.LoadDefaultConfigValues(map[string]interface{}{
		"DEBUG": false,
	})

	viper.Unmarshal(&c)
}
