package config

import (
	"uwwolf/util"

	"github.com/spf13/viper"
)

type appConfig struct {
	Debug bool `mapstructure:"DEBUG"`
}

var App *appConfig

func loadAppConfig() {
	util.LoadDefaultConfigValues(map[string]interface{}{
		"DEBUG": false,
	})

	viper.Unmarshal(&App)
}
