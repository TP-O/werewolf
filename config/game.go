package config

import (
	"github.com/spf13/viper"

	"uwwolf/util"
)

type gameConfig struct {
	MinCapacity  int `mapstructure:"GAME_MIN_CAPACITY"`
	MaxCapacity  int `mapstructure:"GAME_MAX_CAPACITY"`
	MaxStartTime int `mapstructure:"GAME_MAX_START_TIME"`
}

func (c *gameConfig) load() {
	util.LoadDefaultConfigValues(map[string]any{
		"GAME_MIN_CAPACITY":   5,
		"GAME_MAX_CAPACITY":   20,
		"GAME_MAX_START_TIME": 30,
	})

	viper.Unmarshal(&c)
}
