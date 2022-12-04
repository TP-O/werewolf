package config

import (
	"github.com/spf13/viper"
)

type gameConfig struct {
	MinCapacity     uint8 `mapstructure:"min_capacity"`
	MaxCapacity     uint8 `mapstructure:"min_capacity"`
	MinPollCapacity uint8 `mapstructure:"min_poll_capacity"`
	PreparationTime uint8 `mapstructure:"preparation_time"`
}

func loadDefaultGame() {
	viper.SetDefault("game", map[string]interface{}{
		"min_capacity":      5,
		"max_capacicty":     20,
		"min_poll_capacity": 3,
		"preparation_time":  10,
	})
}
