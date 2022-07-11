package config

import (
	"uwwolf/util"

	"github.com/spf13/viper"
)

type gameConfig struct {
	MinCapacity uint `mapstructure:"GAME_MIN_CAPACITY"`
	MaxCapacity uint `mapstructure:"GAME_MAX_CAPACITY"`
}

var Game *gameConfig

func loadGameConfig() {
	util.LoadDefaultConfigValues(map[string]interface{}{
		"GAME_MIN_CAPACITY": 5,
		"GAME_MAX_CAPACITY": 20,
	})

	viper.Unmarshal(&Game)
}
