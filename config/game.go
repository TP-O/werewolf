package config

import (
	"uwwolf/util"

	"github.com/spf13/viper"
)

type gameConfig struct {
	MinCapacity int `mapstructure:"GAME_MIN_CAPACITY"`
	MaxCapacity int `mapstructure:"GAME_MAX_CAPACITY"`
}

func (c *gameConfig) load() {
	util.LoadDefaultConfigValues(map[string]interface{}{
		"GAME_MIN_CAPACITY": 5,
		"GAME_MAX_CAPACITY": 20,
	})

	viper.Unmarshal(&c)
}
