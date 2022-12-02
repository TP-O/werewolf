package config

import (
	"github.com/spf13/viper"

	"uwwolf/util"
)

type gameConfig struct {
	MinCapacity     int  `mapstructure:"GAME_MIN_CAPACITY"`
	MaxCapacity     int  `mapstructure:"GAME_MIN_CAPACITY"`
	MinPollCapacity uint `mapstructure:"GAME_MIN_POLL_CAPACITY"`
	PreparationTime int  `mapstructure:"GAME_PREPARATION_TIME"`
}

func (c *gameConfig) load() {
	util.LoadDefaultConfigValues(map[string]any{
		"GAME_MIN_CAPACITY":      5,
		"GAME_MAX_CAPACITY":      20,
		"GAME_MIN_POLL_CAPACITY": 3,
		"GAME_PREPARATION_TIME":  10,
	})

	viper.Unmarshal(&c)
}
